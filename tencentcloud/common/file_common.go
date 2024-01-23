package common

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const (
	SweeperResourceScanDir = "../../../tmp/resource_scan/"

	KeepResource    = "keep"
	NonKeepResource = "non-keep"
)

// IsResourceKeep check whether to keep resource
func IsResourceKeep(name string) string {
	flag := regexp.MustCompile("^(keep|Default)").MatchString(name)
	if flag {
		return KeepResource
	}
	return NonKeepResource
}

// WriteCsvFileData write data to csv file
func WriteCsvFileData(data [][]string) error {
	log.Printf("[INFO] write csv file data[%v] start", len(data))

	count := 0
	defer func() {
		log.Printf("[INFO] write csv file data success count[%v]", count)
	}()

	if len(data) == 0 {
		return nil
	}

	err := os.MkdirAll(SweeperResourceScanDir, 0755)
	if err != nil {
		log.Printf("[CRITAL] create directory %s error: %v", SweeperResourceScanDir, err.Error())
		return err
	}

	currentDate := time.Now().Format("20060102")
	filePath := filepath.Join(SweeperResourceScanDir, currentDate+".csv")

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		err = GenerateCsvFile(filePath)
		if err != nil {
			log.Printf("[CRITAL] generate csv file error: %v", err.Error())
			return err
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[CRITAL] open csv file error: %v", err.Error())
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	for _, row := range data {
		err = writer.Write(row)
		if err != nil {
			log.Printf("[CRITAL] write data[%v] to csv file error: %v", row, err.Error())
			return err
		}
		count++
	}
	writer.Flush()

	return nil
}

// GenerateCsvFile generate when csv file does not exist
func GenerateCsvFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("[CRITAL] create csv file error: %v", err.Error())
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	header := []string{"ResourceType", "ResourceName", "InstanceId", "InstanceName", "Classification"}
	err = writer.Write(header)
	if err != nil {
		log.Printf("[CRITAL] write header to csv file error: %v", err.Error())
		return err
	}
	writer.Flush()

	return nil
}

func DeleteCsvFile() error {
	currentDate := time.Now().Format("20060102")
	filePath := SweeperResourceScanDir + currentDate + ".csv"

	_, err := os.Stat(filePath)

	if os.IsExist(err) {
		err = os.Remove(filePath)
		if err != nil {
			log.Printf("[CRITAL] delete csv file error: %v", err.Error())
			return err
		}
	}
	return nil
}
