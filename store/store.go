package store

import (
	"github.com/alexsuslov/cms/model"
	"github.com/blevesearch/bleve/v2"
	"github.com/boltdb/bolt"
	"log"
)

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
