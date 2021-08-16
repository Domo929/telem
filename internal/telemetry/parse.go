package telemetry

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/Domo929/telem.git/internal/lap"
)

func parse(r io.Reader) ([]lap.Lap, error) {
	scanner := bufio.NewScanner(r)

	lines := make([]lap.Lap, 0, 10_000)

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

func lineUnmarshall(b []byte) (*lap.Lap, error) {
	fullLap := new(lap.Lap)
	if err := json.Unmarshal(b, fullLap); err == nil {
		return fullLap, nil
	}

	initialLap := new(lap.InitialLap)
	if err := json.Unmarshal(b, initialLap); err == nil {
		return initialLap.ToLap(), nil
	}

	return nil, errors.New("could not unmarshal json into Lap or InitialLap")
}
