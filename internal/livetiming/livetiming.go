package livetiming

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	once     = new(sync.Once)
	sessions = make(Sessions)

	client = http.Client{Timeout: 10 * time.Second}
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

func GetSession(year, round int64) (*Session, error) {
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

func GetData(ctx context.Context, sess *Session) (io.ReadCloser, error) {
	url := fmt.Sprintf("https://livetiming.formula1.com%sTimingData.jsonStream", sess.Path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			b = []byte("could not read response body")
		}
		resp.Body.Close()

		return nil, fmt.Errorf("got unexpected status code %d and resp %s", resp.StatusCode, string(b))
	}

	return resp.Body, nil
}
