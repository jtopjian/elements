package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jtopjian/elements/lib"
)

var cmdServe cli.Command

type httpError struct {
	Error   error
	Message string
	Code    int
}

type httpHandler struct {
	Config lib.Config
	H      func(lib.Config, http.ResponseWriter, *http.Request) *httpError
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
	if err := hh.H(hh.Config, w, r); err != nil {
		http.Error(w, err.Message, err.Code)
		errAndExit(fmt.Errorf("Error serving requests: %s", err.Message))
	}
}

func actionServe(c *cli.Context) {
	config := lib.Config{
		Directory:    c.String("configdir"),
		Listen:       c.String("listen"),
		OutputFormat: c.String("format"),
		Path:         c.String("path"),
	}

	http.Handle("/elements", httpHandler{config, elementsHandler})
	http.Handle("/elements/", httpHandler{config, elementsHandler})
	debug.Printf("%s", http.ListenAndServe(config.Listen, nil))
}

func elementsHandler(config lib.Config, w http.ResponseWriter, r *http.Request) *httpError {
	pathre := regexp.MustCompile("^/elements/?")
	path := pathre.ReplaceAllString(r.URL.Path, "")

	path = strings.Replace(path, "/", ".", -1)
	debug.Printf("Element path requested: %s", path)

	elements := lib.Elements{
		Config: config,
	}

	output, err := elements.Get()
	if err != nil {
		return &httpError{err, "Error processing elements", 500}
	}

	formattedOutput, err := lib.PrintJSON(output)
	if err != nil {
		return &httpError{err, "Error processing elements", 500}
	}

	title := fmt.Sprintf("Elements %s", version)
	w.Header().Set("Server", title)
	if formattedOutput == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Element path not found"))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(formattedOutput))
	}

	return nil
}
