package elements

func (e *Elements) GetDigitalOceanElements() (map[string]interface{}, error) {
	return e.GetElementsFromJsonUrl("http://169.254.169.254/metadata/v1.json")
}
