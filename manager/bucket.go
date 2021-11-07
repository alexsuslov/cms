package manager

import (
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/alexsuslov/cms/model"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Bucket(s *model.Store, path string, o cms.Options) http.HandlerFunc {

	Init()
	onErr := handle.Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		backetname, ok := params["backetname"]
		if !ok {
			onErr(w, fmt.Errorf("404"))
		}

		// todo: replace universal func with filter, limit, offset
		err := s.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(backetname))
			if b == nil {
				return fmt.Errorf("404")
			}
			items := map[string]string{}
			err := b.ForEach(func(k, v []byte) error {
				K := string(k)
				V := string(v)
				items[K] = V
				return nil
			})

			if err != nil {
				return err
			}

			err = t.ExecuteTemplate(w, "bucket", o.Extend(
				cms.Options{
					"URL":   path + "/" + backetname,
					"Items": items,
				}))
			if err != nil {
				logrus.Error(err)
			}

			return nil

		})
		onErr(w, err)
	}
}
