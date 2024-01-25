package common

type ResourceInstance struct {
	Id          string
	Name        string
	DefaultKeep bool
}

func ProcessResources(resources []*ResourceInstance, resourceType, resourceName string) {
	data := make([][]string, len(resources))
	for i, r := range resources {
		isResourceKeep := IsResourceKeep(r.Name)
		// some resources default to keep
		if r.DefaultKeep {
			isResourceKeep = KeepResource
		}

		data[i] = []string{
			resourceType,
			resourceName,
			r.Id,
			r.Name,
			isResourceKeep,
		}
	}
	WriteCsvFileData(data)
}
