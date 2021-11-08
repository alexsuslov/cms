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

func Err(t ITemplate, o IExtend) func(w http.ResponseWriter, err error) bool {
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
			logrus.Errorf("w.Write %v", err)

			w.WriteHeader(500)
			_ = t.ExecuteTemplate(w, "500", o.Extend(cms.Options{
				"Error":err.Error(),
			}))
			return true
		}
		return false
	}
}
