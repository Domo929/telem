package cache

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Domo929/telem.git/internal/lap"
	"github.com/Domo929/telem.git/internal/livetiming"
)

// File handles caching files locally on the computers disk
type File struct {
	basePath string
}

// NewFileCache validates the given path and then returns the File struct
func NewFileCache(path string) (*File, error) {
	if _, err := os.ReadDir(path); err != nil {
		return nil, err
	}
	return &File{basePath: path}, nil
}

// Check checks the File's basepath for a file that matches the sessions signature
func (cf *File) Check(sess *livetiming.Session) bool {
	files, err := os.ReadDir(cf.basePath)
	if err != nil {
		return false
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), sess.String()) {
			return true
		}
	}

	return false
}

// Load loads a session's telemetry from local cache
func (cf *File) Load(sess *livetiming.Session) ([]lap.Lap, error) {
	f, err := os.Open(filepath.Join(cf.basePath, fmt.Sprintf("%s.gob", sess.String())))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := make([]lap.Lap, 0)
	if err := gob.NewDecoder(f).Decode(&lines); err != nil {
		return nil, err
	}

	return lines, nil
}

// Save saves a session's telemetry to local cache
func (cf *File) Save(sess *livetiming.Session, laps []lap.Lap) error {
	f, err := os.Create(filepath.Join(cf.basePath, fmt.Sprintf("%s.gob", sess.String())))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(laps); err != nil {
		return err
	}

	return nil
}
