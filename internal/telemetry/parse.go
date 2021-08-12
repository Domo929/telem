package telemetry

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/Domo929/telem.git/internal/livetiming"
)

// Load returns the aggregated information for the race
func Load(sess *livetiming.Session) ([]Line, error) {
	lines, err := loadFromWeb(sess)
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func loadFromWeb(sess *livetiming.Session) ([]Line, error) {
	r, err := livetiming.GetData(context.Background(), sess)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return parse(r)
}

func parse(r io.Reader) ([]Line, error) {
	scanner := bufio.NewScanner(r)

	lines := make([]Line, 0, 10_000)

	for scanner.Scan() {
		text := scanner.Text()

		if strings.Contains(text, "Segments") {
			continue
		}

		firstCurlyNdx := strings.IndexRune(text, '{')

		// _ is for the session duration, eventually this will be used to map other data streams to laps
		_, jsonStr := text[:firstCurlyNdx], text[firstCurlyNdx:]

		line, err := lineUnmarshall([]byte(jsonStr))
		if err != nil {
			return nil, err
		}

		lines = append(lines, *line)
	}

	return lines, nil
}

func lineUnmarshall(b []byte) (*Line, error) {
	line := new(Line)
	if err := json.Unmarshal(b, line); err == nil {
		return line, nil
	}

	initialLine := new(InitialLine)
	if err := json.Unmarshal(b, initialLine); err == nil {
		return initialLine.ConvertToLine(), nil
	}

	return nil, errors.New("could not unmarshal json into Line or InitialLine")
}
