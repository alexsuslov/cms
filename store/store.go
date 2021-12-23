package store

import (
	"encoding/json"
	"fmt"
	"github.com/alexsuslov/cms/model"
	"github.com/blevesearch/bleve/v2"
	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

var USERS = []byte("_users")
var INVITES = []byte("_invites")

func NewStoreBDB(filename string) (store *Store, err error) {
	search := model.Env("SEARCH", "search")
	index, err := bleve.Open(search)
	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(search, mapping)
		if err != nil {
			log.Println(err)
		}
	}

	DB, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return
	}
	store = &Store{
		filename,
		DB,
		index,
	}
	return store, store.FirstUser()
}

type Store struct {
	Filename string
	DB       *bolt.DB
	Index    bleve.Index
}

func (s Store) Close() error {
	err := s.DB.Close()
	if err != nil {
		return err
	}
	return s.Index.Close()
}

func (s Store) GetUser(username string) (user *model.User, err error) {
	user = &model.User{}
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
		username := model.Env("ADMIN_USER", "root")
		pass := model.Env("ADMIN_USER_PASS", "123456")
		data := b.Get([]byte(username))
		if data == nil || os.Getenv("ADMIN_USER_CREATE") == "YES" {
			u := model.User{
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

func (s Store) RmBucket(name []string) error {
	for _, bucket := range name {
		err := s.DB.Update(func(tx *bolt.Tx) error {
			return tx.DeleteBucket([]byte(bucket))
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Store) RmBucketItem(bucketName []byte, key []byte) error {

	return s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return nil
		}
		return b.Delete(key)
	})
}
