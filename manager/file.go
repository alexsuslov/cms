package manager

import (
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/alexsuslov/cms/vali"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func FileEdit(file string, o cms.Options) http.HandlerFunc {

	Init()
	onErr := handle.Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(string(file))
		if onErr(w, err) {
			return
		}
		err = t.ExecuteTemplate(w, "editor", o.Extend(
			cms.Options{
				"SaveURL":  "/admin/config.yml",
				"BasePath": "https://cdnjs.cloudflare.com/ajax/libs/ace/1.15.0",
				"Theme":    "ace/theme/tomorrow",
				"Mode":     "ace/mode/yaml",
				"Data":     string(data),
			}))
	}
}

func FileUpdate(file string, o cms.Options) http.HandlerFunc {

	Init()
	h := FileEdit(file, o)
	onErr := handle.Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if onErr(w, err) {
			return
		}
		defer r.Body.Close()

		if onErr(w, writeFile(file, data)) {
			return
		}

		o.Refresh(data)
		h(w, r)
	}
}

func writeFile(filename string, data []byte) error {
	// validate by file type
	ext := path.Ext(filename)
	if err := vali.IsValid(ext, data); err != nil {
		return err
	}
	//write file
	return os.WriteFile(filename, data, 0644)
}
