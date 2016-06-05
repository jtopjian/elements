package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/jtopjian/elements/utils"
)

type Elements struct {
	Elements     map[string]interface{}
	ElementsPath string
	ConfigDir    string
	mu           sync.Mutex
}

func New(configDir string, elementsPath string) (*Elements, error) {
	e := &Elements{
		Elements:     make(map[string]interface{}),
		ElementsPath: elementsPath,
		ConfigDir:    configDir,
	}

	systemPathRE := regexp.MustCompile("^System")
	externalPathRE := regexp.MustCompile("^External")

	switch {
	case systemPathRE.MatchString(elementsPath):
		elements := e.GetSystemElements()
		e.Add("System", elements)
	case externalPathRE.MatchString(elementsPath):
		externalElements, err := e.GetExternalElements()
		if err != nil {
			return e, err
		}

		e.Add("External", externalElements)
	default:
		elements := e.GetSystemElements()
		externalElements, err := e.GetExternalElements()
		if err != nil {
			return e, err
		}

		e.Add("System", elements)
		e.Add("External", externalElements)
	}

	return e, nil
}

func (e *Elements) Add(key string, value interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Elements[key] = value
}

func (e *Elements) SystemElements() error {
	var systemElements interface{}
	e.GetSystemElements()
	e.Elements["System"] = systemElements

	return nil
}

func (e *Elements) GetExternalElements() (map[string]interface{}, error) {
	externalElementsDir := fmt.Sprintf("%s/elements.d", e.ConfigDir)
	externalElements := make(map[string]interface{})

	d, err := os.Open(externalElementsDir)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	files, err := d.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("Unable to read %s: %s", externalElementsDir, err)
	}

	executableElements := make([]string, 0)
	staticElements := make([]string, 0)

	for _, fi := range files {
		name := filepath.Join(externalElementsDir, fi.Name())
		if utils.IsExecutable(fi) {
			executableElements = append(executableElements, name)
			continue
		}
		if strings.HasSuffix(name, ".json") {
			staticElements = append(staticElements, name)
		}
	}

	var wg sync.WaitGroup
	for _, p := range staticElements {
		p := p
		wg.Add(1)
		go e.elementsFromFile(p, &wg, externalElements)
	}
	for _, p := range executableElements {
		p := p
		wg.Add(1)
		go e.elementsFromExec(p, &wg, externalElements)
	}
	wg.Wait()

	return externalElements, nil
}

func (e *Elements) elementsFromFile(path string, wg *sync.WaitGroup, externalElements map[string]interface{}) error {
	defer wg.Done()

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Unable to read element from file %s: %s", path, err)
	}

	var result interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal data from file %s: %s", path, err)
	}

	externalElements[strings.TrimSuffix(filepath.Base(path), ".json")] = result

	return nil
}

func (e *Elements) elementsFromExec(path string, wg *sync.WaitGroup, externalElements map[string]interface{}) error {
	defer wg.Done()

	out, err := exec.Command(path).Output()
	if err != nil {
		return fmt.Errorf("Unable to execute command %s: %s", path, err)
	}

	var result interface{}
	err = json.Unmarshal(out, &result)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal data from command %s: %s", path, err)
	}

	externalElements[strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))] = result

	return nil
}

func (e *Elements) ElementsAtPath() (interface{}, error) {
	var data interface{}
	var value interface{}
	path_pieces := strings.Split(e.ElementsPath, ".")

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

func (e *Elements) Elements2JSON() (string, error) {
	elements, err := e.ElementsAtPath()
	if err != nil || elements == nil {
		return "", err
	}

	if _, ok := elements.([]interface{}); ok {
		if j, err := json.MarshalIndent(elements, " ", " "); err != nil {
			return "", err
		} else {
			return string(j), nil
		}
	} else {
		if _, ok := elements.(map[string]interface{}); ok {
			if j, err := json.MarshalIndent(elements, " ", " "); err != nil {
				return "", err
			} else {
				return string(j), nil
			}
		} else {
			return fmt.Sprintf("%s", elements), nil
		}
	}
}
