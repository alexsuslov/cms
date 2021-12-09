package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/alexsuslov/cms/store"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
	"net/http"
	"sort"
)

var VISITS = []byte("visits")

type Visit struct {
	URL   string
	Count int
}

func Visits(s *store.Store, path string, o cms.Options) http.HandlerFunc {

	Init()
	onErr := handle.Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()

		prefix := query.Get("prefix")

		visits := []Visit{}
		err := s.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(VISITS)
			if b == nil {
				return fmt.Errorf("404")
			}
			return b.ForEach(func(k, v []byte) error {
				var Count int
				err := json.Unmarshal(v, &Count)
				if err != nil {
					return nil
				}

				if prefix != "" && !bytes.HasPrefix(k, []byte(prefix)) {
					return nil
				}

				visits = append(visits, Visit{
					string(k), Count,
				})
				return nil
			})
		})

		if onErr(w, err) {
			return
		}

		sort.Slice(visits, func(i, j int) bool {
			return visits[i].Count > visits[j].Count
		})

		err = t.ExecuteTemplate(w, "visits", o.Extend(
			cms.Options{
				"URL":    path + "/",
				"Items":  visits,
				"Prefix": r.URL.Query().Get("prefix"),
				"Value":  r.URL.Query().Get("value"),
			}))
		if err != nil {
			logrus.Error(err)
		}

	}
}
