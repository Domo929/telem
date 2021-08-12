package livetiming

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Sessions contains all possible sessions in a two-depth map. The map is map[year]map[round]Session
type Sessions map[int64]map[int64]Session

type Session struct {
	Season int64
	Round  int64
	Name   string
	Date   time.Time
	Path   string
}

func (s *Session) String() string {
	return fmt.Sprintf("%d_%d_%s", s.Season, s.Round, s.Name)
}

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
