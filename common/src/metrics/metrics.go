/*
 * metrics.go
 * Created on 23.10.2019
 * Copyright (C) 2019 Volkswagen AG, All rights reserved
 *
 */

package metrics

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/go-chi/chi"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	mMemUsage       = stats.Int64("memory_usage", "The memory usage of the process", "mb")
	memoryUsageView = &view.View{
		Name:        "memory_usage",
		Measure:     mMemUsage,
		Description: "The memory usage over time",
		Aggregation: view.LastValue(),
	}

	mGoroutines    = stats.Int64("goroutines", "The number of goroutines running", "")
	goroutinesView = &view.View{
		Name:        "goroutines",
		Measure:     mGoroutines,
		Description: "The number of goroutines running",
		Aggregation: view.LastValue(),
	}

	mUptime    = stats.Int64("uptime", "The time the service runs", "ms")
	uptimeView = &view.View{
		Name:        "uptime",
		Measure:     mUptime,
		Description: "The uptime of the service",
		Aggregation: view.LastValue(),
	}

	// 2XX
	m2XX            = stats.Int64("requests_2XX", "The count of 2XX requests", "")
	requests2XXView = &view.View{
		Name:        "requests_2XX",
		Measure:     m2XX,
		Description: "The count of 2XX requests",
		Aggregation: view.Count(),
	}

	// 3XX
	m3XX            = stats.Int64("requests_3XX", "The count of 3XX requests", "")
	requests3XXView = &view.View{
		Name:        "requests_3XX",
		Measure:     m3XX,
		Description: "The count of 3XX requests",
		Aggregation: view.Count(),
	}

	// 4XX
	m4XX = stats.Int64("requests_4XX", "The count of 4XX requests", "")

	requests4XXView = &view.View{
		Name:        "requests_4XX",
		Measure:     m4XX,
		Description: "The count of 4XX requests",
		Aggregation: view.Count(),
	}

	// 5XX
	m5XX = stats.Int64("requests_5XX", "The count of 5XX requests for the endpoint GetSpecialOffers using Get", "")

	requests5XXView = &view.View{
		Name:        "requests_5XX",
		Measure:     m5XX,
		Description: "The count of 5XX requests",
		Aggregation: view.Count(),
	}

	// Request Latency
	latencyMeasure = stats.Int64("requests_latency", "The latency of successful requests", "")

	latencyView = &view.View{
		Name:        "requests_latency",
		Measure:     latencyMeasure,
		Description: "The latency of successful requests",
		Aggregation: view.Distribution(10, 20, 50, 70, 100, 200, 500, 700, 1000, 2000, 5000, 7000, 10000),
	}
)

// GetHandler returns an http.Handler that exposes Prometheus metrics
// or an error if registering the MetricViews returns an error
func RegisterHandler(namespace string, r chi.Router) error {

	view.SetReportingPeriod(500 * time.Millisecond)

	err := view.Register(memoryUsageView, goroutinesView, uptimeView)
	if err != nil {
		return fmt.Errorf("error registering metric views: %s", err)
	}

	// The method can`t return an error, see source code
	pe, _ := prometheus.NewExporter(prometheus.Options{
		Namespace: namespace,
	})

	view.RegisterExporter(pe)

	// start reporting passive stats in the background
	reportStats()

	r.Get("/metrics", pe.ServeHTTP)

	return nil
}

func GetMetricsRecordingHandler(next http.Handler) http.Handler {
	view.Register(
		requests2XXView,
		requests3XXView,
		requests4XXView,
		requests5XXView,
		latencyView)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := statusWriter{ResponseWriter: w, status: 200}

		start := time.Now()
		next.ServeHTTP(&sw, r)
		duration := time.Since(start)

		switch sw.getStatusCategory() {
		case "2XX":
			stats.Record(r.Context(), m2XX.M(1))
			stats.Record(r.Context(), latencyMeasure.M(duration.Nanoseconds()/1000000))
			break
		case "3XX":
			stats.Record(r.Context(), m3XX.M(1))
			break
		case "4XX":
			stats.Record(r.Context(), m4XX.M(1))
			break
		case "5XX":
			stats.Record(r.Context(), m5XX.M(1))
			break
		}
	})
}

