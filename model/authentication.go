package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type AuthenticationMiddleware struct {
	Store *Store
	Roles []string
}

func NewAuthMid(Store *Store, roles ...string) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{Store, roles}
}

func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(401)
			return
		}
		user, err := amw.Store.GetUser(name)
		if err != nil {
			logrus.Errorf("Store.GetUser:%v", err)
			time.Sleep(2 * time.Second)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if !user.PassEQ(user.GenToken(pass)) {
			logrus.Errorf("PassEQ:%v", err)
			time.Sleep(2 * time.Second)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if !user.IsAllow(amw.Roles) {
			logrus.
				WithField("user", name).
				WithField("q", r.URL.String()).
				Warning("user try open some restricted")
			err = fmt.Errorf("401")
		}

		r.Header.Set("roles", strings.Join(user.GetRoles(), "|"))
		next.ServeHTTP(w, r)

	})
}
