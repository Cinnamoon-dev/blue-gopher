package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *statusRecorder) WriteHeader(status int) {
	if r.Status == 0 {
		r.Status = status
		r.ResponseWriter.WriteHeader(status)
	}
}

func Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{ResponseWriter: w, Status: 0}

		start := time.Now()
		next.ServeHTTP(recorder, r)
		duration := time.Since(start)

		log.Printf("%s %s %d -> %v\n", r.Method, r.URL.Path, recorder.Status, duration)
	}
}
