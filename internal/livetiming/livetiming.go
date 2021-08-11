package livetiming

import (
	"encoding/json"
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
		defer f.Close()

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
