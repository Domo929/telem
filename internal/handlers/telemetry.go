package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Domo929/telem.git/internal/livetiming"
	"github.com/Domo929/telem.git/internal/telemetry"
	"github.com/gorilla/mux"
)

// SessionTelemetry is the handler that returns the formatted/aggregated data for the race
func SessionTelemetry(w http.ResponseWriter, r *http.Request) {
	vals := mux.Vars(r)

	year, round, err := parseSessionInfo(vals)
	if err != nil {
		if _, wErr := w.Write([]byte(err.Error())); wErr != nil {
			log.Println(wErr)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sess, err := livetiming.GetSession(year, round)
	if err != nil {
		if _, wErr := w.Write([]byte(err.Error())); wErr != nil {
			log.Println(wErr)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	lines, err := telemetry.Load(sess)
	if err != nil {
		if _, wErr := w.Write([]byte(err.Error())); wErr != nil {
			log.Println(wErr)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(lines); err != nil {
		if _, wErr := w.Write([]byte(err.Error())); wErr != nil {
			log.Println(wErr)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
