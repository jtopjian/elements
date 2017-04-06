package elements

func (e *Elements) GetOpenStackElements() (map[string]interface{}, error) {
	return e.GetElementsFromJsonUrl("http://169.254.169.254/openstack/2013-10-17/meta_data.json")
}
