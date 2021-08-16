package telemetry

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

func parse(r io.Reader) ([]Lap, error) {
	scanner := bufio.NewScanner(r)

	lines := make([]Lap, 0, 10_000)

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

func lineUnmarshall(b []byte) (*Lap, error) {
	line := new(Lap)
	if err := json.Unmarshal(b, line); err == nil {
		return line, nil
	}

	initialLine := new(InitialLap)
	if err := json.Unmarshal(b, initialLine); err == nil {
		return initialLine.toLine(), nil
	}

	return nil, errors.New("could not unmarshal json into Lap or InitialLap")
}
