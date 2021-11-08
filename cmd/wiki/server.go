package main

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/alexsuslov/cms/manager"
	"github.com/alexsuslov/cms/model"
	"github.com/alexsuslov/godotenv"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"syscall"
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
	log.Printf("Starting WIKI" + getMessage())

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

	Store, err := model.NewStoreBDB(Env("STORE", "db.bolt"))
	if err != nil {
		panic(err)
	}
	defer Store.Close()

	r := mux.NewRouter()

	// home
	r.HandleFunc("/",
		handle.Logger(
			handle.Home(Store, Templates, *Options)))

	// manager
	sub := r.PathPrefix("/admin").Subrouter()
	mid := model.NewAuthMid(Store, "admin")
	sub.Use(mid.Middleware)

	// config
	sub.HandleFunc("/config.yml",
		manager.FileEdit(Env("CONFIG", "config.yml"), *Options)).Methods("GET")
	sub.HandleFunc("/config.yml",
		manager.FileUpdate(Env("CONFIG", "config.yml"), *Options)).Methods("POST")

	manager.Editor(sub, "css", Options)
	manager.Editor(sub, "js", Options)
	manager.Imager(sub, "images", Options)
	manager.Templater(sub, "templates", Options)
	manager.Bucketer(Store, sub, "buckets", Options)

	// page
	r.HandleFunc("/{filename}.html",
		handle.Logger(
			handle.Page(Templates, *Options)))

	// wiki page
	r.HandleFunc("/{key}",
		handle.Logger(
			WikiPage(Templates, Store, *Options)))

	static := Env("STATIC", "static")
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(static))))

	httpAddr := fmt.Sprintf("%s:%s",
		Env("HTTP_HOST", "0.0.0.0"),
		Env("PORT", "8080"))
	log.Println("listen", httpAddr)

	server := http.Server{Addr: httpAddr, Handler: r}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

