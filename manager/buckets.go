package manager

import (
	"bytes"
	"encoding/json"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/alexsuslov/cms/model"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type ITemplate interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
}

type Item struct {
	Name    string
	Counter int
}

func (Item Item) ToBytes() []byte {
	data, _ := json.Marshal(Item)
	return data
}

func Buckets(s *model.Store, path string, o cms.Options) http.HandlerFunc {

	Init()
	onErr := handle.Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {
		items := []Item{}
		query := r.URL.Query()

		// rm bucket
		rm, ok := query["rm"]
		if ok {
			for _, bucket := range rm {
				s.DB.Update(func(tx *bolt.Tx) error {
					return tx.DeleteBucket([]byte(bucket))
				})
			}
		}
		// todo: replace universal func with filter, limit, offset
		err := s.DB.View(func(tx *bolt.Tx) error {
			err := tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
				b := tx.Bucket(name)
				stats := b.Stats()
				if !bytes.HasPrefix(name, []byte("_")) {
					items = append(items, Item{
						string(name),
						stats.KeyN,
					})

				}

				return nil
			})
			if err != nil {
				return err
			}

			err = t.ExecuteTemplate(w, "buckets", o.Extend(
				cms.Options{
					"URL":     path,
					"Buckets": items,
				}))
			if err != nil {
				logrus.Error(err)
			}

			return nil

		})
		onErr(w, err)
	}
}
