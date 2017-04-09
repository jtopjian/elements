package cloud

import (
	"fmt"
	"net/http"
)

const digitalocean_metadata_url = "http://169.254.169.254/metadata/v1.json"

func GetDigitalOceanElements() (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", digitalocean_metadata_url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving metadata: %s", err)
	}

	elements, err := GetElementsFromJsonUrl(req)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving digitalocean metadata: %s", err)
	}
	delete(elements, "vendor_data")
	return elements, nil
}
