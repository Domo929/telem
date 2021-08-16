package livetiming

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Sessions contains all possible sessions in a two-depth map. The map is map[year]map[round]Session
type Sessions map[int64]map[int64]Session

// Session contains all the information for a session
type Session struct {
	Season int64
	Round  int64
	Name   string
	Date   time.Time
	Path   string
}

// String returns the unique identifier for the Session
func (s *Session) String() string {
	name := strings.ReplaceAll(s.Name, " ", "_")
	return fmt.Sprintf("%d_%d_%s", s.Season, s.Round, name)
}

// UnmarshalJSON is the overload function to handle the unique format of the json stream
func (s *Session) UnmarshalJSON(b []byte) error {
	var (
		rawStr string
		ok     bool
		err    error
	)

	raw := make(map[string]interface{})

	if err = json.Unmarshal(b, &raw); err != nil {
		return err
	}

	if rawStr, ok = raw["season"].(string); !ok {
		return errors.New("could not unmarshal season field of session")
	}
	s.Season, err = strconv.ParseInt(rawStr, 10, 64)
	if err != nil {
		return err
	}

	if rawStr, ok = raw["round"].(string); !ok {
		return errors.New("could not unmarshal round field of session")
	}
	s.Round, err = strconv.ParseInt(rawStr, 10, 64)
	if err != nil {
		return err
	}

	if s.Name, ok = raw["raceName"].(string); !ok {
		return errors.New("could not unmarshal raceName field of session")
	}

	if s.Path, ok = raw["apiPath"].(string); !ok {
		return errors.New("could not unmarshal apiPath field of session")
	}
	if rawStr, ok = raw["date"].(string); !ok {
		return errors.New("could not unmarshal date field of session")
	}

	s.Date, err = time.Parse("2006-01-02", rawStr)
	if err != nil {
		return err
	}

	return nil
}
