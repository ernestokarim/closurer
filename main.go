package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/ernestokarim/closurer/app"
	"github.com/ernestokarim/closurer/config"
	"github.com/ernestokarim/closurer/js"
	"github.com/ernestokarim/closurer/test"

	"github.com/gorilla/mux"
)

func main() {
	flag.Parse()

	if err := config.Load(config.ConfPath); err != nil {
		log.Fatal(err)
	}

	if config.Build {
		if err := build(); err != nil {
			err.(*app.AppError).Log()
		}
	} else {
		serve()
	}
}

func serve() {
	r := mux.NewRouter().StrictSlash(true)
	http.Handle("/", r)

	r.Handle("/", app.Handler(home))
	r.Handle("/compile", app.Handler(compile))
	r.Handle("/input/{name:.+}", app.Handler(Input))
	r.Handle("/test/all", app.Handler(test.TestAll))
	r.Handle("/test/list", app.Handler(test.TestList))
	r.Handle("/test/{name:.+}", app.Handler(test.Main))

	log.Printf("Started closurer server on http://localhost%s/\n", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, nil))
}

func home(r *app.Request) error {
	return r.ExecuteTemplate([]string{"home"}, nil)
}

func compile(r *app.Request) error {
	conf := config.Current()
	target := conf.Js.CurTarget()

	if target.Mode == "RAW" {
		return RawOutput(r)
	}

	return js.CompiledJs(r)
}
