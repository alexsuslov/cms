package list

import (
	"embed"
	"github.com/alexsuslov/cms"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

//go:embed templates/*
var templates embed.FS

type ITemplate interface {
	ExecuteTemplate(wr io.Writer, data interface{}) error
}

type IError interface {
	IsError(w http.ResponseWriter, err error) bool
}

func Row(o cms.Options, t ITemplate, e IError) func(http.ResponseWriter, interface{}) {
	return func(w http.ResponseWriter, item interface{}) {
		err := t.ExecuteTemplate(w, o.Extend(
			cms.Options{
				"Item": item,
			}))
		if err != nil {
			logrus.Error(err)
		}
	}
}
