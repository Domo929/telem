package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/Domo929/telem.git/internal/handlers"
	"github.com/Domo929/telem.git/internal/livetiming"
	"github.com/Domo929/telem.git/internal/telemetry"
)

var addr = flag.String("addr", "0.0.0.0:8080", "the local address of the daemon")

func main() {
	flag.Parse()

	if err := livetiming.Init(); err != nil {
		log.Fatalln("error starting live timing :", err)
	}

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

	s := r.PathPrefix("/session/{year}/{round}").Subrouter()
	s.HandleFunc("/info", telemetry.SessionInfo)

	return r
}
