package middleware

import (
	"fmt"
	"net/http"
)

func ExampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("OK WORKS")
		next.ServeHTTP(w, r)
	})
}
