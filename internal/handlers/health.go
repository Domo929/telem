package handlers

import (
	"log"
	"net/http"
)

// Health is the handler that simply returns "OK" to check for healthiness
func Health(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Println("issue writing ping output: ", err.Error())
	}
}
