package fotia

import (
	"math/rand"
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
