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
	memCache  = flag.Bool("memCache", false, "whether to cache races in memory vs on disk")
	cachePath = flag.String("cache", "cache", "the path to the folder to use for a cache path")
)

func main() {
	flag.Parse()

	var (
		c   cache.Cache
		err error
	)

	if *memCache {
		c = cache.NewMemoryCache()
	} else {
		c, err = cache.NewFileCache(*cachePath)
	}
	if err != nil {
		log.Fatalln("issue setting up cache: ", err)
	}
	cache.SetCache(c)

	if err = livetiming.Init(); err != nil {
		log.Fatalln("error starting live timing: ", err)
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
