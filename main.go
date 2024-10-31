package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	dummyCounterLabels = [][]string{ // interval_seconds, step
		{"3", "1"},
		{"5", "2"},
		{"10", "3"},
	}

	dummyCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "dummy_counter",
		Help: "My Dummy Counter",
	},
		[]string{"interval_seconds", "step"},
	)
)

func recordDummyCounter(label []string) {
	intervalSeconds, _ := strconv.ParseUint(label[0], 10, 64)
	step, _ := strconv.ParseFloat(label[1], 64)

	for {
		dummyCounter.WithLabelValues(label...).Add(step)
		time.Sleep(time.Duration(intervalSeconds) * time.Second)
	}
}

func main() {
	for _, dummyCounterLabel := range dummyCounterLabels {
		go recordDummyCounter(dummyCounterLabel)
	}

	r := prometheus.NewRegistry()
	r.MustRegister(dummyCounter)
	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	http.Handle("/metrics", handler)
	http.ListenAndServe(":2112", nil)
}
