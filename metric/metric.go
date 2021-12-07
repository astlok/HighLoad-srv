package metric

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	hits    *prometheus.CounterVec
	errors  *prometheus.CounterVec
	cpu     *prometheus.CounterVec
	Timings *prometheus.HistogramVec
)

func Destroy() {
	prometheus.Unregister(hits)
	prometheus.Unregister(errors)
	prometheus.Unregister(Timings)
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

func New() {
	hits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "hits",
	}, []string{"status", "path"})

	cpu = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "cpu",
	}, []string{"persent"})

	errors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "errors",
	}, []string{"error"})

	Timings = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "timings",
		Buckets: []float64{0.01, 0.02, 0.03, 0.04, 0.05, 0.1,
			0.15, 0.2, 0.4, 0.6, 0.8, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12, 13, 15, 20},
	}, []string{"method", "URL"})
	prometheus.MustRegister(hits, cpu, errors, Timings)

	go func() {
		idle0, total0 := getCPUSample()
		time.Sleep(3 * time.Second)
		idle1, total1 := getCPUSample()

		idleTicks := float64(idle1 - idle0)
		totalTicks := float64(total1 - total0)
		cpuUsage := 100 * (totalTicks - idleTicks) / totalTicks

		cpu.WithLabelValues(fmt.Sprintf("%f", cpuUsage)).Inc()
		time.Sleep(time.Second * 1)
	}()
}

func CrateRequestHits(status int, r *http.Request) {
	route := mux.CurrentRoute(r)
	path, _ := route.GetPathTemplate()
	hits.WithLabelValues(strconv.Itoa(status), path).Inc()
}

func CrateRequestError(err error) {
	if err != nil {
		errors.WithLabelValues(err.Error()).Inc()
	}
}
