package cloud

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func GetElements(provider string) (map[string]interface{}, error) {
	var cloudElements map[string]interface{}
	var err error

	switch provider {
	case "aws":
		cloudElements, err = GetAWSElements()
		if err != nil {
			return nil, err
		}
	case "azure":
		cloudElements, err = GetAzureElements()
		if err != nil {
			return nil, err
		}
	case "digitalocean":
		cloudElements, err = GetDigitalOceanElements()
		if err != nil {
			return nil, err
		}
	case "gce":
		cloudElements, err = GetGoogleComputeElements()
		if err != nil {
			return nil, err
		}
	case "openstack":
		cloudElements, err = GetOpenStackElements()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Cloud provider '%s' is not supported.", provider)
	}

	cloudElements["provider"] = provider
	return cloudElements, nil
}

func GetElementsFromJsonUrl(req *http.Request) (map[string]interface{}, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error retrieving JSON metadata from '%s': %s", req.URL, err)
	}

	jsonDecoder := json.NewDecoder(resp.Body)

	elements := make(map[string]interface{})
	if err := jsonDecoder.Decode(&elements); err != nil {
		return nil, fmt.Errorf("Error processing metadata JSON: %s", err)
	}

	return elements, nil
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
