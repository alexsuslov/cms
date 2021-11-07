package handle

import (
	"fmt"
	"github.com/alexsuslov/cms"

	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type IUser interface {
	GenToken(pass string) []byte
	PassEQ(token []byte) bool
	IsAllow(roles []string) bool
	GetRoles() []string
}

type IStore interface {
	GetUser(user string) (IUser, error)
}

type BasicAuthFunc func(h http.HandlerFunc, roles ...string) http.HandlerFunc

func BasicAuth(t ITemplate, s IStore, o *cms.Options) BasicAuthFunc {
	onErr := Err(t, o)
	return func(h http.HandlerFunc, roles ...string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			name, pass, ok := r.BasicAuth()
			if !ok {
				onErr(w, fmt.Errorf("401"))

			}
			user, err := s.GetUser(name)
			if onErr(w, err) {
				return
			}
			if !user.PassEQ(user.GenToken(pass)) {
				onErr(w, fmt.Errorf("401"))
			}
			if !user.IsAllow(roles) {
				logrus.
					WithField("user", name).
					WithField("q", r.URL.String()).
					Warning("user try open some restricted")
				err = fmt.Errorf("401")
			}
			r.Header.Set("roles", strings.Join(user.GetRoles(), "|"))
			h(w, r)
		}
	}
}
