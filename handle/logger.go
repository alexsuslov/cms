package handle

import (
	"encoding/json"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/store"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
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

func LoggingMiddlewareDB(s *store.Store, o *cms.Options) func(next http.Handler) http.Handler {
	hint := counter(s, []byte("visits"), o)
	hintIP := counter(s, []byte("ips"), o)
	hint404 := counter(s, []byte("notfound"), o)
	hint403 := counter(s, []byte("forbidden"), o)
	return func(next http.Handler) http.Handler {

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


			key:=strings.Split(ip,":")

			switch wrapped.status {
			case 404:
				go hint404(ip, []byte(r.URL.EscapedPath()))
			case 403:
				go hint403(ip, []byte(key[0]))
			default:
				go hint(ip, []byte(r.URL.EscapedPath()))
				go hintIP(ip, []byte(key[0]))
			}

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
}

type skip []string
func NewSkipper(data interface{})(s skip){
	s=[]string{}
	ips, ok := data.([]interface{})
	if !ok{
		return
	}
	for _, ip := range ips{
		s= append(s, ip.(string))
	}
	return
}
func (skip skip) Is(ip string) bool {
	for _, prefix := range skip {
		if strings.HasPrefix(ip, prefix) {
			return true
		}
	}
	return false
}

func counter(s *store.Store, bucketName []byte, o *cms.Options) func(ip string, key []byte) {
	Skiper := NewSkipper(o.Get("skip"))

	return func(ip string, key []byte) {
		if Skiper.Is(ip) {
			return
		}
		err := s.DB.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists(bucketName)
			if err != nil {
				return err
			}
			count := 0
			data := b.Get(key)
			if data != nil {
				if err := json.Unmarshal(data, &count); err != nil {
					return err
				}
			}
			count++
			data, err = json.Marshal(count)
			if err != nil {
				return err
			}
			return b.Put(key, data)
		})
		if err != nil {
			log.Println("counter:", err)
		}
	}
}
