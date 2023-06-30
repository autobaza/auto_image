package main

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	photos := readCsvFile(envVar("SOURCE_FILE"))
	doJobs(photos, envVar("TARGET_DIR"))
}

func envVar(key string) string {
	viper.SetConfigFile("app.env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}
	return value
}

func doJobs(photos []string, dir string) {
	maxJobs := 5
	guard := make(chan struct{}, maxJobs)
	for _, name := range photos {
		guard <- struct{}{}
		go job(name, dir, guard)
	}
}

func job(name string, dir string, guard chan struct{}) {

	path := fmt.Sprintf("https://autobaza.kg/uploads/%s/%s/%s/%s", name[0:2], name[2:4], "1024x768", name)
	resp, err := http.Get(path)
	if err != nil {
		log.Println(err)
	}

	dirPath := fmt.Sprintf("%s/%s/%s", dir, name[0:2], name[2:4])
	err = os.MkdirAll(dirPath, 0700)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(dirPath + "/" + name)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Println(err)
	}

	<-guard
}

func readCsvFile(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	var result []string
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, rec[0])
	}

	return result
}
