package manager

import (
	"embed"
	"github.com/Masterminds/sprig"
	"html/template"
)

//go:embed templates/*
var templates embed.FS

var t *template.Template

func Init() {
	if t == nil {
		var err error
		t = template.New("base").
			Funcs(sprig.FuncMap())
		t, err = t.ParseFS(templates, "templates/*")
		if err != nil {
			panic(err)
		}
	}

}

func GetTemplate()*template.Template{
	return t
}
