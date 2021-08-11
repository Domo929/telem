package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/Domo929/telem.git/internal/handlers"

	"github.com/gorilla/mux"
)

var addr = flag.String("addr", "0.0.0.0:8080", "the local address of the daemon")

func main() {
	flag.Parse()

	r := router()

	srv := http.Server{
		Addr:         *addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

func router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/ping", handlers.Health)

	return r
}
