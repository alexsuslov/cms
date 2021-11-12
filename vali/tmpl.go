package vali

import "html/template"

func Tmpl(data []byte)(err error){
	_, err = template.New("validate").Parse(string(data))
	return
}