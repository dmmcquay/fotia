package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

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
	rand.Seed(time.Now().UnixNano())
}

// Time is a function that makes it simple to add one-line timings to function
// calls.
func Time() func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start)
		pc := make([]uintptr, 10)
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])

		Func.WithLabelValues(f.Name()).Observe(float64(elapsed / time.Microsecond))
	}
}

var addr = flag.String("addr", ":8080", "The address to listen on for HTTP requests.")

func demo(w http.ResponseWriter, r *http.Request) {
	defer Time()()
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func sleep(w http.ResponseWriter, r *http.Request) {
	defer Time()()
	i := rand.Intn(10)
	fmt.Fprintf(w, "sleeping for %d", i)
	time.Sleep(time.Second * time.Duration(i))
}

func main() {
	flag.Parse()
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/demo", demo)
	http.HandleFunc("/sleep", sleep)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
