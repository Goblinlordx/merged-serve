package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
)

type writeTracker struct {
	writer     http.ResponseWriter
	size       int
	statusCode int
}

func (w *writeTracker) Header() http.Header {
	return w.writer.Header()
}

func (w *writeTracker) Write(b []byte) (int, error) {
	n, err := w.writer.Write(b)
	w.size = w.size + n

	return n, err
}

func (w *writeTracker) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.writer.WriteHeader(statusCode)
}

func AccessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		tracker := writeTracker{writer: w}
		next.ServeHTTP(&tracker, r)
		url := r.URL
		d := time.Since(t).Truncate(1 * time.Millisecond)
		fmt.Printf("[%s] %s %s [%d] (%s) %s\n", time.Now().Format("2006-01-02T15:04:05Z"), r.Method, url, tracker.statusCode, humanize.Bytes(uint64(tracker.size)), d)
	})
}
