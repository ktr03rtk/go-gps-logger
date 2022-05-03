package datawriter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ktr03rtk/go-gps-logger/pkg/dataconverter"
)

var (
	data_path = os.Getenv("DATA_PATH")
)

func init() {
	if err := os.MkdirAll(data_path, 0755); err != nil {
		log.Fatalf("Failed to create directory: %s", err)
	}
}

func Write(data_set dataconverter.ConvertedDataSet) {
	for k, v := range data_set {
		file_path := createFilePath(k, v.Timestamp)

		f, err := os.OpenFile(file_path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to create file: %s", err)
		}
		defer f.Close()

		bs, err := json.Marshal(v)
		if err != nil {
			log.Fatalf("Failed to json marshal: %s", err)
		}
		fmt.Fprintln(f, string(bs))
	}
}

func createFilePath(k string, s string) string {
	dir_path := filepath.Join(data_path, k)
	if err := os.MkdirAll(dir_path, 0755); err != nil {
		log.Fatalf("Failed to create diretory: %s", err)
	}
	time := timeFormat(s)
	file_name := fmt.Sprintf("%s-%s.raw", k, time)
	file_path := filepath.Join(dir_path, file_name)
	return file_path
}

func timeFormat(s string) string {
	time, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		log.Fatalf("Failed to convert time: %s", err)
	}
	formatted_time := time.Format(("2006-01-02-15-04-05"))
	return formatted_time
}
