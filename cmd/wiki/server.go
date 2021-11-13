package main

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/auth"
	"github.com/alexsuslov/cms/handle"
	"github.com/alexsuslov/cms/manager"
	"github.com/alexsuslov/cms/store"
	"github.com/alexsuslov/godotenv"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
	Templates := template.New("base").
		Funcs(sprig.FuncMap())

	Templates, err = Templates.ParseGlob(Env("TEMPLATES", "templates") + "/*.tmpl")
	if err!= nil{
		logrus.Error(err)
		manager.Init()
		Templates = manager.GetTemplate()
	}

	Store, err := store.NewStoreBDB(Env("STORE", "db.bolt"))
	if err != nil {
		panic(err)
	}
	defer Store.Close()

	r := mux.NewRouter()
	r.Use(handle.LoggingMiddleware)

	httpAddr := fmt.Sprintf("%s:%s",
		Env("HTTP_HOST", "0.0.0.0"),
		Env("PORT", "8080"))
	log.Println("listen", httpAddr)

	server := http.Server{Addr: httpAddr, Handler: r}

	// home
	r.HandleFunc("/",
			handle.HomeSearch(Store, Templates, Options))

	// manager
	sub := r.PathPrefix("/admin").Subrouter()
	mid := auth.NewAuthMid(Store, "admin")
	sub.Use(mid.Middleware)
	sub.HandleFunc("/shutdown", manager.Shut()).Methods("GET")
	sub.HandleFunc("/shutdown", manager.Cancel()).Methods("PUT")
	sub.HandleFunc("/shutdown", manager.Down(&server)).Methods("POST")

	// User
	subUser := r.PathPrefix("/user").Subrouter()
	midUser := auth.NewAuthMid(Store, "user", "owner")
	subUser.Use(midUser.Middleware)

	// config
	sub.HandleFunc("/config.yml",
		manager.FileEdit(Env("CONFIG", "config.yml"), *Options)).Methods("GET")
	sub.HandleFunc("/config.yml",
		manager.FileUpdate(Env("CONFIG", "config.yml"), *Options)).Methods("POST")

	manager.Editor(sub, "css", Options)
	manager.Editor(sub, "js", Options)
	manager.Imager(sub, "images", Options)
	manager.Templater(sub, "templates", Options)
	manager.HR(Store, sub, "users", Options)
	manager.Bucketer(Store, sub, "buckets", Options)

	// page
	r.HandleFunc("/{filename}.html",

			handle.Page(Templates, *Options))

	// wiki page
	r.HandleFunc("/{key}",

			handle.WikiPage(Templates, Store, Options))

	static := Env("STATIC", "static")
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(static))))


	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
