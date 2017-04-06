package elements

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func getAzureData(url string) (data []string) {
	fmt.Println("getAzureData('%s')", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return make([]string, 0)
	}
	req.Header.Set("Metadata", "true")
	resp, err := client.Do(req)
	if err != nil {
		return make([]string, 0)
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		data = append(data, strings.TrimRight(scanner.Text(), "\n"))
		if err != nil {
			break
		}
	}
	return
}

func crawlAzureData(url string) map[string]interface{} {
	data := make(map[string]interface{})
	urlData := getAzureData(url)

	var key string
	for _, line := range urlData {

		// replace hyphens with underscores for JSON keys
		key = strings.Replace(line, "-", "_", -1)

		d := getAzureData(url + line)
		if len(d) > 0 {
			data[key] = d[0]
		}
	}
	return data
}

func azureJsonData(url string) ([]byte, error) {
	data, err := json.MarshalIndent(crawlAzureData(url), "", "    ")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *Elements) GetAzureElements() (map[string]interface{}, error) {
	url := "http://169.254.169.254/metadata/latest/instance/"

	data, err := azureJsonData(url)
	if err != nil {
		return nil, fmt.Errorf("Error crawling azure data: %s", err)
	}

	ec2Elements := make(map[string]interface{})
	err = json.Unmarshal(data, &ec2Elements)
	if err != nil {
		return nil, fmt.Errorf("Error crawling azure data: %s", err)
	}

	return ec2Elements, nil
}
