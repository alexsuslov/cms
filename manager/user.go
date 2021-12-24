package manager

import (
	"encoding/json"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/alexsuslov/cms/model"
	"github.com/alexsuslov/cms/store"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var INVITE_URL = "/invite"

func Users(s *store.Store, path string, o cms.Options) http.HandlerFunc {

	Init()
	onErr := handle.Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		if email := query.Get("email"); email != "" {
			invite, err := s.CreateOrUpdateInvite(email)
			if err != nil {
				logrus.WithField("email", email).Errorf("inviter:%v", err)
			} else {
				if err := store.SendInvite(os.Getenv("HOST")+INVITE_URL, invite); err != nil {
					logrus.Warningf("send invite:%v", err)
				}
			}
		}

		opt := store.NewSelectOptions().FromQuery(r.URL.Query())

		keyvalues, err := store.Select(s, store.USERS)(*opt)
		if onErr(w, err) {
			return
		}

		items := map[string]model.User{}
		for k, v := range keyvalues {
			user := model.User{}
			err := json.Unmarshal(v, &user)
			if onErr(w, err) {
				return
			}
			items[k] = user
		}

		keyvalues = nil

		err = t.ExecuteTemplate(w, "users", o.Extend(
			cms.Options{
				"URL":   path,
				"Items": items,
			}))
		if err != nil {
			logrus.Error(err)
		}

	}
}
