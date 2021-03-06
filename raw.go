package main

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ernestokarim/closurer/app"
	"github.com/ernestokarim/closurer/config"
	"github.com/ernestokarim/closurer/gss"
	"github.com/ernestokarim/closurer/hooks"
	"github.com/ernestokarim/closurer/js"
	"github.com/ernestokarim/closurer/soy"
)

func RawOutput(r *app.Request) error {
	if err := hooks.PreCompile(); err != nil {
		return err
	}

	if err := gss.Compile(); err != nil {
		return err
	}

	if err := soy.Compile(); err != nil {
		return err
	}

	_, namespaces, err := js.GenerateDeps("input-production")
	if err != nil {
		return err
	}

	log.Println("Output RAW mode")

	conf := config.Current()
	content := bytes.NewBuffer(nil)

	base := path.Join(conf.Library.Root, "closure", "goog", "base.js")
	if err := addFile(content, base); err != nil {
		return err
	}

	if err := addFile(content, path.Join(conf.Build, config.RENAMING_MAP_NAME)); err != nil {
		return err
	}

	if err := addFile(content, path.Join(conf.Build, config.DEPS_NAME)); err != nil {
		return err
	}

	if err := hooks.PostCompile(); err != nil {
		return err
	}

	css := make([]byte, 0)
	if conf.Gss != nil {
		css, err = ioutil.ReadFile(filepath.Join(conf.Build, config.CSS_NAME))
		if err != nil {
			return app.Error(err)
		}
	}

	data := map[string]interface{}{
		"Content":    template.HTML(string(content.Bytes())),
		"Port":       config.Port,
		"LT":         template.HTML("<"),
		"Namespaces": template.HTML("'" + strings.Join(namespaces, "', '") + "'"),
		"Css":        template.HTML(template.JSEscapeString(string(css))),
	}
	r.W.Header().Set("Content-Type", "text/javascript")
	return r.ExecuteTemplate([]string{"raw"}, data)
}

func addFile(w io.Writer, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return app.Error(err)
	}
	defer f.Close()

	io.Copy(w, f)

	return nil
}
