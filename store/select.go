package store

import (
	"fmt"
	"github.com/boltdb/bolt"
	"net/url"
	"strconv"
)

type SelectOptions struct {
	Limit  *int
	Offset *int
	Prefix *string
	Reverse bool
}

func NewSelectOptions()*SelectOptions{
	return &SelectOptions{}
}

func (Opt *SelectOptions)FromQuery(values url.Values)*SelectOptions{
	if l, err := strconv.Atoi(values.Get("limit")); err == nil {
			Opt.Limit = &l
		}

		if f, err := strconv.Atoi(values.Get("offset")); err == nil {
			Opt.Offset = &f
		}

		if p := values.Get("prefix"); p != "" {
			Opt.Prefix = &p
		}
		return Opt
}


func (Opt *SelectOptions) IsLimit() (r bool) {
	if r = *Opt.Limit == 0; r {
		return
	}
	*Opt.Limit--
	return
}

func (Opt *SelectOptions) SetLimit(limit *int) *SelectOptions {
	if limit != nil {
		Opt.Limit = limit
	}
	return Opt
}

func (Opt *SelectOptions) IsOffset() (r bool) {
	if r = *Opt.Offset == 0; r {
		return
	}
	*Opt.Offset--
	return
}

func (Opt *SelectOptions) SetOffset(Offset *int) *SelectOptions {
	if Offset != nil {
		Opt.Offset = Offset
	}
	return Opt
}

func Select(s *Store, bucketName []byte) func(...SelectOptions) (map[string][]byte, error) {
	Limit := 100
	Offset := 0
	return func(ops ...SelectOptions) (result map[string][]byte, err error) {
		result = map[string][]byte{}

		option := SelectOptions{Limit: &Limit, Offset: &Offset}
		for _, op := range ops {
			option.
				SetLimit(op.Limit).
				SetOffset(op.Offset)
		}

		err = s.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucketName)
			if b == nil {
				return fmt.Errorf("no bucket")
			}

			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {

				if *option.Offset != 0 {
					*option.Offset--
					continue
				}

				if *option.Limit == 0 {
					break
				}
				result[string(k)] = v
			}
			return nil
		})
		return
	}
}

func Read(c bolt.Cursor, option SelectOptions) map[string][]byte {
	result := map[string][]byte{}
	first := func()([]byte, []byte){
		if option.Prefix!= nil{
			return c.Seek([]byte(*option.Prefix))
		}
		if option.Reverse{
			return c.Last()
		}
		return c.First()
	}

	getter := func()([]byte, []byte){
		if option.Reverse{
			return c.Prev()
		}
		return c.Next()
	}

	for k, v := first(); k != nil; k, v = getter() {
		if !option.IsOffset() {
			continue
		}
		if option.IsLimit() {
			break
		}

		result[string(k)] = v
	}
	return result
}


