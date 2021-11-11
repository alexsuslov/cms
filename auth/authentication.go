package auth

import (
	"fmt"
	"github.com/alexsuslov/cms/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type IStore interface{
	GetUser(username string) (user *model.User, err error)
}

type AuthenticationMiddleware struct {
	Store IStore
	Roles []string
}

func NewAuthMid(Store IStore, roles ...string) *AuthenticationMiddleware {
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
		expire := time.Now().AddDate(0, 0, 1)
		cookie := http.Cookie{
			Name:"editor",
			Value:"Yes",
			Path: "/",
			Domain: r.URL.Host,
			Expires: expire,
			RawExpires: expire.Format(time.UnixDate),
		}
		http.SetCookie(w, &cookie)

		r.Header.Set("roles", strings.Join(user.GetRoles(), "|"))
		next.ServeHTTP(w, r)

	})
}
