package main

import (
	_ "encoding/json"
	"net/http"
)

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/postform", postform)
	mux.HandleFunc("/test1", test1)
	mux.HandleFunc("/test2", test2)
	mux.HandleFunc("/test3", test3)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	render(w, r, "./ui/html/home.page.tmpl", nil)
}

func postform(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error! Locked.", 423)
		return
	}
	name := r.FormValue("username")
	var dat templateData = name
	render(w, r, "./ui/html/postform.page.tmpl", dat)
}

func test1(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			render(w, r, "./ui/html/test1.page.tmpl", nil)
		}
	default:
		{
			http.Error(w, "Error! Locked.", 423)
			return
		}
	}
}

func test2(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			render(w, r, "./ui/html/test2.page.tmpl", nil)
		}
	default:
		{
			http.Error(w, "Error! Locked.", 423)
			return
		}
	}
}

func test3(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			render(w, r, "./ui/html/test3.page.tmpl", nil)
		}
	default:
		{
			http.Error(w, "Error! Locked.", 423)
			return
		}
	}
} //нужно будет исправить
