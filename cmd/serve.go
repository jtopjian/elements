package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
	elements "github.com/jtopjian/elements/lib"
)

var cmdServe cli.Command

type httpError struct {
	Error   error
	Message string
	Code    int
}

type httpHandler struct {
	configDir string
	H         func(string, http.ResponseWriter, *http.Request) *httpError
}

func init() {
	cmdServe = cli.Command{
		Name:   "serve",
		Usage:  "Serve elements over HTTP",
		Action: actionServe,
		Flags: []cli.Flag{
			&flagConfigDir,
			&flagDebug,
			&flagListen,
		},
	}
}

func (hh httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := hh.H(hh.configDir, w, r); err != nil {
		http.Error(w, err.Message, err.Code)
		errAndExit(fmt.Errorf("Error serving requests: %s", err.Message))
	}
}

func actionServe(c *cli.Context) {
	configDir := c.String("configdir")
	http.Handle("/elements", httpHandler{configDir, elementsHandler})
	http.Handle("/elements/", httpHandler{configDir, elementsHandler})
	debug.Printf("%s", http.ListenAndServe(c.String("listen"), nil))
}

func elementsHandler(configDir string, w http.ResponseWriter, r *http.Request) *httpError {
	pathre := regexp.MustCompile("^/elements/?")
	path := pathre.ReplaceAllString(r.URL.Path, "")

	path = strings.Replace(path, "/", ".", -1)
	debug.Printf("Element path requested: %s", path)

	e, err := elements.New(configDir, path)
	if err != nil {
		return &httpError{err, "Error processing elements", 500}
	}

	elements, err := e.Elements2JSON()
	if err != nil {
		return &httpError{err, "Error processing elements", 500}
	}

	title := fmt.Sprintf("Elements %s", version)
	w.Header().Set("Server", title)
	if elements == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Element path not found"))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(elements))
	}

	return nil
}
