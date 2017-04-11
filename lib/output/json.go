package output

import (
	"encoding/json"
	"fmt"
)

type JSONOutput struct {
	Config Config
}

func (o *JSONOutput) GenerateOutput(elements interface{}) (string, error) {
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
