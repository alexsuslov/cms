package handle

import (
	"github.com/alexsuslov/cms"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type ITemplate interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
}


func Home(T ITemplate, o cms.Options) func(w http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		err := T.ExecuteTemplate(w, "home", o)
		if err!= nil{
			logrus.Error(err)
		}
	}
}