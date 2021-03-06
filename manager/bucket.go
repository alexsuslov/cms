package manager

import (
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/alexsuslov/cms/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Bucket(s *store.Store, path string, o cms.Options) http.HandlerFunc {

	Init()
	onErr := handle.Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		bucketname, ok := params["bucket"]
		if !ok {
			onErr(w, fmt.Errorf("404"))
		}

		// rm bucket
		rm := r.URL.Query().Get("rm")
		if rm != "" {
			s.RmBucketItem([]byte(bucketname), []byte(rm))
		}

		opt := store.NewSelectOptions().FromQuery(r.URL.Query())

		keyvalues, err := store.Select(s, []byte(bucketname))(*opt)
		if onErr(w, err) {
			return
		}

		items := map[string]string{}
		for k, v := range keyvalues {
			items[k] = string(v)
		}

		keyvalues = nil

		err = t.ExecuteTemplate(w, "bucket", o.Extend(
			cms.Options{
				"URL":    path + "/" + bucketname,
				"Items":  items,
				"Prefix": r.URL.Query().Get("prefix"),
				"Value":  r.URL.Query().Get("value"),
			}))
		if err != nil {
			logrus.Error(err)
		}

	}
}
