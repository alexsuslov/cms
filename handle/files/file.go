package files

import (
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func FileEdit(file string, o cms.Options) http.HandlerFunc {
	if t == nil {
		Init()
	}
	onErr := handle.Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(string(file))
		if onErr(w, err) {
			return
		}
		err = t.ExecuteTemplate(w, "editor", o.Extend(
			cms.Options{
				"SaveURL":  "/admin/config",
				"BasePath": "https://pagecdn.io/lib/ace/1.4.12",
				"Theme":    "ace/theme/tomorrow",
				"Mode":     "ace/mode/yaml",
				"Data":     string(data),
			}))
	}
}

func FilePost(file string, o cms.Options) http.HandlerFunc {
	if t == nil {
		Init()
	}
	h := FileEdit(file, o)
	onErr := handle.Err(t, o)
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if onErr(w, err) {
			return
		}
		defer r.Body.Close()

		err = cms.Check(data)

		if onErr(w, err) {
			return
		}

		f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0666)
		if onErr(w, err) {
			return
		}
		_, err = f.Write(data)
		if onErr(w, err) {
			return
		}
		o.Refresh(data)
		h(w, r)
	}
}