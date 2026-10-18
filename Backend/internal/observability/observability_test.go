package observability

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// Regression test: with an OTLP endpoint configured (as in Kubernetes),
// tracing setup failed with "conflicting Schema URL" and crashed every pod.
func TestInitTracingWithEndpoint(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://127.0.0.1:4317")
	t.Setenv("OTEL_EXPORTER_OTLP_INSECURE", "true")

	shutdown, err := InitTracing(context.Background(), "test-service", "test")
	if err != nil {
		t.Fatalf("InitTracing returned error: %v", err)
	}
	_ = shutdown(context.Background())
}

func TestInitTracingWithoutEndpointIsNoop(t *testing.T) {
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	shutdown, err := InitTracing(context.Background(), "test-service", "test")
	if err != nil || shutdown(context.Background()) != nil {
		t.Fatalf("expected no-op tracing, got err=%v", err)
	}
}

// Panics must still be counted as 500s, and routes are recorded by template
// (/items/:id) rather than raw path, to keep label cardinality bounded.
func TestMetricsCountsPanicsAsServerErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Metrics(), gin.Recovery())
	r.GET("/items/:id", func(c *gin.Context) { panic("boom") })

	before := testutil.ToFloat64(httpRequests.WithLabelValues("GET", "/items/:id", "500"))
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/items/42", nil))
	after := testutil.ToFloat64(httpRequests.WithLabelValues("GET", "/items/:id", "500"))

	if after-before != 1 {
		t.Fatalf("expected one 500 recorded for route template, got %v", after-before)
	}
	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range families {
		for _, m := range f.GetMetric() {
			for _, l := range m.GetLabel() {
				if strings.Contains(l.GetValue(), "/items/42") {
					t.Fatalf("raw path leaked into metric %s", f.GetName())
				}
			}
		}
	}
}
