package middleware

import (
	"log"
	"net/http"
	"os"
)

func UseAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedUser := os.Getenv("basicuser")
		expectedPass := os.Getenv("basicpass")

		// only care about basic auth if it's set
		if expectedUser != "" && expectedPass != "" {
			user, pass, ok := r.BasicAuth()
			if !ok || user != expectedUser || pass != expectedPass {
				log.Printf("unauthorized request: %+v\n", r)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
