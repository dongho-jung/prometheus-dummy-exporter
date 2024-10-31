package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	dummyCounterLabels = [][]string{ // interval_seconds, step
		{"3", "1"},
		{"15", "2"},
		{"30", "3"},
		{"10-20", "0-5"},
		{"120-300", "1-10"},
	}

	dummyCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "dummy_counter",
		Help: "My Dummy Counter",
	},
		[]string{"interval_seconds", "step"},
	)
)

// getValueFromString takes an input string in either "number-number" format or "number" format
// and returns a random integer in the range or the integer itself.
func getValueFromString(input string) (int, error) {
	// Compile the regular expression to match "number-number" format
	re := regexp.MustCompile(`^(\d+)-(\d+)$`)

	// Check if input matches "number-number" format
	if matches := re.FindStringSubmatch(input); len(matches) == 3 {
		// Convert matches to integers
		start, err1 := strconv.Atoi(matches[1])
		end, err2 := strconv.Atoi(matches[2])
		if err1 != nil || err2 != nil {
			return 0, fmt.Errorf("error converting string to integer")
		}

		// Ensure start is less than or equal to end
		if start > end {
			return 0, fmt.Errorf("start should be less than or equal to end")
		}

		// Seed the random number generator
		rand.Seed(time.Now().UnixNano())

		// Return a random integer between start and end (inclusive)
		return rand.Intn(end-start+1) + start, nil
	}

	// Check if input is a single number
	if singleNumber, err := strconv.Atoi(input); err == nil {
		return singleNumber, nil
	}

	return 0, fmt.Errorf("input does not match expected format")
}

func recordDummyCounter(label []string) {
	for {
		intervalSeconds, _ := getValueFromString(label[0])
		step, _ := getValueFromString(label[1])
		dummyCounter.WithLabelValues(label...).Add(float64(step))
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
