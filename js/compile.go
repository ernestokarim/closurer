package js

import (
	"fmt"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/ernestokarim/closurer/app"
	"github.com/ernestokarim/closurer/config"
	"github.com/ernestokarim/closurer/gss"
	"github.com/ernestokarim/closurer/hooks"
	"github.com/ernestokarim/closurer/soy"
)

func FullCompile() error {
	if err := hooks.PreCompile(); err != nil {
		return err
	}

	if err := gss.Compile(); err != nil {
		return err
	}

	if err := soy.Compile(); err != nil {
		return err
	}

	if err := Compile(); err != nil {
		return err
	}

	if err := hooks.PostCompile(); err != nil {
		return err
	}

	return nil
}

func Compile() error {
	conf := config.Current()
	target := conf.Js.CurTarget()

	if conf.Js == nil {
		return nil
	}

	if len(conf.Js.Inputs) == 0 {
		return nil
	}

	deps, _, err := GenerateDeps("compile")
	if err != nil {
		return err
	}

	args := []string{
		"-jar", path.Join(conf.Js.Compiler, "build", "compiler.jar"),
		"--js_output_file", path.Join(conf.Build, config.JS_NAME),
	}

	if conf.Library != nil {
		args = append(args,
			"--js", path.Join(conf.Library.Root, "closure", "goog", "base.js"),
			"--js", path.Join(conf.Library.Root, "closure", "goog", "deps.js"),
		)
	}

	args = append(args,
		"--js", filepath.Join(conf.Build, config.DEPS_NAME),
		"--js", filepath.Join(conf.Build, config.RENAMING_MAP_NAME),
	)

	if conf.Js.SideEffects == "" {
		args = append(args, "--output_wrapper", `(function(){%output%})();`)
	}

	for _, dep := range deps {
		if !strings.Contains(dep.Filename, "_test.js") {
			args = append(args, "--js", dep.Filename)
		}
	}

	if target.Defines != nil {
		for _, define := range target.Defines {
			// If it's not a boolean, quote it
			if define.Value != "true" && define.Value != "false" {
				define.Value = "\"" + define.Value + "\""
			}
			args = append(args, "--define", define.Name+"="+define.Value)
		}
	}

	if target.Mode == "ADVANCED" {
		args = append(args, "--compilation_level", "ADVANCED_OPTIMIZATIONS")
	} else if target.Mode == "SIMPLE" {
		args = append(args, "--compilation_level", "SIMPLE_OPTIMIZATIONS")
	} else if target.Mode == "WHITESPACE" {
		args = append(args, "--compilation_level", "WHITESPACE_ONLY")
	} else {
		return app.Errorf("RAW mode not allowed while compiling")
	}

	args = append(args, "--warning_level", target.Level)

	if conf.Js.Checks != nil {
		for _, check := range conf.Js.Checks.Errors {
			args = append(args, "--jscomp_error", check.Name)
		}
		for _, check := range conf.Js.Checks.Warnings {
			args = append(args, "--jscomp_warning", check.Name)
		}
		for _, check := range conf.Js.Checks.Offs {
			args = append(args, "--jscomp_off", check.Name)
		}
	}

	for _, extern := range conf.Js.Externs {
		args = append(args, "--externs", extern.File)
	}

	if conf.Js.Language != "" {
		args = append(args, "--language_in", conf.Js.Language)
	}

	if conf.Js.Formatting != "" {
		args = append(args, "--formatting", conf.Js.Formatting)
		args = append(args, "--debug", "true")
	}

	log.Println("Compiling JS:", target.Name)

	// Prepare the command
	cmd := exec.Command("java", args...)

	// Output it if asked to
	if config.OutputCmd {
		fmt.Println("java", strings.Join(cmd.Args, " "))
	}

	// Run the JS compiler
	output, err := cmd.CombinedOutput()
	if err != nil {
		if len(output) != 0 {
			fmt.Println(string(output))
		}

		return app.Errorf("exec error: %s", err)
	}

	if len(output) > 0 {
		log.Println("Output from JS compiler:\n", string(output))
	}

	log.Println("Done compiling JS!")

	return nil
}
