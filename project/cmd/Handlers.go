package main

import (
	_ "encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (aps *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", aps.home)
	mux.HandleFunc("/postform", aps.postform)
	mux.HandleFunc("/quotesform", aps.quotesform)
	mux.HandleFunc("/quotesform/create", aps.createQuote)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

func (aps *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	aps.render(w, r, "./ui/html/home.page.tmpl", nil)
}

func (aps *application) postform(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error! Locked.", 423)
		return
	}
	name := r.FormValue("username")
	var dat templateData = name
	aps.render(w, r, "./ui/html/postform.page.tmpl", dat)
}

//выводим заметки по айди
func (aps *application) quotesform(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("id") != "" {
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		qu, err := aps.data.GetQuote(id)
		if err != nil {
			return
		}
		var dat templateData = qu
		aps.render(w, r, "./ui/html/showquotes.page.tmpl", dat)

	} else if r.Method == http.MethodPost {
		text := r.FormValue("text")
		title := r.FormValue("title")
		id := aps.data.AddQuote(title, text)
		http.Redirect(w, r, fmt.Sprintf("/quotesform?id=%d", id), http.StatusSeeOther)

	} else {
		http.Error(w, "Error! Locked.", 423)
		return
	}
}

func (aps *application) createQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Error! Locked.", 423)
		return
	}
	aps.render(w, r, "./ui/html/quotes.page.tmpl", nil)
}
