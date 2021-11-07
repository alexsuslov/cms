package handle

import (
	"github.com/alexsuslov/cms"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type IExtend interface {
	Extend(m cms.Options) cms.Options
}

func Err(t ITemplate, o interface{}) func(w http.ResponseWriter, err error) bool {
	return func(w http.ResponseWriter, err error) bool {
		if err != nil {
			if err.Error() == "401" {
				time.Sleep(2 * time.Second)
				w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
				w.WriteHeader(401)
				return true
			}
			if err.Error() == "404" {
				w.WriteHeader(404)
				_ = t.ExecuteTemplate(w, "404", o)
				return true
			}
			w.WriteHeader(500)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				logrus.Errorf("w.Write %v", err)
			}
			return true
		}
		return false
	}
}
