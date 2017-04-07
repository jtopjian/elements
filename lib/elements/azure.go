package elements

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const azure_metadata_url = "http://169.254.169.254/metadata/latest/instance/"

func getAzureData(url string) (data []string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return make([]string, 0)
	}
	// Microsoft requires this header for metadata queries.
	req.Header.Set("Metadata", "true")
	return getCloudData(req)
}

func crawlAzureData(url string) map[string]interface{} {
	data := make(map[string]interface{})
	urlData := getAzureData(url)

	var key string
	for _, line := range urlData {

		// replace hyphens with underscores for JSON keys
		key = strings.Replace(line, "-", "_", -1)

		if strings.HasSuffix(line, "/") {
			data[key[:len(line)-1]] = crawlAzureData(url + line)
		} else {
			d := getAzureData(url + line)
			if len(d) > 0 {
				data[key] = d[0]
			}
		}
	}
	return data
}

func (e *Elements) GetAzureElements() (map[string]interface{}, error) {

	data, err := json.MarshalIndent(crawlAzureData(azure_metadata_url), "", "    ")
	if err != nil {
		return nil, fmt.Errorf("Error crawling azure metadata: %s", err)
	}

	elements := make(map[string]interface{})
	err = json.Unmarshal(data, &elements)
	if err != nil {
		return nil, fmt.Errorf("Error crawling azure metadata: %s", err)
	}

	return elements, nil
}
