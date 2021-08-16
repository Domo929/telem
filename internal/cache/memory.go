package cache

import (
	"errors"
	"time"

	"github.com/Domo929/telem.git/internal/lap"
	"github.com/Domo929/telem.git/internal/livetiming"
	memCache "github.com/patrickmn/go-cache"
)

// Mem handles caching things in memory
type Mem struct {
	m *memCache.Cache
}

// NewMemoryCache creates a new instance of a Mem struct
func NewMemoryCache() *Mem {
	return &Mem{m: memCache.New(5*time.Minute, 10*time.Minute)}
}

// Check returns whether the associated session is in memory
func (m *Mem) Check(sess *livetiming.Session) bool {
	_, has := m.m.Get(sess.Name)
	return has
}

// Load returns the laps saved from memory, or an error if the lap has expired or a non-valid structure was saved to the cache
func (m *Mem) Load(sess *livetiming.Session) ([]lap.Lap, error) {
	lapsI, has := m.m.Get(sess.Name)
	if !has {
		return nil, errors.New("attempted to load laps that have expired or do not exist")
	}
	laps, ok := lapsI.([]lap.Lap)
	if !ok {
		return nil, errors.New("attempted to load laps that are not of type []laps.lap")
	}
	return laps, nil
}

// Save saves the race laps into the memory cache
func (m *Mem) Save(sess *livetiming.Session, laps []lap.Lap) error {
	m.m.Set(sess.Name, laps, memCache.DefaultExpiration)
	return nil
}
