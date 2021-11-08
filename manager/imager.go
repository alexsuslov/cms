package manager

import (
	"github.com/alexsuslov/cms"
	"github.com/gorilla/mux"
)

func Imager(sub *mux.Router, ext string, Options *cms.Options) {
	static := Env("STATIC", "static/")

	p := "/" + ext
	l := static + "/" + ext
	w := "/admin/" + ext

	sub.HandleFunc(p,
		Files(l, w, *Options)).Methods("GET")
	sub.HandleFunc(p,
		FileUpload(l, w, *Options)).
		Methods("POST")
}
