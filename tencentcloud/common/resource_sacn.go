package common

type ResourceInstance struct {
	Id   string
	Name string
}

func (r ResourceInstance) GetId() string {
	return r.Id
}

func (r ResourceInstance) GetName() string {
	return r.Name
}

func ProcessResources(resources []ResourceInstance, resourceType, resourceName string) {
	data := make([][]string, len(resources))
	for i, v := range resources {
		data[i] = []string{
			resourceType,
			resourceName,
			v.GetId(),
			v.GetName(),
			IsResourceKeep(v.GetName()),
		}
	}
	WriteCsvFileData(data)
}
