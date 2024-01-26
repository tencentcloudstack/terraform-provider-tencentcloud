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

func ProcessResources(resources []*ResourceInstance, resourceType, resourceName string) {
	data := make([][]string, len(resources))
	for i, r := range resources {
		isResourceKeep := IsResourceKeep(r.Name)
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
	WriteCsvFileData(data)
}

// IsResourceKeep check whether to keep resource
func IsResourceKeep(name string) string {
	if name == "" {
		return NonKeepResource
	}

	flag := regexp.MustCompile("^(keep|Default)").MatchString(name)
	if flag {
		return KeepResource
	}
	return NonKeepResource
}

// DaysSinceCreation compute resource creation duration
func DaysSinceCreation(creatTime string) (string, error) {
	if creatTime == "" {
		return "", nil
	}

	var parsedTime time.Time
	var err error

	timestamp, err := strconv.ParseInt(creatTime, 10, 64)
	if err == nil {
		parsedTime = time.Unix(timestamp, 0)
	} else {
		// try parsing input strings using different time formats
		for _, format := range TimeFormats {
			parsedTime, err = time.Parse(format, creatTime)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		return "", fmt.Errorf("unable to parse creat time: %v", err.Error())
	}

	duration := time.Since(parsedTime)
	days := duration.Hours() / 24

	return fmt.Sprintf("%.2f", days), nil
}
