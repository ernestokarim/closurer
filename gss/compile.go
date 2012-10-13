package gss

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/ernestokarim/closurer/app"
	"github.com/ernestokarim/closurer/cache"
	"github.com/ernestokarim/closurer/config"
	"github.com/ernestokarim/closurer/hooks"
	"github.com/ernestokarim/closurer/scan"
)

const (
	CSS_NAME          = "compiled.css"
	RENAMING_MAP_NAME = "renaming-map.js"
)

// Serves the compiled CSS file through the 
func CompiledCss(r *app.Request) error {
	r.W.Header().Set("Content-Type", "text/css")
	conf := config.Current()

	if err := hooks.PreCompile(); err != nil {
		return err
	}

	if err := Compile(); err != nil {
		return err
	}

	if err := hooks.PostCompile(); err != nil {
		return err
	}

	f, err := os.Open(filepath.Join(conf.Build, CSS_NAME))
	if os.IsNotExist(err) {
		fmt.Fprintln(r.W, "")
	} else if err != nil {
		return app.Error(err)
	}
	defer f.Close()

	if _, err := io.Copy(r.W, f); err != nil {
		return app.Error(err)
	}

	return nil
}

// Compiles the .gss files
func Compile() error {
	conf := config.Current()

	// Create/Clean the renaming map file to avoid compilation errors (the JS
	// compiler assumes there's a file with this name there).
	f, err := os.Create(path.Join(conf.Build, RENAMING_MAP_NAME))
	if err != nil {
		return err
	}
	f.Close()

	// Output early if there's no GSS files.
	if conf.RootGss == "" {
		return nil
	}

	gss, err := scan.Do(conf.RootGss, ".gss")
	if err != nil {
		return err
	}

	// No results, no compiling
	if len(gss) == 0 {
		return nil
	}

	// Check if the cached version is still ok
	modified := false
	for _, filepath := range gss {
		if m, err := cache.Modified("compile", filepath); err != nil {
			return err
		} else if m {
			modified = true
			break
		}
	}

	if !modified && !config.Build {
		return nil
	}

	log.Println("Compiling GSS...")

	// Run the soy compiler
	cmd := exec.Command(
		"java",
		"-jar", path.Join(conf.ClosureStylesheets, "build", "closure-stylesheets.jar"),
		"--output-file", filepath.Join(conf.Build, CSS_NAME),
		"--output-renaming-map-format", "CLOSURE_COMPILED",
		"--rename", "CLOSURE",
		"--output-renaming-map", path.Join(conf.Build, RENAMING_MAP_NAME))
	cmd.Args = append(cmd.Args, gss...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return app.Errorf("exec error: %s\n%s", err, string(output))
	}

	log.Println("Done compiling GSS!")

	return nil
}
