package cloud

import (
	"fmt"
	"net/http"
)

const googlecompute_metadata_url = "http://metadata.google.internal/computeMetadata/v1/?recursive=true"

func GetGoogleComputeElements() (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", googlecompute_metadata_url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving metadata: %s", err)
	}
	// Google requires this header for metadata queries.
	req.Header.Set("Metadata-Flavor", "Google")
	return GetElementsFromJsonUrl(req)
}
