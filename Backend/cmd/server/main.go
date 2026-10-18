package main

import (
	"context"
	"errors"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/BHAV0207/Requip-assignment/internal/config"
	"github.com/BHAV0207/Requip-assignment/internal/database"
	"github.com/BHAV0207/Requip-assignment/internal/observability"
	"github.com/BHAV0207/Requip-assignment/internal/response"
	"github.com/BHAV0207/Requip-assignment/internal/user"
	"github.com/BHAV0207/Requip-assignment/internal/version"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const serviceName = "requip-backend"

func main() {

	cfg := config.Load()
	log := observability.NewLogger(cfg.LogLevel)

	// SIGTERM is what Kubernetes sends before killing a pod.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	// Tracing is diagnostic, not part of serving users: if it cannot start,
	// run without it rather than taking the API down.
	shutdownTracing, err := observability.InitTracing(ctx, serviceName, version.Version)
	if err != nil {
		log.Warn("tracing disabled: setup failed", "error", err)
		shutdownTracing = func(context.Context) error { return nil }
	}

	db, err := database.Connect(ctx, cfg, log)
	if err != nil {
		log.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	sqlDB, _ := db.DB()

	if cfg.RunMigrations {
		if err := db.AutoMigrate(&user.User{}); err != nil {
			log.Error("migration failed", "error", err)
			os.Exit(1)
		}
		log.Info("migrations applied")
	}

	prometheus.MustRegister(collectors.NewDBStatsCollector(sqlDB, "user_management"))
	observability.BuildInfo.WithLabelValues(version.Version, version.Commit).Set(1)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(
		otelgin.Middleware(serviceName, otelgin.WithFilter(func(r *http.Request) bool {
			return !strings.HasPrefix(r.URL.Path, "/healthz") &&
				!strings.HasPrefix(r.URL.Path, "/readyz") &&
				r.URL.Path != "/metrics"
		})),
		observability.Metrics(),
		observability.RequestLogger(log),
		// Recovery sits inside metrics and logging so a panic is still counted as a 500.
		gin.Recovery(),
		cors.New(cors.Config{
			AllowOrigins: cfg.CORSAllowedOrigins,
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Content-Type"},
			MaxAge:       12 * time.Hour,
		}),
	)

	// Liveness: the process is up. It must not check the database, or a
	// database outage would make Kubernetes restart every healthy API pod.
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Readiness: this pod can serve traffic. When MySQL is unreachable the pod
	// is removed from the Service endpoints instead of returning errors.
	router.GET("/readyz", func(c *gin.Context) {
		pingCtx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(pingCtx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "database": "down"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready", "database": "up"})
	})

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	versionInfo := gin.H{"version": version.Version, "commit": version.Commit, "buildTime": version.BuildTime}
	router.GET("/version", func(c *gin.Context) { c.JSON(http.StatusOK, versionInfo) })

	api := router.Group("/api/v1")
	api.GET("/version", func(c *gin.Context) { c.JSON(http.StatusOK, versionInfo) })
	if cfg.FaultErrorRate > 0 {
		log.Warn("fault injection enabled", "error_rate", cfg.FaultErrorRate)
		api.Use(func(c *gin.Context) {
			if rand.Float64() < cfg.FaultErrorRate {
				response.Error(c, http.StatusInternalServerError, "injected fault")
				c.Abort()
			}
		})
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	user.RegisterRoutes(api, userHandler)

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Info("server started", "port", cfg.AppPort, "version", version.Version, "commit", version.Commit)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("shutdown signal received, draining connections")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "error", err)
	}
	_ = shutdownTracing(shutdownCtx)
	_ = sqlDB.Close()
	log.Info("server stopped")
}
