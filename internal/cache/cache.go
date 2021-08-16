package cache

import (
	"github.com/Domo929/telem.git/internal/lap"
	"github.com/Domo929/telem.git/internal/livetiming"
)

type Cache interface {
	Check(sess *livetiming.Session) bool
	Load(sess *livetiming.Session) ([]lap.Lap, error)
	Save(sess *livetiming.Session, laps []lap.Lap) error
}

var (
	// Current is the cache that will be used by the Telemetry package
	Current Cache
)

// SetCache sets the cache that will be used by the Telemetry package
func SetCache(c Cache) {
	Current = c
}
