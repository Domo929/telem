package handlers

import (
	"log"
	"net/http"
)

func Health(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Println("issue writing ping output: ", err.Error())
	}
}
