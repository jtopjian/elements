package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/urfave/cli"

	e "github.com/jtopjian/elements/lib/elements"
	o "github.com/jtopjian/elements/lib/output"
)

var cmdServe cli.Command

type httpConfig struct {
	EConfig e.Config
	OConfig o.Config
}

type httpError struct {
	Error   error
	Message string
	Code    int
}

type httpHandler struct {
	C httpConfig
	H func(httpConfig, http.ResponseWriter, *http.Request) *httpError
}

func init() {
	cmdServe = cli.Command{
		Name:   "serve",
		Usage:  "Serve elements over HTTP",
		Action: actionServe,
		Flags: []cli.Flag{
			&flagConfigDir,
			&flagDebug,
			&flagFormat,
			&flagListen,
		},
	}
}

func (hh httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := hh.H(hh.C, w, r); err != nil {
		http.Error(w, err.Message, err.Code)
		errAndExit(fmt.Errorf("Error serving requests: %s", err.Message))
	}
}

func actionServe(c *cli.Context) {
	eConfig := e.Config{
		Directory: c.String("configdir"),
		Listen:    c.String("listen"),
		Path:      c.String("path"),
	}

	oConfig := o.Config{
		Format: c.String("format"),
	}

	config := httpConfig{
		eConfig,
		oConfig,
	}

	http.Handle("/elements", httpHandler{config, elementsHandler})
	http.Handle("/elements/", httpHandler{config, elementsHandler})
	debug.Printf("%s", http.ListenAndServe(eConfig.Listen, nil))
}

func elementsHandler(config httpConfig, w http.ResponseWriter, r *http.Request) *httpError {
	pathre := regexp.MustCompile("^/elements/?")
	path := pathre.ReplaceAllString(r.URL.Path, "")

	path = strings.Replace(path, "/", ".", -1)
	debug.Printf("Element path requested: %s", path)

	if v := r.FormValue("format"); v != "" {
		debug.Printf("Format override requested: %s", v)
		config.OConfig.Format = v
	}

	output := o.Output{
		Config: config.OConfig,
	}

	config.EConfig.Path = path
	elements := e.Elements{
		Config: config.EConfig,
	}

	collectedElements, err := elements.Get()
	if err != nil {
		return &httpError{err, "Error collecting elements", 500}
	}

	title := fmt.Sprintf("Elements %s", version)
	w.Header().Set("Server", title)

	formattedOutput, outputErr := output.Generate(collectedElements)

	if outputErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(outputErr.Error()))
	} else if formattedOutput == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Element path not found"))
	} else {
		switch config.OConfig.Format {
		case "json":
			w.Header().Set("Content-Type", "application/json")
		case "shell":
			w.Header().Set("Content-Type", "text/plain")
		}

		w.Write([]byte(formattedOutput))
	}

	return nil
}
