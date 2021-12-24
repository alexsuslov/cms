package handle

import (
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/store"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Login(s *store.Store, T ITemplate, o *cms.Options) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		err := T.ExecuteTemplate(w, "home", o)
		if err != nil {
			logrus.Error(err)
		}
		return

	}
}

func InviteGET(s *store.Store, T ITemplate, o *cms.Options) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		err := T.ExecuteTemplate(w, "invite_form", o)
		if err != nil {
			logrus.Error(err)
		}
		return

	}
}

func InvitePOST(s *store.Store, T ITemplate, o *cms.Options) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		err := T.ExecuteTemplate(w, "home", o)
		if err != nil {
			logrus.Error(err)
		}
		return

	}
}
