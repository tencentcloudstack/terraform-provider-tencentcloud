package common

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	SweeperResourceScanDir        = "../../../tmp/resource_scan/"
	SweeperNonKeepResourceScanDir = "../../../tmp/non_keep_resource_scan/"
)

var ResourceScanHeader = []string{"资源类型", "资源名称", "实例ID", "实例名称", "分类", "创建时长(天)", "创建者用户ID", "创建者用户名"}
var NonKeepResourceScanHeader = []string{"ResourceType", "ResourceName", "InstanceId", "InstanceName", "PrincipalId", "UserName"}

// WriteCsvFileData write data to csv file
func WriteCsvFileData(dirPath string, header []string, data [][]string) error {
	log.Printf("[INFO] write csv file data[%v] to path[%v] start", len(data), dirPath)

	count := 0
	defer func() {
		log.Printf("[INFO] write csv file data to path[%v] success count[%v]", dirPath, count)
	}()

	if len(data) == 0 {
		return nil
	}

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		log.Printf("[CRITAL] create directory %s error: %v", dirPath, err.Error())
		return err
	}

	currentDate := time.Now().Format("20060102")
	filePath := filepath.Join(dirPath, currentDate+".csv")

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		err = GenerateCsvFile(filePath, header)
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
func GenerateCsvFile(filePath string, header []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("[CRITAL] create csv file error: %v", err.Error())
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.Write(header)
	if err != nil {
		log.Printf("[CRITAL] write header to csv file error: %v", err.Error())
		return err
	}
	writer.Flush()

	return nil
}
