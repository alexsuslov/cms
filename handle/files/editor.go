package files

import (
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

var modes = map[string]string{
	".css":  "ace/mode/css",
	".js":   "ace/mode/javascript",
	".json": "ace/mode/json",
	".md":   "ace/mode/markdown",
	".tmpl": "ace/mode/html",
}

func PathEdit(localPath string, webPath string, o cms.Options) http.HandlerFunc {

	Init()
	onErr := handle.Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		filename, ok := params["filename"]
		if !ok {
			err := fmt.Errorf("404")
			if onErr(w, err) {
				return
			}
		}
		filename = path.Base(filename)
		data, err := ioutil.ReadFile(localPath + "/" + filename)
		if err != nil {
			logrus.Warning(err)
			data = []byte("")
		}

		mode, ok := modes[path.Ext(filename)]
		if !ok {
			mode = "ace/mode/text"
		}

		err = t.ExecuteTemplate(w, "editor", o.Extend(
			cms.Options{
				"SaveURL":  webPath + "/" + filename,
				"BasePath": "https://pagecdn.io/lib/ace/1.4.12",
				"Theme":    "ace/theme/tomorrow",
				"Mode":     mode,
				"Data":     string(data),
			}))
	}
}

func PathUpdate(localPath string, webPath string, o cms.Options) http.HandlerFunc {

	Init()
	onErr := handle.Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		filename, ok := params["filename"]
		if !ok {
			err := fmt.Errorf("404")
			if onErr(w, err) {
				return
			}
		}
		filename = path.Base(filename)

		data, err := io.ReadAll(r.Body)
		if onErr(w, err) {
			return
		}
		defer r.Body.Close()

		f, err := os.OpenFile(localPath+"/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
		if onErr(w, err) {
			return
		}
		defer f.Close()
		_, err = f.Write(data)
		if onErr(w, err) {
			return
		}
	}
}
