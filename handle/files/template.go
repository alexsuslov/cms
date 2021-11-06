package files

import (
	"embed"
	"html/template"
)


//go:embed templates/*
var templates embed.FS

var t *template.Template

func Init(){
	var err error
	t, err = template.ParseFS(templates, "templates/*")
    if err != nil {
        panic(err)
    }
}