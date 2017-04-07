package elements

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	case "azure":
		cloudElements, err = e.GetAzureElements()
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

func getCloudData(req *http.Request) (data []string) {
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return make([]string, 0)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		data = append(data, strings.TrimRight(scanner.Text(), "\n"))
		if err != nil {
			break
		}
	}
	return
}
