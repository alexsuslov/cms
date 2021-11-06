package handle

import (
	"github.com/alexsuslov/cms"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Page(T ITemplate, o cms.Options) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		name, _ := params["filename"]
		err := T.ExecuteTemplate(w, name, o)
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(404)
		}
	}
}
