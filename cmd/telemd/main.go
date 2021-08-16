package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/Domo929/telem.git/internal/cache"

	"github.com/Domo929/telem.git/internal/handlers"
	"github.com/Domo929/telem.git/internal/livetiming"
	"github.com/gorilla/mux"
)

var (
	addr      = flag.String("addr", "0.0.0.0:8080", "the local address of the daemon")
	cachePath = flag.String("cache", "cache", "the path to the folder to use for a cache path")
)

func main() {
	flag.Parse()

	localCache, err := cache.New(*cachePath)
	if err != nil {
		log.Fatalln("error setting up local cache : ", err)
	}
	cache.SetLocal(localCache)

	if err = livetiming.Init(); err != nil {
		log.Fatalln("error starting live timing :", err)
	}

	r := router()

	srv := http.Server{
		Addr:         *addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

func router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/ping", handlers.Health)

	s := r.PathPrefix("/telemetry/{year}/{round}").Subrouter()
	s.HandleFunc("/", handlers.SessionTelemetry)
	s.HandleFunc("/info", handlers.SessionInfo)
	s.HandleFunc("/raw", handlers.SessionRaw)

	return r
}
