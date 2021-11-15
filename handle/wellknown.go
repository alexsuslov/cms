package handle

import (
	"github.com/alexsuslov/cms"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

//WELL-KNOWN
func WellKnown(r *mux.Router) {
	security:=cms.Env("SECURITY", "security.txt")
	data, err := ioutil.ReadFile(security)
	if err==nil{
		log.Println("well-known")
		r.HandleFunc("/.well-known/security.txt", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			w.Write(data)
		})
	}
}
