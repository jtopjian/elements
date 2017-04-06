package elements

func (e *Elements) GetOpenStackElements() (map[string]interface{}, error) {
	return e.GetElementsFromJsonUrl("http://169.254.169.254/openstack/latest/meta_data.json")
}
