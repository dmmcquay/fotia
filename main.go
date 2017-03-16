package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Errors tracks http status codes for problematic requests.
	Errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Number of upstream errors",
		},
		[]string{"status"},
	)

	// Func tracks time spent in a function.
	Func = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "function_microseconds",
			Help: "function timing.",
		},
		[]string{"route"},
	)

	// DB tracks timing of interactions with the file system.
	Msgs = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "msgs",
			Help: "silly msgs.",
		},
		[]string{"msgs"},
	)
)

func init() {
	prometheus.MustRegister(Errors)
	prometheus.MustRegister(Func)
	prometheus.MustRegister(Msgs)
}

var addr = flag.String("addr", ":8080", "The address to listen on for HTTP requests.")

func demo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	flag.Parse()
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/demo", demo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
