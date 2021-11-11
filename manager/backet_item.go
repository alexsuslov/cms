package manager

import (
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/alexsuslov/cms/store"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func BucketItem(s *store.Store, path string, o cms.Options) http.HandlerFunc {

	Init()
	onErr := handle.Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		bucketName, ok := params["bucket"]
		if !ok {
			err := fmt.Errorf("404")
			if onErr(w, err) {
				return
			}
		}
		Key, ok := params["item"]
		if !ok {
			err := fmt.Errorf("404")
			if onErr(w, err) {
				return
			}
		}
		value := ""
		s.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return nil
			}
			Value := b.Get([]byte(Key))
			if Value != nil {
				value = string(Value)
			}
			return nil
		})
		mode := "ace/mode/json"
		err := t.ExecuteTemplate(w, "editor", o.Extend(
			cms.Options{
				"SaveURL":  fmt.Sprintf("%s/%s/%s", path, bucketName, Key),
				"BasePath": "https://pagecdn.io/lib/ace/1.4.12",
				"Theme":    "ace/theme/tomorrow",
				"Mode":     mode,
				"Data":     value,
			}))
		if err != nil {
			logrus.Error(err)
		}

	}
}

func BucketItemUpdate(s *store.Store, path string, o cms.Options) http.HandlerFunc {

	Init()
	onErr := handle.Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		bucketName, ok := params["bucket"]
		if !ok {
			err := fmt.Errorf("404")
			if onErr(w, err) {
				return
			}
		}
		Key, ok := params["item"]
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

		err = s.DB.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
			if err != nil {
				return err
			}
			return b.Put([]byte(Key), data)
		})
		onErr(w, err)
	}
}
