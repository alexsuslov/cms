package model

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
)

var USERS = []byte("_users")

func NewStoreBDB(filename string) (store *Store, err error) {
	DB, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return
	}
	store = &Store{
		filename,
		DB,
	}
	return store, store.FirstUser()
}

type Store struct {
	Filename string
	DB       *bolt.DB
}

func (s Store) GetUser(username string) (user *User, err error) {
	user = &User{}
	err = s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(USERS)
		if b == nil {
			return fmt.Errorf("no bucket")
		}
		data := b.Get([]byte(username))
		if data == nil {
			return fmt.Errorf("no user")
		}
		return json.Unmarshal(data, user)
	})
	return user, err
}

func (s Store) FirstUser() error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(USERS)
		if err != nil {
			return err
		}
		username := Env("ADMIN_USER", "root")
		pass := Env("ADMIN_USER_PASS", "123456")
		data := b.Get([]byte(username))
		if data == nil {
			u := User{
				Username: "admin",
				Roles:    []string{"admin"},
			}
			u.Token = u.GenToken(pass)
			logrus.
				WithField("user", username).
				Info("user created")
			return b.Put([]byte(username), u.ToBytes())
		}
		return nil
	})
}
