package manager

import (
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/gorilla/mux"
)

func Templater(sub *mux.Router, ext string, Options *cms.Options) {
	p := "/" + ext
	l := ext
	w := "/admin/" + ext

	sub.HandleFunc(p,
		Files(l, w, *Options)).
		Methods("GET")

	sub.HandleFunc(p,
		FileUpload(l, w, *Options)).
		Methods("POST")

	p = fmt.Sprintf("/%s/{filename}", ext)
	l = ext
	w = fmt.Sprintf("/admin/%s", ext)

	sub.HandleFunc(p,
		PathEdit(l, w, *Options)).Methods("GET")
	sub.HandleFunc(p,
		PathUpdate(l, w, *Options)).Methods("POST")

}