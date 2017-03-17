package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"s.mcquay.me/dm/vain/metrics"

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

	ECount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ecount",
		Help: "Number of times ecount has been called",
	})
)

func init() {
	prometheus.MustRegister(Errors)
	prometheus.MustRegister(Func)
	prometheus.MustRegister(ECount)
	rand.Seed(time.Now().UnixNano())
}

type failure struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func NewFailure(msg string) *failure {
	return &failure{
		Success: false,
		Error:   msg,
	}
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

func up(w http.ResponseWriter, r *http.Request) {
	searchreq := r.URL.Path[len("/up/"):]
	if len(searchreq) == 0 {
		metrics.Errors.WithLabelValues(fmt.Sprintf("%d", http.StatusBadRequest)).Add(1)
		b, _ := json.Marshal(NewFailure("url could not be parsed"))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	}
	if searchreq[len(searchreq)-1] != '/' {
		http.Redirect(w, r, "/up/"+searchreq+"/", http.StatusMovedPermanently)
		return
	}
	searchReqParsed := strings.Split(searchreq, "/")
	s, err := strconv.Atoi(searchReqParsed[0])
	if err != nil {
		metrics.Errors.WithLabelValues(fmt.Sprintf("%d", http.StatusBadRequest)).Add(1)
		b, _ := json.Marshal(NewFailure(fmt.Sprintf("could not convert %v to an int", searchReqParsed[0])))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	}
	ECount.Add(float64(s))
}

func down(w http.ResponseWriter, r *http.Request) {
	searchreq := r.URL.Path[len("/down/"):]
	if len(searchreq) == 0 {
		metrics.Errors.WithLabelValues(fmt.Sprintf("%d", http.StatusBadRequest)).Add(1)
		b, _ := json.Marshal(NewFailure("url could not be parsed"))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	}
	if searchreq[len(searchreq)-1] != '/' {
		http.Redirect(w, r, "/down/"+searchreq+"/", http.StatusMovedPermanently)
		return
	}
	searchReqParsed := strings.Split(searchreq, "/")
	s, err := strconv.Atoi(searchReqParsed[0])
	if err != nil {
		metrics.Errors.WithLabelValues(fmt.Sprintf("%d", http.StatusBadRequest)).Add(1)
		b, _ := json.Marshal(NewFailure(fmt.Sprintf("could not convert %v to an int", searchReqParsed[0])))
		http.Error(w, string(b), http.StatusBadRequest)
		return
	}
	ECount.Sub(float64(s))
}

func main() {
	flag.Parse()
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/demo", demo)
	http.HandleFunc("/sleep", sleep)
	http.HandleFunc("/down/", down)
	http.HandleFunc("/up/", up)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
