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
