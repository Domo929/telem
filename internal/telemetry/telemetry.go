package telemetry

import (
	"context"

	"github.com/Domo929/telem.git/internal/livetiming"
)

// Load returns the aggregated information for the race
func Load(sess *livetiming.Session) ([]Lap, error) {
	lines, err := loadFromWeb(sess)
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func loadFromWeb(sess *livetiming.Session) ([]Lap, error) {
	r, err := livetiming.GetData(context.Background(), sess)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	rawLaps, err := parse(r)
	if err != nil {
		return nil, err
	}

	return aggregate(rawLaps)
}
