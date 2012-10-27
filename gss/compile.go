package gss

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/ernestokarim/closurer/app"
	"github.com/ernestokarim/closurer/cache"
	"github.com/ernestokarim/closurer/config"
	"github.com/ernestokarim/closurer/scan"
)

// Compiles the .gss files
func Compile() error {
	conf := config.Current()

	// Output early if there's no GSS files.
	if conf.RootGss == "" {
		if err := cleanRenamingMap(); err != nil {
			return err
		}

		return nil
	}

	gss, err := scan.Do(conf.RootGss, ".gss")
	if err != nil {
		return err
	}

	// No results, no compiling
	if len(gss) == 0 {
		if err := cleanRenamingMap(); err != nil {
			return err
		}

		return nil
	}

	// Check if the cached version is still ok
	modified := false
	for _, filepath := range gss {
		if m, err := cache.Modified("compile", filepath); err != nil {
			return err
		} else if m {
			modified = true
		}
	}

	if !modified && !config.Build {
		return nil
	}

	log.Println("Compiling GSS...")

	if err := cleanRenamingMap(); err != nil {
		return err
	}

	// Prepare the list of non-standard functions.
	funcs := []string{}
	if len(conf.NonStandardCssFuncs) > 0 {
		for _, f := range conf.NonStandardCssFuncs {
			funcs = append(funcs, "--allowed-non-standard-function")
			funcs = append(funcs, f)
		}
	}

	// Prepare the renaming map args
	renaming := []string{}
	if conf.RenameCss == "true" {
		renaming = []string{
			"--output-renaming-map-format", "CLOSURE_COMPILED",
			"--rename", "CLOSURE",
			"--output-renaming-map", path.Join(conf.Build, config.RENAMING_MAP_NAME),
		}
	}

	// Run the gss compiler
	cmd := exec.Command(
		"java",
		"-jar", path.Join(conf.ClosureStylesheets, "build", "closure-stylesheets.jar"),
		"--output-file", filepath.Join(conf.Build, config.CSS_NAME))
	cmd.Args = append(cmd.Args, funcs...)
	cmd.Args = append(cmd.Args, renaming...)
	cmd.Args = append(cmd.Args, gss...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) != 0 {
			fmt.Println(string(output))
			os.Exit(1)
		}

		return app.Errorf("exec error: %s", err)
	}

	log.Println("Done compiling GSS!")

	return nil
}

func cleanRenamingMap() error {
	conf := config.Current()

	// Create/Clean the renaming map file to avoid compilation errors (the JS
	// compiler assumes there's a file with this name there).
	f, err := os.Create(path.Join(conf.Build, config.RENAMING_MAP_NAME))
	if err != nil {
		return app.Error(err)
	}
	f.Close()

	return nil
}
