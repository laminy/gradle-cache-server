package middleware

import (
	"github.com/laminy/gradle-cache-server/handlers"
	"net/http"
)

func FileAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := new(LoggingResponseWriter)
		lrw.ResponseWriter = w
		lrw.RequestId = string(lrw.RandomBytes())
		next.ServeHTTP(lrw, r)
		handlers.ModifyFileAccessTime(r.RequestURI, lrw.StatusCode, r.Method)
	})
}
