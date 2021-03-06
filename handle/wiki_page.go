package handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/store"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"html/template"
	"net/http"
)

var WIKI = []byte("wiki_pages")
var VALUES = []byte("wiki_values")

func WikiPage(t ITemplate, s *store.Store, o *cms.Options) func(w http.ResponseWriter, r *http.Request) {
	s.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(WIKI)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(VALUES)
		return err
	})

	onErr := Err(t, o)
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		key, keyOk := params["key"]
		if !keyOk {
			onErr(w, fmt.Errorf("404"))
			return
		}
		Page := ""
		Values := map[string]interface{}{}

		err := s.DB.View(func(tx *bolt.Tx) error {
			wiki := tx.Bucket(WIKI)
			values := tx.Bucket(VALUES)
			if wiki == nil || values == nil {
				return fmt.Errorf("302")
			}

			data := wiki.Get([]byte(key))
			if data == nil {
				return fmt.Errorf("302")
			}

			data1 := values.Get([]byte(key))
			if data1 != nil {
				err := json.Unmarshal(data1, &Values)
				if err != nil {
					return err
				}
				tmpT, err := template.New("").Parse(string(data))
				if err != nil {
					return err
				}
				var buf bytes.Buffer
				err = tmpT.Execute(&buf, Values)
				if err != nil {
					return err
				}
				data = buf.Bytes()

			}
			//indexer
			err := s.Index.Index(r.URL.String(), string(data))
			if err != nil {
				logrus.WithField("index.Index", r.URL.String()).Warning(err)
			}

			//log.Println("page indexed", key)

			output := markdown.ToHTML(data, nil, nil)

			Page = string(output)
			return nil
		})
		if err != nil {
			//	on 300 redirect to editor
			cookie, _ := r.Cookie("editor")
			if cookie != nil && cookie.Value == "Yes" && err.Error() == "302" {
				http.Redirect(w, r, "/admin/buckets/wiki_pages/"+key, 302)
				return
			} else {
				err = fmt.Errorf("404")
			}
			onErr(w, err)
			return
		}

		err = t.ExecuteTemplate(w, "wiki_page", o.Extend(cms.Options{
			"Values": Values,
			"HTML":   template.HTML(Page),
		}))
		if err != nil {
			logrus.Error(err)
		}
	}
}
