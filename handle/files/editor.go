package files

import (
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func PathEditor(localPath string, webPath string, o cms.Options) http.HandlerFunc {
	if t == nil {
		Init()
	}
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
		data, err := ioutil.ReadFile(localPath+"/"+ filename)
		if onErr(w, err) {
			return
		}

		ext := path.Ext(filename)
		err = t.ExecuteTemplate(w, "editor", o.Extend(
			cms.Options{
				"SaveURL":  webPath+"/"+ filename,
				"BasePath": "https://pagecdn.io/lib/ace/1.4.12",
				"Theme":    "ace/theme/tomorrow",
				"Mode":     "ace/mode/"+ext[1:],
				"Data":     string(data),
			}))
	}
}

func PathUpdate(localPath string, webPath string, o cms.Options) http.HandlerFunc {
	if t == nil {
		Init()
	}
	h := PathEditor(localPath, webPath, o)
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
		data, err := io.ReadAll(r.Body)
		if onErr(w, err) {
			return
		}
		defer r.Body.Close()

		f, err := os.OpenFile(localPath+"/"+filename, os.O_WRONLY|os.O_CREATE, 0666)
		if onErr(w, err) {
			return
		}
		_, err = f.Write(data)
		if onErr(w, err) {
			return
		}
		h(w, r)
	}
}