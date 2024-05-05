package main

import (
	"context"
	"fmt"
	"github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HOSTNAME = strings.ReplaceAll(os.Getenv("HOSTNAME"), "-", "_")

	requestsNode = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("myapp_requests_node_%s", HOSTNAME),
			Help: fmt.Sprintf("Number of requests from node %s", HOSTNAME),
		},
	)

	requestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "myapp_requests_total",
			Help: "Total number of requests",
		},
	)

	provider *sdktrace.TracerProvider
)

func init() {
	prometheus.MustRegister(requestsNode)
	prometheus.MustRegister(requestsTotal)
}

func endpointHandler1(w http.ResponseWriter, r *http.Request) {
	requestsNode.Inc()
	requestsTotal.Inc()

	ctx := r.Context()
	tr := provider.Tracer("example-tracer")

	ctx, span := tr.Start(ctx, "endpoint1")
	defer span.End()

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	// Respond
	fmt.Fprintln(w, fmt.Sprintf("Response from endpoint 1 from node %s", HOSTNAME))
	logrus.Info("request successfully", "handler", "endpointHandler1", "host", HOSTNAME)
}

func endpointHandler2(w http.ResponseWriter, r *http.Request) {
	requestsNode.Inc()
	requestsTotal.Inc()

	ctx := r.Context()
	tr := provider.Tracer("example-tracer")

	ctx, span := tr.Start(ctx, "endpoint2")
	defer span.End()

	// Simulate some work
	time.Sleep(50 * time.Millisecond)

	// Respond
	fmt.Fprintln(w, fmt.Sprintf("Response from endpoint 2 from node %s", HOSTNAME))
	logrus.Info("request successfully", "handler", "endpointHandler2", "host", HOSTNAME)
}

func initTracer() (*sdktrace.TracerProvider, error) {
	// Create OTLP exporter
	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("otel-collector:4317"),
	)
	if err != nil {
		return nil, err
	}

	// Create trace provider with OTLP exporter
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
	)

	return provider, nil
}

func main() {
	// Initialize the tracer provider
	prov, err := initTracer()
	if err != nil {
		log.Fatalf(err.Error())
	}
	provider = prov
	defer func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	hook := graylog.NewGraylogHook("graylog:12201", map[string]interface{}{"this": "is logged every time"})
	defer hook.Flush()
	logrus.AddHook(hook)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/endpoint1", endpointHandler1)
	http.HandleFunc("/endpoint2", endpointHandler2)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
