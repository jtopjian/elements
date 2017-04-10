package cloud

import "fmt"

func GetDigitalOceanElements() (map[string]interface{}, error) {
	elements, err := GetElementsFromJsonUrl("http://169.254.169.254/metadata/v1.json")
	if err != nil {
		return nil, fmt.Errorf("Error retrieving digitalocean metadata: %s", err)
	}
	delete(elements, "vendor_data")
	return elements, nil
}
