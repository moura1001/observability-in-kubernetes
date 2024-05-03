package main

import (
	"fmt"
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
}

func endpointHandler2(w http.ResponseWriter, r *http.Request) {
	requestsNode.Inc()
	requestsTotal.Inc()

	// Simulate some work
	time.Sleep(50 * time.Millisecond)

	// Respond
	fmt.Fprintln(w, fmt.Sprintf("Response from endpoint 2 from node %s", HOSTNAME))
}

func main() {
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
