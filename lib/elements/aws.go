package elements

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const aws_metadata_url = "http://169.254.169.254/latest/"

func getAWSData(url string) (data []string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return make([]string, 0)
	}
	return getCloudData(req)
}

func crawlAWSData(url string) map[string]interface{} {
	data := make(map[string]interface{})
	urlData := getAWSData(url)

	var key string
	for _, line := range urlData {

		// replace hyphens with underscores for JSON keys
		key = strings.Replace(line, "-", "_", -1)

		switch {
		default:
			d := getAWSData(url + line)
			if len(d) > 0 {
				data[key] = d[0]
			}
		case line == "dynamic":
			break
		case line == "meta-data":
			data[key] = crawlAWSData(url + line + "/")
		case line == "user-data":
			data[key] = strings.Join(getAWSData(url+line+"/"), "")
		case strings.HasSuffix(line, "/"):
			data[key[:len(line)-1]] = crawlAWSData(url + line)
		case strings.HasSuffix(url, "public-keys/"):
			keyId := strings.SplitN(line, "=", 2)[0]
			data[key] = crawlAWSData(url + keyId + "/")
		}
	}
	return data
}

func (e *Elements) GetAWSElements() (map[string]interface{}, error) {
	data, err := json.MarshalIndent(crawlAWSData(aws_metadata_url), "", "    ")
	if err != nil {
		return nil, fmt.Errorf("Error crawling aws metadata: %s", err)
	}

	elements := make(map[string]interface{})
	err = json.Unmarshal(data, &elements)
	if err != nil {
		return nil, fmt.Errorf("Error crawling aws metadata: %s", err)
	}

	return elements, nil
}

/*
Attribution: much of the code in this file was lifted from the ec2_metadata_dump
project: https://github.com/thbishop/ec2_metadata_dump, which is:

Copyright (c) 2013 Thomas Bishop
*/
