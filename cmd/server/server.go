package main

import (
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/alexsuslov/cms/handle/files"
	"github.com/alexsuslov/godotenv"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"syscall"
	"github.com/Masterminds/sprig"
)

var version = "developer preview"

func getMessage() string {
	return os.Getenv("MESSAGE") + "(" + version + ")"
}

func Env(key string, def string) string {
	v, _ := syscall.Getenv(key)
	if v == "" {
		return def
	}
	return v

}

func main() {
	log.Printf("Starting " + getMessage())

	// load env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("warrning load env", err)
	}

	// load options
	Options, err := cms.Load(Env("CONFIG", "config.yml"))
	if err != nil {
		panic(err)
	}

	//templates
	Templates := template.Must(
		template.New("base").
			Funcs(sprig.FuncMap()).
				ParseGlob(Env("TEMPLATES", "templates") + "/*.tmpl"))
	r := mux.NewRouter()

	// home
	r.HandleFunc("/",
		handle.Logger(
			handle.Home(Templates, *Options)))
	// page
	r.HandleFunc("/{filename}.html",
		handle.Logger(
			handle.Page(Templates, *Options)))

	static := Env("STATIC", "static/")
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(static))))

	sub := r.PathPrefix("/admin").Subrouter()

	// edit config
	sub.HandleFunc("/config.yml",
		files.FileEdit("config.yml", *Options)).Methods("GET")
	sub.HandleFunc("/config.yml",
		files.FileUpdate("config.yml", *Options)).Methods("POST")

	// css
	Editor(sub, "css", Options)
	Editor(sub, "js", Options)
	Imager(sub, "images", Options)
	Tmpl(sub, "templates", Options)

	httpAddr := fmt.Sprintf("%s:%s",
		Env("HTTP_HOST", "0.0.0.0"),
		Env("HTTP_PORT", "8080"))
	log.Println("listen", httpAddr)

	server := http.Server{Addr: httpAddr, Handler: r}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func Editor(sub *mux.Router, ext string, Options *cms.Options) {
	p := "/" + ext
	l := "static/" + ext
	w := "/admin/" + ext

	sub.HandleFunc(p,
		files.Files(l, w, *Options)).
		Methods("GET")

	sub.HandleFunc(p,
		files.FileUpload(l, w, *Options)).
		Methods("POST")

	p = fmt.Sprintf("/%s/{filename}", ext)
	l = fmt.Sprintf("static/%s", ext)
	w = fmt.Sprintf("/admin/%s", ext)

	sub.HandleFunc(p,
		files.PathEdit(l, w, *Options)).Methods("GET")
	sub.HandleFunc(p,
		files.PathUpdate(l, w, *Options)).Methods("POST")

}

func Imager(sub *mux.Router, ext string, Options *cms.Options) {
	p := "/" + ext
	l := "static/" + ext
	w := "/static/" + ext

	sub.HandleFunc(p,
		files.Files(l, w, *Options)).Methods("GET")
	sub.HandleFunc(p,
		files.FileUpload(l, w, *Options)).
		Methods("POST")
}

func Tmpl(sub *mux.Router, ext string, Options *cms.Options) {
	p := "/" + ext
	l := ext
	w := "/admin/" + ext

	sub.HandleFunc(p,
		files.Files(l, w, *Options)).
		Methods("GET")

	sub.HandleFunc(p,
		files.FileUpload(l, w, *Options)).
		Methods("POST")

	p = fmt.Sprintf("/%s/{filename}", ext)
	l = ext
	w = fmt.Sprintf("/admin/%s", ext)

	sub.HandleFunc(p,
		files.PathEdit(l, w, *Options)).Methods("GET")
	sub.HandleFunc(p,
		files.PathUpdate(l, w, *Options)).Methods("POST")

}
