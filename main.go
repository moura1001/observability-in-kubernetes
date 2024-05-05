package main

import (
	"fmt"
	"github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/sirupsen/logrus"
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
)

func init() {
	prometheus.MustRegister(requestsNode)
	prometheus.MustRegister(requestsTotal)
}

func endpointHandler1(w http.ResponseWriter, r *http.Request) {
	requestsNode.Inc()
	requestsTotal.Inc()

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	// Respond
	fmt.Fprintln(w, fmt.Sprintf("Response from endpoint 1 from node %s", HOSTNAME))
	logrus.Info("request successfully", "handler", "endpointHandler1", "host", HOSTNAME)
}

func endpointHandler2(w http.ResponseWriter, r *http.Request) {
	requestsNode.Inc()
	requestsTotal.Inc()

	// Simulate some work
	time.Sleep(50 * time.Millisecond)

	// Respond
	fmt.Fprintln(w, fmt.Sprintf("Response from endpoint 2 from node %s", HOSTNAME))
	logrus.Info("request successfully", "handler", "endpointHandler2", "host", HOSTNAME)
}

func main() {
	hook := graylog.NewGraylogHook("localhost:12201", map[string]interface{}{"this": "is logged every time"})
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
