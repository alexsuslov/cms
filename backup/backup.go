package backup

import (
	"archive/zip"
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

type Backup struct {
	Paths []string
}

func New(Paths interface{}) (*Backup) {
	paths, ok :=Paths.([]interface{})
	if !ok{
		logrus.Warning("no backup paths")
	}
	ps :=[]string{}
	for _, p := range paths{
		s, ok := p.(string)
		if !ok{
			continue
		}
 		ps= append(ps, s)
	}


	return &Backup{ps}
}

func (Backup Backup) Write(w http.ResponseWriter) error {
	zw := zip.NewWriter(w)
	for _, Path := range Backup.Paths {
		err := addFiles(zw, Path, "backup")
		if err != nil {
			logrus.WithField("Path", Path).Error(err)
			return err
		}
	}

	return zw.Close()
}

func addFiles(w *zip.Writer, basePath, baseInZip string) (err error) {
	file, err := os.Open(basePath)
	if err != nil {
		return err
	}
	defer file.Close()
	Info, err := file.Stat()
	if err != nil {
		return err
	}
	if !Info.IsDir() {
		return addFile(w, basePath, baseInZip)
	}

	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() {
			if strings.HasPrefix(file.Name(), ".") {
				continue
			}

			err = addFile(w, basePath + "/" + file.Name(), baseInZip)
			if err != nil {
				return err
			}

		} else if file.IsDir() {
			newBase := strings.Join([]string{basePath, file.Name()}, "/")
			err = addFiles(w, newBase, baseInZip)
			if err != nil {
				return
			}
		}
	}
	return
}

func addFile(w *zip.Writer, filePath, baseInZip string) (err error) {
	source, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer source.Close()

	newBaseInZip := baseInZip + "/" + filePath
	f, err := w.Create(newBaseInZip)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, source)
	return err
}

type ITemplate interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
}

func HandleFunc(T ITemplate, o *cms.Options) func(http.ResponseWriter, *http.Request) {
	backup := New(o.Get("Backup"))
	onErr := handle.Err(T, o)
	return func(w http.ResponseWriter, r *http.Request) {
		host:=r.Host
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition",
			fmt.Sprintf("attachment; filename=\"%s_%s.zip\"", host, "backup"))
		err := backup.Write(w)
		if onErr(w, err) {
			return
		}
	}
}
