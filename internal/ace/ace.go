package ace

import (
	"embed"
	"github.com/alexsuslov/cms"
	"html/template"
	"net/http"
)

//go:embed templates/*
var templates embed.FS

type IFiler interface {
	GetData(*http.Request) ([]byte, error)
}

type IError interface {
	IsError(w http.ResponseWriter, err error) bool
}

func Editor(o cms.Options, f IFiler, e IError) func(http.ResponseWriter, *http.Request) {
	t, err := template.ParseFS(templates, "templates/*")
	if err != nil {
        panic(err)
    }

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := f.GetData(r)
		if e.IsError(w, err) {
			return
		}
		err = t.ExecuteTemplate(w, "editor", o.Extend(
			cms.Options{
				"BasePath": "https://pagecdn.io/lib/ace/1.4.12",
				"Theme":    "ace/theme/tomorrow",
				"Mode":     "ace/mode/javascript",
				"Data":     data,
			}))
	}
}
