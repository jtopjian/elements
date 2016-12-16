package lib

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(output interface{}) (string, error) {
	if _, ok := output.([]interface{}); ok {
		if j, err := json.MarshalIndent(output, " ", " "); err != nil {
			return "", err
		} else {
			return string(j), nil
		}
	} else {
		if _, ok := output.(map[string]interface{}); ok {
			if j, err := json.MarshalIndent(output, " ", " "); err != nil {
				return "", err
			} else {
				return string(j), nil
			}
		} else {
			return fmt.Sprintf("%s", output), nil
		}
	}
}
