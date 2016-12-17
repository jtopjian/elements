package elements

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jtopjian/elements/utils"
)

func (e *Elements) GetExternalElements() (map[string]interface{}, error) {
	externalElementsDir := fmt.Sprintf("%s/elements.d", e.Config.Directory)
	externalElements := make(map[string]interface{})
	externalElementsDirExists := true

	_, err := os.Stat(externalElementsDir)
	if err != nil {
		if os.IsNotExist(err) {
			externalElementsDirExists = false
		} else {
			return nil, err
		}
	}

	if externalElementsDirExists {
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
	}

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
