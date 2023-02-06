//
// (c) 2022 Alex Suslov
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is furnished to do
// so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package handle

import (
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/model"
	"github.com/alexsuslov/cms/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
	"gopkg.in/yaml.v3"
	"net/http"
	"reflect"
)

type kv struct {
	key   string
	value interface{}
}

func Items(t ITemplate, s *store.Store, o cms.Options) http.HandlerFunc {

	onErr := Err(t, o)
	dict := model.Dict{}
	iter := reflect.ValueOf(o["Dict"]).MapRange()

	for iter.Next() {
		key := iter.Key().Interface().(string)
		value := iter.Value().Interface().(string)
		dict[key] = value
	}

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		q := r.URL.Query()
		tag := q.Get("tag")
		bucketname, ok := params["bucket"]
		if !ok {
			onErr(w, fmt.Errorf("404"))
		}
		template := bucketname

		opt := store.NewSelectOptions().FromQuery(r.URL.Query())
		keyValues, err := store.Select(s, []byte(bucketname))(*opt)
		if onErr(w, err) {
			return
		}

		items := map[string]model.Item{}
		rows := map[int]map[string]model.Item{}
		rowItems := 4
		i := 0
		row := 0
		for k, v := range keyValues {
			item := model.Item{}
			err := yaml.Unmarshal(v, &item)
			if err != nil {
				logrus.WithField("data", string(v)).
					Error(err)
			}
			if !item.IsTag(tag) {
				continue
			}

			if !item.Disable {
				items[k] = item
				if rows[row] == nil {
					rows[row] = map[string]model.Item{}
				}
				rows[row][k] = item
				i++
				if i == rowItems {
					i = 0
					row++
				}
			}
		}

		keyValues = nil

		if q.Get("card") == "0" {
			template = "wheels_table"
		}

		err = t.ExecuteTemplate(w, template, o.Extend(
			cms.Options{
				"Name":  dict[bucketname],
				"URL":   fmt.Sprintf("/%s", bucketname),
				"Items": items,
				"Rows":  rows,
			}))
		if err != nil {
			logrus.Error(err)
		}

	}
}

func Item(t ITemplate, s *store.Store, o cms.Options) http.HandlerFunc {
	onErr := Err(t, o)

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		bucketName, ok := params["bucket"]
		if !ok {
			err := fmt.Errorf("404")
			if onErr(w, err) {
				return
			}
		}
		Key, ok := params["key"]
		if !ok {
			err := fmt.Errorf("404")
			if onErr(w, err) {
				return
			}
		}
		item := model.Item{}
		err := s.DB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketName))
			if b == nil {
				return nil
			}
			return yaml.Unmarshal(b.Get([]byte(Key)), &item)
		})

		if onErr(w, err) {
			return
		}

		err = t.ExecuteTemplate(w, "wheel", o.Extend(
			cms.Options{
				"Key":     Key,
				"EditURL": fmt.Sprintf("/admin/buckets/%s/%s", bucketName, Key),
				"Item":    item,
			}))
		if err != nil {
			logrus.Error(err)
		}

	}
}
