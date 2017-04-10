package cloud

func GetOpenStackElements() (map[string]interface{}, error) {
	return GetElementsFromJsonUrl("http://169.254.169.254/openstack/latest/meta_data.json")
}
