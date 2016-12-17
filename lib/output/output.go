package output

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Format string
}

type Output struct {
	Config Config
}

func (o *Output) Generate(elements interface{}) (string, error) {
	switch o.Config.Format {
	case "json":
		return o.JSONOutput(elements)
	}

	return "", fmt.Errorf("Unrecognized output format: %s", o.Config.Format)
}

func (o *Output) JSONOutput(elements interface{}) (string, error) {
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
