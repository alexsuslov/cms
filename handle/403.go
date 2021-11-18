package handle

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Forbidden(w http.ResponseWriter, _ *http.Request){
		w.WriteHeader(403)
	}


func ForbiddenScan(sub *mux.Router, paths []string){
	for _, path := range paths{
		sub.HandleFunc(path, Forbidden)
	}
}
