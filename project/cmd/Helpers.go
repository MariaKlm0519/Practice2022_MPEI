package main

import (
	"html/template"
	"net/http"
)

type templateData interface{}

func (aps *application) render(w http.ResponseWriter, r *http.Request, name string, td templateData) {
	files := []string{
		name,
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	rs, err := template.ParseFiles(files...)
	if err != nil {
		aps.infoLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = rs.Execute(w, td)
	if err != nil {
		aps.infoLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
