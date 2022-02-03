package store

import (
	"encoding/json"
	"fmt"
	"github.com/alexsuslov/cms/mail"
	"github.com/alexsuslov/cms/model"
	uuid "github.com/satori/go.uuid"
	bolt "go.etcd.io/bbolt"
	"gopkg.in/gomail.v2"
	"os"
	"time"
)

var INVITES = []byte("_invites")

func (s *Store) InviteToUser(key string, user IUser) error {
	// get invite by id
	invite, err := s.GetInvite(key)
	if err != nil {
		return err
	}
	// create user
	err = s.CreateUser(user)
	if err != nil {
		return err
	}
	// delete invite
	return s.DeleteInvite(invite)
}

type IInvite interface {
	GetUUID() []byte
	GetEmail() string
	IsExpired() bool
}

func (s *Store) DeleteInvite(invite IInvite) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(INVITES)
		if err != nil {
			return err
		}

		err = b.Delete(invite.GetUUID())
		if err != nil {
			return err
		}
		return b.Delete([]byte(invite.GetEmail()))
	})
}

func (s *Store) GetInvite(key string) (*model.Invite, error) {
	invite := &model.Invite{}
	return invite, s.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(INVITES)
		if err != nil {
			return err
		}

		data := b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("401:no invite")
		}

		err = json.Unmarshal(data, invite)
		if err != nil {
			b.Delete([]byte(key))
			return err
		}
		// check expire
		if invite.IsExpired() {
			b.Delete([]byte(key))
			b.Delete([]byte(invite.Email))
			return fmt.Errorf("401:invite expired")
		}
		return nil
	})
}

func (s *Store) CreateOrUpdateInvite(email string) (invite *model.Invite, err error) {
	invite = &model.Invite{
		UUID:    uuid.NewV4(),
		Email:   email,
		Expired: time.Now().Add(8 * time.Hour),
		Count:   3,
	}
	return invite, s.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(INVITES)
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