func GetMetricsRecordingHandlerForEndpoint(endpointName string) func(next http.Handler) http.Handler {

	// 2XX
	m2XX := stats.Int64(fmt.Sprintf("requests_%s_2XX", endpointName), "The count of 2XX requests", "")
	requests2XXView := &view.View{
		Name:        fmt.Sprintf("requests_%s_2XX", endpointName),
		Measure:     m2XX,
		Description: "The count of 2XX requests for endpoint " + endpointName,
		Aggregation: view.Count(),
	}

	// 3XX
	m3XX := stats.Int64(fmt.Sprintf("requests_%s_3XX", endpointName), "The count of 3XX requests", "")
	requests3XXView := &view.View{
		Name:        fmt.Sprintf("requests_%s_3XX", endpointName),
		Measure:     m3XX,
		Description: "The count of 3XX requests",
		Aggregation: view.Count(),
	}

	// 4XX
	m4XX := stats.Int64(fmt.Sprintf("requests_%s_4XX", endpointName), "The count of 4XX requests", "")

	requests4XXView := &view.View{
		Name:        fmt.Sprintf("requests_%s_4XX", endpointName),
		Measure:     m4XX,
		Description: "The count of 4XX requests",
		Aggregation: view.Count(),
	}

	// 5XX
	m5XX := stats.Int64(fmt.Sprintf("requests_%s_5XX", endpointName), "The count of 5XX requests for the endpoint GetSpecialOffers using Get", "")

	requests5XXView := &view.View{
		Name:        fmt.Sprintf("requests_%s_5XX", endpointName),
		Measure:     m5XX,
		Description: "The count of 5XX requests",
		Aggregation: view.Count(),
	}

	// Request Latency
	latencyMeasure := stats.Int64(fmt.Sprintf("requests_%s_latency", endpointName), "The latency of successful requests", "")

	latencyView := &view.View{
		Name:        fmt.Sprintf("requests_%s_latency", endpointName),
		Measure:     latencyMeasure,
		Description: "The latency of successful requests",
		Aggregation: view.Distribution(10, 20, 50, 70, 100, 200, 500, 700, 1000, 2000, 5000, 7000, 10000),
	}

	view.Register(
		requests2XXView,
		requests3XXView,
		requests4XXView,
		requests5XXView,
		latencyView)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sw := statusWriter{ResponseWriter: w, status: 200}

			start := time.Now()
			next.ServeHTTP(&sw, r)
			duration := time.Since(start)

			switch sw.getStatusCategory() {
			case "2XX":
				stats.Record(r.Context(), m2XX.M(1))
				stats.Record(r.Context(), latencyMeasure.M(duration.Nanoseconds()/1000000))
				break
			case "3XX":
				stats.Record(r.Context(), m3XX.M(1))
				break
			case "4XX":
				stats.Record(r.Context(), m4XX.M(1))
				break
			case "5XX":
				stats.Record(r.Context(), m5XX.M(1))
				break
			}
		})
	}
}

func reportStats() {

	startTime := time.Now()

	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			stats.Record(context.Background(), mMemUsage.M(int64(m.Sys)/1024/1024))

			stats.Record(context.Background(), mGoroutines.M(int64(runtime.NumGoroutine())))

			stats.Record(context.Background(), mUptime.M(int64(time.Since(startTime).Nanoseconds()/1000000)))

			time.Sleep(500 * time.Millisecond)
		}
	}()
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	return n, err
}

func (sw statusWriter) getStatusCategory() string {

	s := sw.status / 100

	switch s {
	case 2:
		return "2XX"
	case 3:
		return "3XX"
	case 4:
		return "4XX"
	case 5:
		return "5XX"
	default:
		return "XXX"
	}
}
