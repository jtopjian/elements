package output

import (
	"fmt"
)

type InvalidOutput struct {
	Config Config
}

func (o *InvalidOutput) GenerateOutput(elements interface{}) (string, error) {
	return "", fmt.Errorf("Unrecognized output format: %s", o.Config.Format)
}
