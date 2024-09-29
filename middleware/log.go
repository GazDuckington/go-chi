package middleware

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func LogCalls(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		log.Printf("%v → %v %v\n", r.URL, r.RemoteAddr, ww.Status())
	})
}
