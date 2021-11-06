package files

import (
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type FileInfo struct {
	Name string
	Size int64
	Created string
}


type Editables []string

func (E Editables) Is(ext string) (result bool) {
	for _, v := range E {
		if ext == "."+v {
			return true
		}
	}
	return
}

var ru = "2006-01-02T15:04:05"


var Editable = Editables{
	"tmpl",
	"js",
	"css",
}

func Files(localPath string, path string, o cms.Options)http.HandlerFunc{
	if t == nil {
		Init()
	}
	onErr := handle.Err(t, o)
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// rm file
		if f, ok := query["rm"]; ok {
			pathFile := fmt.Sprintf("%s/%s", localPath, f[0])
			logrus.Info(pathFile)
			err := os.Remove(pathFile)
			if err != nil {
				logrus.Error(err)
			}
		}

		files, err := ioutil.ReadDir(localPath)
		if onErr(w, err) {
			return
		}
		var Files []FileInfo
		for _, f := range files {
			// skip dot files
			if strings.HasPrefix(f.Name(), "."){
				continue
			}

			Files = append(Files, FileInfo{
				f.Name(),
				f.Size(),
				f.ModTime().Format(ru),
			})
		}
		err = t.ExecuteTemplate(w, "files", o.Extend(
			cms.Options{
				"URL": path,
				"Files": Files,

			}))
		if err!= nil{
			logrus.Error(err)
		}
	}
}


func FileUpload(localPath string, path string, o cms.Options)http.HandlerFunc{
	onErr := handle.Err(t, o)
	h:=Files(localPath, path, o)
	return func(w http.ResponseWriter, r *http.Request) {
		if file, handler, err := r.FormFile("file"); err == nil {
			defer file.Close()

			filePath := fmt.Sprintf("%s/%s", localPath, handler.Filename)
			f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
			if err == nil {
				_, err = io.Copy(f, file)
				if onErr(w, err) {
					return
				}
			} else {
				onErr(w, err)
				return
			}
		}
		h(w, r)
	}
}