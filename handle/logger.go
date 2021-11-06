package handle

import (
	"github.com/sirupsen/logrus"
	"net/http"
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
