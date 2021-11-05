package files

import "html/template"

var t *template.Template

func Init(){
	var err error
	t, err = template.ParseFS(templates, "templates/*")
    if err != nil {
        panic(err)
    }
}