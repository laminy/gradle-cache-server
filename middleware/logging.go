package middleware

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	RequestId  string
	Request    []byte
	Response   []byte
}

func (o *LoggingResponseWriter) WriteHeader(code int) {
	o.StatusCode = code
	o.ResponseWriter.WriteHeader(code)
}

func (o *LoggingResponseWriter) Write(response []byte) (int, error) {
	o.Response = response
	return o.ResponseWriter.Write(response)
}

func (o *LoggingResponseWriter) RandomBytes() []byte {
	token := make([]byte, 4)
	rand.Read(token)
	return token
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := new(LoggingResponseWriter)
		lrw.ResponseWriter = w
		lrw.RequestId = string(lrw.RandomBytes())
		logStr := fmt.Sprintf("[%x] %s %s %s", lrw.RequestId, r.RemoteAddr, r.Method, r.URL)
		start := time.Now().UnixNano()
		next.ServeHTTP(lrw, r)
		end := time.Now().UnixNano()
		var length int64 = getLength(w, r)
		log.Printf("%s %dbs %dns %d %s\n", logStr, length, (end-start)/1e3, lrw.StatusCode, http.StatusText(lrw.StatusCode))
	})
}

func getLength(w http.ResponseWriter, r *http.Request) (length int64) {
	switch r.Method {
	case "PUT":
		length = r.ContentLength
		return
	case "GET":
		getLen := w.Header().Get("Content-Length")
		if getLen == "" {
			break
		}
		l, err := strconv.ParseInt(getLen, 10, 64)
		if err == nil {
			length = l
		}
	}
	return
}
