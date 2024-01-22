package common

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	SWEEPER_RESOURCE_SCAN_DIR = "./tmp/resource_scan/"
)

func WriteCsvFileData(data [][]string) error {
	log.Printf("[INFO] write csv file data[%v] start", len(data))
	if len(data) == 0 {
		return nil
	}

	err := os.MkdirAll(SWEEPER_RESOURCE_SCAN_DIR, 0755)
	if err != nil {
		log.Printf("[CRITAL] create directory %s error: %v", SWEEPER_RESOURCE_SCAN_DIR, err.Error())
		return err
	}

	currentDate := time.Now().Format("20060102")
	filePath := filepath.Join(SWEEPER_RESOURCE_SCAN_DIR, currentDate+".csv")

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
		err := writer.Write(row)
		if err != nil {
			log.Printf("[CRITAL] write data to csv file error: %v", err.Error())
			return err
		}
	}
	writer.Flush()

	log.Printf("[INFO] write csv file data end")
	return nil
}

func GenerateCsvFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("[CRITAL] create csv file error: %v", err.Error())
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	header := []string{"CloudProduct", "Resource", "Id", "Name"}
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
	filePath := SWEEPER_RESOURCE_SCAN_DIR + currentDate + ".csv"

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

func PrintFile() {
	_, err := os.Stat(SWEEPER_RESOURCE_SCAN_DIR)
	if os.IsNotExist(err) {
		log.Printf("#############路径不存在: %s\n", SWEEPER_RESOURCE_SCAN_DIR)
		return
	}
	log.Printf("#############路径存在: %s\n", SWEEPER_RESOURCE_SCAN_DIR)

	err = filepath.Walk(SWEEPER_RESOURCE_SCAN_DIR, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("#############遇到错误: %v\n", err)
			return err
		}
		log.Printf("#############文件路径: %s, 文件名: %s\n", path, info.Name())
		return nil
	})
	if err != nil {
		log.Printf("#############错误: %v\n", err)
	}
}
