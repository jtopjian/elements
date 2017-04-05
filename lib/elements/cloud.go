package elements

import (
	"fmt"
)

func (e *Elements) GetCloudElements(provider string) (map[string]interface{}, error) {

	var cloudElements map[string]interface{}
	var err error

	switch provider {
	case "aws":
		cloudElements, err = e.GetAWSElements()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Cloud provider '%s' is not supported.", provider)
	}

	cloudElements["provider"] = provider
	return cloudElements, nil
}
