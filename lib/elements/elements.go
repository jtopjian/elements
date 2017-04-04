package elements

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/erichs/cloudsysfs"
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
	cloudPathRE := regexp.MustCompile("^cloud")

	cloud := cloudsysfs.Detect()

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
	case cloudPathRE.MatchString(e.Config.Path):
		if cloud != "" {
			cloudElements, err := e.GetCloudElements(cloud)
			if err != nil {
				return nil, err
			}
			e.Add("cloud", cloudElements)
		}
	default:
		elements, err := e.GetSystemElements()
		if err != nil {
			return nil, err
		}
		e.Add("system", elements)

		externalElements, err := e.GetExternalElements()
		if err != nil {
			return nil, err
		}
		e.Add("external", externalElements)

		if cloud != "" {
			if cloudElements, err := e.GetCloudElements(cloud); err == nil {
				e.Add("cloud", cloudElements)
			}
		}
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
				} else {
					value = nil
				}
			} else {
				if _, ok := data.([]interface{}); ok {
					if len(data.([]interface{})) > i {
						value = data.([]interface{})[i]
					} else {
						value = nil
					}
				}
			}
			data = value
		}
	}

	return data, nil
}
