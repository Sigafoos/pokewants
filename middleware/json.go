package middleware

import (
	"net/http"
)

const ApplicationJSON = "application/json"

func UseJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != ApplicationJSON {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		if (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete) && r.Header.Get("Content-Type") != ApplicationJSON {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
