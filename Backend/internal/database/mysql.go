package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/BHAV0207/Requip-assignment/internal/config"
	"github.com/XSAM/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens the MySQL pool. It retries for about a minute because in
// Kubernetes the API pod can start before MySQL is ready to accept connections.
func Connect(ctx context.Context, cfg *config.Config, log *slog.Logger) (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC&timeout=5s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	// otelsql wraps the database/sql driver, so every query becomes a child
	// span of the HTTP request that issued it.
	sqlDB, err := otelsql.Open("mysql", dsn, otelsql.WithAttributes(semconv.DBSystemMySQL))
	if err != nil {
		return nil, err
	}

	var db *gorm.DB
	for attempt := 1; attempt <= 12; attempt++ {
		db, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
		if err == nil {
			break
		}
		log.Warn("database not reachable, retrying", "attempt", attempt, "error", err)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
	if err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	log.Info("database connected", "host", cfg.DBHost, "database", cfg.DBName)
	return db, nil
}
