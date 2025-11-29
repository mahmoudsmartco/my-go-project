package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("[%s] %s %s\n", start.Format(time.RFC3339), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		fmt.Printf("Completed in %v\n", time.Since(start))
	})
}
