package main

import (
	"log"
	"net/http"
)

/*func TakeAddr() (srv *http.Server, Rest *http.Server) {
	addr := 4000
	for {

		break
	}
	return srv, nil
}*/

func main() {

	addr1 := ":4000"
	addr2 := ":4001"
	srv := &http.Server{
		Addr:    addr1,
		Handler: routes(),
	}
	Rest := &http.Server{
		Addr:    addr2,
		Handler: restroutes(),
	}

	log.Println("Запуск сервера на http://127.0.0.1" + addr1)
	go Rest.ListenAndServe()
	srv.ListenAndServe()
}
