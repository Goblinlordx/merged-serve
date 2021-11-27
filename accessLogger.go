package main

import (
	"fmt"
	"net/http"
	"time"
)

func AccessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)
		url := r.URL
		d := time.Since(t).Truncate(1 * time.Millisecond)
		fmt.Println(time.Now().Format("[2006-01-02T15:04:05Z]"), r.Method, url, d)
	})
}
