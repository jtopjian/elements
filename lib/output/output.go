package output

import (
	"encoding/json"
	"fmt"
	"sort"
)

type Config struct {
	Format string
}

type Output struct {
	Config Config
}

func (o *Output) Generate(elements interface{}) (string, error) {
	if elements == nil {
		return "", nil
	}

	switch o.Config.Format {
	case "json":
		return o.JSONOutput(elements)
	case "shell":
		return o.ShellOutput(elements)
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

func (o *Output) ShellOutput(elements interface{}) (string, error) {
	var keys []string
	var output string

	parsed := make(map[string]string)
	_, _, parsed = parseElementsForShell("elements", elements, parsed)

	for k := range parsed {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		output = fmt.Sprintf("%s%s=\"%s\"\n", output, k, parsed[k])
	}

	return output, nil
}

func parseElementsForShell(label string, elements interface{}, parsed map[string]string) (string, interface{}, map[string]string) {
	var e interface{}
	switch elements.(type) {
	case map[string]interface{}:
		for key, subElements := range elements.(map[string]interface{}) {
			l := fmt.Sprintf("%s_%s", label, key)
			_, e, parsed = parseElementsForShell(l, subElements, parsed)
		}
	case []interface{}:
		for i, subElements := range elements.([]interface{}) {
			l := fmt.Sprintf("%s_%d", label, i)
			_, e, parsed = parseElementsForShell(l, subElements, parsed)
		}
	default:
		parsed[label] = fmt.Sprintf("%v", elements)
		return label, nil, parsed
	}

	return label, e, parsed
}
