package elements

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Config struct {
	Directory string
	Listen    string
	Path      string
}

type Elements struct {
	Config   Config
	Elements map[string]interface{}
	mu       sync.Mutex
}

func (e *Elements) Add(key string, value interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.Elements == nil {
		e.Elements = make(map[string]interface{})
	}

	e.Elements[key] = value
}

func (e *Elements) Get() (interface{}, error) {
	systemPathRE := regexp.MustCompile("^system")
	externalPathRE := regexp.MustCompile("^external")

	switch {
	case systemPathRE.MatchString(e.Config.Path):
		elements, err := e.GetSystemElements()
		if err != nil {
			return nil, err
		}
		e.Add("system", elements)
	case externalPathRE.MatchString(e.Config.Path):
		externalElements, err := e.GetExternalElements()
		if err != nil {
			return nil, err
		}
		e.Add("external", externalElements)
	default:
		elements, err := e.GetSystemElements()
		if err != nil {
			return nil, err
		}

		externalElements, err := e.GetExternalElements()
		if err != nil {
			return nil, err
		}
		e.Add("system", elements)
		e.Add("external", externalElements)
	}

	return e.ElementsAtPath()
}

func (e *Elements) ElementsAtPath() (interface{}, error) {
	var data interface{}
	var value interface{}
	path_pieces := strings.Split(e.Config.Path, ".")

	// Convert the Element structure into a generic interface{}
	// by first converting it to JSON and then decoding it.
	j, err := json.Marshal(e.Elements)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(strings.NewReader(string(j)))
	d.UseNumber()
	if err := d.Decode(&data); err != nil {
		return nil, err
	}

	// Walk through the given path.
	// If there's a result, print it.
	if len(path_pieces) > 1 {
		for _, p := range path_pieces {
			i, err := strconv.Atoi(p)
			if err != nil {
				if _, ok := data.(map[string]interface{}); ok {
					value = data.(map[string]interface{})[p]
				}
			} else {
				if _, ok := data.([]interface{}); ok {
					if len(data.([]interface{})) >= i {
						value = data.([]interface{})[i]
					}
				}
			}
			data = value
		}
	}

	return data, nil
}
