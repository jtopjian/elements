package elements

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func getData(url string) (data []string) {
	resp, err := http.Get(url)
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

func crawlData(url string) map[string]interface{} {
	data := make(map[string]interface{})
	urlData := getData(url)

	var key string
	for _, line := range urlData {

		// replace hyphens with underscores for JSON keys
		key = strings.Replace(line, "-", "_", -1)

		switch {
		default:
			d := getData(url + line)
			if len(d) > 0 {
				data[key] = d[0]
			}
		case line == "dynamic":
			break
		case line == "meta-data":
			data[key] = crawlData(url + line + "/")
		case line == "user-data":
			data[key] = strings.Join(getData(url+line+"/"), "")
		case strings.HasSuffix(line, "/"):
			data[key[:len(line)-1]] = crawlData(url + line)
		case strings.HasSuffix(url, "public-keys/"):
			keyId := strings.SplitN(line, "=", 2)[0]
			data[key] = crawlData(url + keyId + "/")
		}
	}
	return data
}

func jsonData(url string) ([]byte, error) {
	data, err := json.MarshalIndent(crawlData(url), "", "    ")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (e *Elements) GetAWSElements() (map[string]interface{}, error) {
	url := "http://169.254.169.254/latest/"

	data, err := jsonData(url)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal data from ec2 crawl: %s", err)
	}

	ec2Elements := make(map[string]interface{})
	err = json.Unmarshal(data, &ec2Elements)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal data from ec2 JSON: %s", err)
	}

	return ec2Elements, nil
}

/*
Attribution: much of the code in this file was lifted from the ec2_metadata_dump
project: https://github.com/thbishop/ec2_metadata_dump, which is:

Copyright (c) 2013 Thomas Bishop
*/
