package main

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/ernestokarim/closurer/app"
	"github.com/ernestokarim/closurer/cache"
	"github.com/ernestokarim/closurer/config"
)

var loadCacheOnce sync.Once

// Called before each compilation task. It load the caches
// and reload the confs if needed.
func PreCompileActions() error {
	if err := config.ReadFromFile(config.ConfPath); err != nil {
		return err
	}

	conf := config.Current()

	if err := os.MkdirAll(conf.Build, 0755); err != nil {
		return app.Error(err)
	}

	var err error
	loadCacheOnce.Do(func() {
		err = cache.Load(filepath.Join(conf.Build, "cache"))
	})

	return err
}

// Called after each compilation tasks. It saves the caches.
func PostCompileActions() error {
	conf := config.Current()
	return cache.Dump(filepath.Join(conf.Build, "cache"))
}
