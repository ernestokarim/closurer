package cache

import (
	"os"
	"time"

	"github.com/ernestokarim/closurer/app"
	"github.com/ernestokarim/closurer/config"
)

var modificationCache = map[string]time.Time{}

// Checks if filename has been modified since the last time
// it was scanned. It so, or if it's not present in the cache,
// it returns true and stores the new time.
func Modified(dest, filename string) (bool, error) {
	if config.NoCache {
		return true, nil
	}

	name := dest + filename

	info, err := os.Lstat(filename)
	if err != nil {
		return false, app.Error(err)
	}

	modified, ok := modificationCache[name]

	if !ok || info.ModTime() != modified {
		modificationCache[name] = info.ModTime()
		return true, nil
	}

	return false, nil
}
