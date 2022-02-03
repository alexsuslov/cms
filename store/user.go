package store

import (
	"encoding/json"
	"fmt"
	"github.com/alexsuslov/cms/model"
	bolt "go.etcd.io/bbolt"
)

var USERS = []byte("_users")

type IUser interface {
	ToBytes() []byte
	GetName() string
	SetCreated()
	SetUpdated()
}

func (s Store) CreateUser(user IUser) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(USERS)
		if err != nil {
			return err
		}
		data := b.Get([]byte(user.GetName()))
		if data != nil {
			return fmt.Errorf("user exists")
		}
		user.SetCreated()
		user.SetUpdated()
		return b.Put([]byte(user.GetName()), user.ToBytes())
	})
}

func (s Store) GetUser(username string, user interface{}) (err error) {
	return s.DB.View(func(tx *bolt.Tx) error {
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
}

func (s Store) FirstUser() error {
	err := s.CreateUser(
		model.NewAdminUser())
	if err != nil && err.Error() == "user exists" {
		err = nil
	}
	return err
}
