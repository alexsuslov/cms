package handle

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
	"time"
)

func Logger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _, _ := r.BasicAuth()
		ip := ReadUserIP(r)
		q := r.URL.String()
		start := time.Now()
		h(w, r)
		logrus.
			WithField("user", user).
			WithField("ip", ip).
			WithField("q", q).
			WithField("t", time.Since(start).String()).
			Info("req")
	}
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

func LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logrus.WithFields(logrus.Fields{
					"err":   err,
					"trace": debug.Stack(),
				})
			}
		}()
		ip := ReadUserIP(r)
		user, _, _ := r.BasicAuth()

		start := time.Now()
		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)
		logrus.WithFields(logrus.Fields{
			"user":     user,
			"ip":       ip,
			"status":   wrapped.status,
			"method":   r.Method,
			"path":     r.URL.EscapedPath(),
			"duration": time.Since(start),
		}).Info("income")
	}

	return http.HandlerFunc(fn)
}
