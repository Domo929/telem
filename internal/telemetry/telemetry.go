package telemetry

import (
	"context"
	"log"

	"github.com/Domo929/telem.git/internal/cache"
	"github.com/Domo929/telem.git/internal/lap"

	"github.com/Domo929/telem.git/internal/livetiming"
)

// Load returns the aggregated information for the race
func Load(sess *livetiming.Session) ([]lap.Lap, error) {
	if cache.Local.Check(sess) {
		log.Printf("Loaded %s from local cache", sess.String())
		return cache.Local.Load(sess)
	}
	laps, err := loadFromWeb(sess)
	if err != nil {
		return nil, err
	}

	log.Printf("Loaded %s from web", sess.String())

	if err := cache.Local.Save(sess, laps); err != nil {
		return nil, err
	}

	return laps, nil
}

func loadFromWeb(sess *livetiming.Session) ([]lap.Lap, error) {
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
