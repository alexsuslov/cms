package manager

import (
	"fmt"
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/store"
	"github.com/gorilla/mux"
)

func Bucketer(store *store.Store, sub *mux.Router, ext string, Options *cms.Options) {

	p := "/" + ext
	w := "/admin/" + ext
	sub.HandleFunc(p,
		Buckets(store, w, *Options)).
		Methods("GET")

	p = fmt.Sprintf("/%s/{bucket}", ext)

	sub.HandleFunc(p,
		Bucket(store, w, *Options)).
		Methods("GET")

	p = fmt.Sprintf("/%s/{bucket}/{item}", ext)

	sub.HandleFunc(p,
		BucketItem(store, w, *Options)).
		Methods("GET")

	sub.HandleFunc(p,
		BucketItemUpdate(store, w, *Options)).Methods("POST")
}


func HR(store *store.Store, sub *mux.Router, ext string, Options *cms.Options) {

	p := "/" + ext
	w := "/admin/" + ext

	sub.HandleFunc(p,
		Users(store, w, *Options)).
		Methods("GET")
	//
	//p = fmt.Sprintf("/%s/{bucket}/{item}", ext)
	//
	//sub.HandleFunc(p,
	//	BucketItem(store, w, *Options)).
	//	Methods("GET")
	//
	//sub.HandleFunc(p,
	//	BucketItemUpdate(store, w, *Options)).Methods("POST")
}