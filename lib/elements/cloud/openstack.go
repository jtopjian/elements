package cloud

import (
	"fmt"
	"net/http"
)

const openstack_metadata_url = "http://169.254.169.254/openstack/latest/meta_data.json"

func GetOpenStackElements() (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", openstack_metadata_url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving metadata: %s", err)
	}

	return GetElementsFromJsonUrl(req)
}
