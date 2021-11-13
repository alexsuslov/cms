package manager

import (
	"github.com/teris-io/shortid"
	"log"
	"net/http"
	"time"
)

var restartID string

func Shut() func(w http.ResponseWriter, r *http.Request) {
	Init()
	return func(w http.ResponseWriter, r *http.Request) {
		restartID, _ = shortid.Generate()
		t.ExecuteTemplate(w, "shutdown", map[string]string{
			"id": restartID,
		})
	}
}

func Cancel() func(w http.ResponseWriter, r *http.Request) {
	Init()
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("cancel shutdown process")
		restartID = ""
		w.Write([]byte("ok"))
	}
}

func close(s *http.Server) {
	time.Sleep(10 * time.Second)
	if restartID != "" {
		restartID = ""
		s.Close()
	}
}

func Down(server *http.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("start shutdown process")
		if restartID == "" {
			http.Redirect(w, r, "/", 302)
		}
		go close(server)
	}
}
