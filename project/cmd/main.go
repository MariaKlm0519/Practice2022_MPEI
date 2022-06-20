package main

import (
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	data     *QuotesStore
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	aps := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		data:     NewStore(),
	}

	srv := &http.Server{
		Addr:    ":4000",
		Handler: aps.routes(),
	}

	infoLog.Println("Запуск сервера на http://127.0.0.1:4000")
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
