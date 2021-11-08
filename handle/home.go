package handle

import (
	"github.com/alexsuslov/cms"
	"github.com/alexsuslov/cms/model"
	"github.com/blevesearch/bleve/v2"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
)

type ITemplate interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
}

func Home(S *model.Store, T ITemplate, o cms.Options) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q:= r.URL.Query()
		query := bleve.NewQueryStringQuery(q.Get("q"))
		searchRequest := bleve.NewSearchRequest(query)
		searchResult, err := S.Index.Search(searchRequest)
		if err!= nil{
			log.Println(err)
		}else{
			if searchResult.Total > 0{
				err = T.ExecuteTemplate(w, "search_result", o.Extend(cms.Options{
					"Hits": searchResult.Hits,
					"Total":searchResult.Total,
				}))
				if err != nil {
					logrus.Error(err)
				}
				return
			}
		}

		err = T.ExecuteTemplate(w, "home", o)
		if err != nil {
			logrus.Error(err)
		}
	}
}
