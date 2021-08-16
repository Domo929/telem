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

type Cache struct {
	basePath string
}

var (
	Local *Cache
)

func New(path string) (*Cache, error) {
	if _, err := os.ReadDir(path); err != nil {
		return nil, err
	}
	return &Cache{basePath: path}, nil
}

func SetLocal(c *Cache) {
	Local = c
}

func (c *Cache) Check(sess *livetiming.Session) bool {
	files, err := os.ReadDir(c.basePath)
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

func (c *Cache) Load(sess *livetiming.Session) ([]lap.Lap, error) {
	f, err := os.Open(filepath.Join(c.basePath, fmt.Sprintf("%s.gob", sess.String())))
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

func (c *Cache) Save(sess *livetiming.Session, lines []lap.Lap) error {
	f, err := os.Create(filepath.Join(c.basePath, fmt.Sprintf("%s.gob", sess.String())))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := gob.NewEncoder(f).Encode(lines); err != nil {
		return err
	}

	return nil
}
