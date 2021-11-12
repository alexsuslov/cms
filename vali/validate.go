package vali

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

func IsValid(ext string, data []byte) (err error) {
	tmp := map[string]interface{}{}
	switch ext {
	case ".tmpl":
		err= Tmpl(data)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &tmp)
	case ".json":
		err= json.Unmarshal(data, &tmp)
	}
	if err!= nil{
		err = fmt.Errorf("500:%v", err)
	}
	return
}
