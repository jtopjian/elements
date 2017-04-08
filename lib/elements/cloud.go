package elements

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	case "digitalocean":
		cloudElements, err = e.GetDigitalOceanElements()
		if err != nil {
			return nil, err
		}
	case "openstack":
		cloudElements, err = e.GetOpenStackElements()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Cloud provider '%s' is not supported.", provider)
	}

	cloudElements["provider"] = provider
	return cloudElements, nil
}

func (e *Elements) GetElementsFromJsonUrl(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving metadata from '%s': %s", url, err)
	}
	defer resp.Body.Close()

	jsonDecoder := json.NewDecoder(resp.Body)

	osElements := make(map[string]interface{})
	if err := jsonDecoder.Decode(&osElements); err != nil {
		return nil, fmt.Errorf("Error processing metadata JSON: %s", err)
	}

	return osElements, nil
}
