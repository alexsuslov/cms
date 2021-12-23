package manager

import (
	"encoding/json"
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/handle"
	"github.com/alexsuslov/cms/mail"
	"github.com/alexsuslov/cms/model"
	"github.com/alexsuslov/cms/store"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"net/http"
	"os"
	"time"
)

var INVITE_URL = "/invite"

func Users(s *store.Store, path string, o cms.Options) http.HandlerFunc {

	Init()
	onErr := handle.Err(t, o)
	inviter := CreateOrUpdateInvite(s)

	return func(w http.ResponseWriter, r *http.Request) {

		query := r.URL.Query()
		if email := query.Get("email"); email != "" {
			invite, err := inviter(email)
			if err != nil {
				logrus.WithField("email", email).Errorf("inviter:%v", err)
			} else {
				if err := SendInvite(os.Getenv("HOST")+INVITE_URL, invite); err != nil {
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
func SendInvite(url string, invite *model.Invite) error {
	d, err := mail.GetDialer()
	if err != nil {
		return err
	}
	url += "/" + invite.UUID.String()
	m := gomail.NewMessage()
	from := os.Getenv("INVITE_FROM")
	if from == "" {
		return fmt.Errorf("no from email")
	}
	m.SetHeader("From", from)
	m.SetHeader("To", invite.Email)
	m.SetHeader("Subject", "invite")
	m.SetBody("text/html", fmt.Sprintf("<a href='%s'>вход</a>", url))
	return d.DialAndSend(m)
}

func CreateOrUpdateInvite(s *store.Store) func(email string) (*model.Invite, error) {
	return func(email string) (invite *model.Invite, err error) {
		invite = &model.Invite{
			UUID:    uuid.NewV4(),
			Email:   email,
			Expired: time.Now().Add(8 * time.Hour),
			Count:   3,
		}
		return invite, s.DB.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists(store.INVITES)
			if err != nil {
				return err
			}
			key := []byte(invite.Email)
			if tmp := b.Get(key); tmp != nil {
				tmpInvite := &model.Invite{}
				err = json.Unmarshal(tmp, tmpInvite)
				if err == nil && invite.Expired.After(time.Now()) {
					invite = tmpInvite
					return nil
				}
			}
			if err = b.Put([]byte(invite.Email),
				invite.ToBytes()); err != nil {
				return err
			}
			return b.Put(invite.UUID.Bytes(), invite.ToBytes())
		})
	}
}
