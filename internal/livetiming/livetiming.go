package livetiming

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

var (
	once     = new(sync.Once)
	sessions = make(Sessions)
)

func Init() error {
	var err error

	once.Do(func() {
		var f *os.File
		f, err = os.Open("paths.json")
		if err != nil {
			return
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Println("issue closing paths.json file: ", err)
			}
		}(f)

		sessSlice := make([]Session, 0)
		if err = json.NewDecoder(f).Decode(&sessSlice); err != nil {
			return
		}

		for _, sess := range sessSlice {
			if _, ok := sessions[sess.Season]; !ok {
				sessions[sess.Season] = make(map[int64]Session)
			}
			sessions[sess.Season][sess.Round] = sess
		}
	})

	return err
}

func Info(year, round int64) (*Session, error) {
	var (
		sess Session
		ok   bool
	)
	if _, ok = sessions[year]; !ok {
		return nil, errors.New("year is not within a valid range")
	}
	if sess, ok = sessions[year][round]; !ok {
		return nil, errors.New("round is not within a valid range for that year")
	}

	return &sess, nil
}
