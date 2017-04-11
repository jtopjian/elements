package output

import (
	"fmt"
	"sort"
	"strings"
)

type ShellOutput struct {
	Config Config
}

func (o *ShellOutput) GenerateOutput(elements interface{}) (string, error) {
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
			l := fmt.Sprintf("%s_%s", label, sanitizeKeyForShell(key))
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

func sanitizeKeyForShell(key string) string {
	// replace non-compliant shell characters with underscores,
	// per IEEE Std 1003.1-2001
	r := strings.NewReplacer(
		"`", "_",
		"~", "_",
		"!", "_",
		"@", "_",
		"#", "_",
		"$", "_",
		"%", "_",
		"^", "_",
		"&", "_",
		"*", "_",
		"(", "_",
		")", "_",
		"-", "_",
		"+", "_",
		"=", "_",
		"{", "_",
		"}", "_",
		"[", "_",
		"]", "_",
		"|", "_",
		":", "_",
		";", "_",
		"'", "_",
		"\"", "_",
		"/", "_",
		"?", "_",
		",", "_",
		".", "_",
		"<", "_",
		">", "_",
	)
	return r.Replace(key)
}
