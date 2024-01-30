package common

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

const (
	KeepResource    = "keep"
	NonKeepResource = "non-keep"
)

// TimeFormats add all possible time formats
var TimeFormats = []string{
	time.RFC3339, //ISO8601 UTC time
	"2006-01-02 15:04:05",
	// add other time formats here
}

type ResourceInstance struct {
	Id          string
	Name        string
	CreatTime   string
	DefaultKeep bool
}

func ProcessScanCloudResources(resources, nonKeepResources []*ResourceInstance, resourceType, resourceName string) {
	ProcessResources(resources, resourceType, resourceName)

	ProcessNonKeepResources(nonKeepResources, resourceType, resourceName)
}

// ProcessResources Process all scanned cloud resources
func ProcessResources(resources []*ResourceInstance, resourceType, resourceName string) {
	data := make([][]string, len(resources))
	for i, r := range resources {
		isResourceKeep := CheckResourceNameKeep(r.Name)
		// some resources default to keep
		if r.DefaultKeep {
			isResourceKeep = KeepResource
		}

		creationDuration, err := DaysSinceCreation(r.CreatTime)
		if err != nil {
			log.Printf("[CRITAL] compute resource creation duration error: %v", err.Error())
		}

		data[i] = []string{
			resourceType,
			resourceName,
			r.Id,
			r.Name,
			isResourceKeep,
			creationDuration,
		}
	}
	err := WriteCsvFileData(SweeperResourceScanDir, ResourceScanHeader, data)
	if err != nil {
		log.Printf("[CRITAL] write csv file data error: %v", err.Error())
	}
}

// ProcessNonKeepResources Processing scanned non-keep cloud resources
func ProcessNonKeepResources(nonKeepResources []*ResourceInstance, resourceType, resourceName string) {
	data := make([][]string, len(nonKeepResources))
	for i, r := range nonKeepResources {
		data[i] = []string{
			resourceType,
			resourceName,
			r.Id,
			r.Name,
		}
	}
	err := WriteCsvFileData(SweeperNonKeepResourceScanDir, NonKeepResourceScanHeader, data)
	if err != nil {
		log.Printf("[CRITAL] write csv file data error: %v", err.Error())
	}
}

// CheckResourceNameKeep check whether to keep resource name
func CheckResourceNameKeep(name string) string {
	flag := CheckResourcePersist(name, "")
	if flag {
		return KeepResource
	}
	return NonKeepResource
}

// CheckResourcePersist check whether to persist resource
func CheckResourcePersist(name, createTime string) bool {
	if name == "" && createTime == "" {
		return false
	}
	parsedTime, _ := ParsedTime(createTime)

	createdWithin30Minutes := false
	if parsedTime != nil {
		createdWithin30Minutes = parsedTime.Add(time.Minute * 30).After(time.Now())
	}

	flag := regexp.MustCompile("^(keep|Default)").MatchString(name)
	return flag || createdWithin30Minutes
}

// DaysSinceCreation compute resource creation duration
func DaysSinceCreation(createTime string) (string, error) {
	parsedTime, err := ParsedTime(createTime)
	if err != nil {
		return "", err
	}

	duration := time.Since(*parsedTime)
	days := duration.Hours() / 24

	return fmt.Sprintf("%.2f", days), nil
}

// ParsedTime parse time
func ParsedTime(createTime string) (*time.Time, error) {
	if createTime == "" {
		return nil, nil
	}

	var parsedTime time.Time
	var err error

	timestamp, err := strconv.ParseInt(createTime, 10, 64)
	if err == nil {
		parsedTime = time.Unix(timestamp, 0)
	} else {
		// try parsing input strings using different time formats
		for _, format := range TimeFormats {
			parsedTime, err = time.Parse(format, createTime)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		log.Printf("[CRITAL] unable to parse create time[%s]", createTime)
		return nil, fmt.Errorf("unable to parse create time: %v", err.Error())
	}
	return &parsedTime, nil
}
