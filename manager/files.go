package manager

import (
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/sirupsen/logrus"
	"gopkg.in/validator.v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

type FileInfo struct {
	Name    string
	Size    int64
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
	"md",
}

func addPath(localPath string, filenames []string) (paths []string) {
	paths = []string{}
	for _, filename := range filenames {
		f := path.Base(filename)
		paths = append(paths, fmt.Sprintf("%s/%s", localPath, f))
	}
	return
}

func rmFiles(filePaths []string) error {
	for _, filePath := range filePaths {
		logrus.WithField("filepath", filePath).Info("remove file")
		err := os.Remove(filePath)
		if err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}

func Files(localPath string, path string, o cms.Options) http.HandlerFunc {

	Init()
	onErr := handle.Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		// rm fill
		rm, ok := query["rm"]
		if ok {
			rmFiles(addPath(localPath, rm))
		}

		files, err := ioutil.ReadDir(localPath)
		if onErr(w, err) {
			return
		}

		var Files []FileInfo
		for _, f := range files {
			// skip dot files
			if strings.HasPrefix(f.Name(), ".") {
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
				"URL":   path,
				"Files": Files,
			}))
		if err != nil {
			logrus.Error(err)
		}
	}
}

type vUpload struct {
	Filename string `validate:"regexp=^[\w\-. ]+$"`
}

func FileUpload(localPath string, path string, o cms.Options) http.HandlerFunc {
	onErr := handle.Err(t, o)
	h := Files(localPath, path, o)

	return func(w http.ResponseWriter, r *http.Request) {
		if file, handler, err := r.FormFile("file"); err == nil {
			defer file.Close()
			q := vUpload{handler.Filename}
			err := validator.Validate(q)
			if onErr(w, err) {
				return
			}

			filePath := fmt.Sprintf("%s/%s", localPath, q.Filename)
			f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
			if err == nil {
				_, err = io.Copy(f, file)
				defer f.Close()
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
