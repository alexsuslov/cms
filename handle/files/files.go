package files

import (
	"embed"
	"github.com/alexsuslov/cms"
	"github.com/gorilla/mux"
)

//go:embed templates/*
var templates embed.FS

func Files(Router *mux.Router, o cms.Options){
	//t, err := template.ParseFS(templates, "templates/*")
    //if err != nil {
    //    panic(err)
    //}
}