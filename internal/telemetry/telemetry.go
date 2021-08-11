package telemetry

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Domo929/telem.git/internal/livetiming"
)

// SessionInfo returns information about a session, including Season, Round, Name, Date, and the API URL path
func SessionInfo(w http.ResponseWriter, r *http.Request) {
	vals := mux.Vars(r)

	year, round, err := parseSessionInfo(vals)
	if err != nil {
		if _, wErr := w.Write([]byte(err.Error())); wErr != nil {
			log.Println(wErr)
		}
		w.WriteHeader(http.StatusBadRequest)
	}

	sess, err := livetiming.Info(year, round)
	if err != nil {
		if _, wErr := w.Write([]byte(err.Error())); wErr != nil {
			log.Println(wErr)
		}
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(sess); err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusOK)
}

func parseSessionInfo(values map[string]string) (int64, int64, error) {
	yearStr, ok := values["year"]
	if !ok {
		return 0, 0, errors.New("year not set")
	}

	year, err := strconv.ParseInt(yearStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	roundStr, ok := values["round"]
	if !ok {
		return 0, 0, errors.New("round not set")
	}

	round, err := strconv.ParseInt(roundStr, 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return year, round, nil
}