package main

import (
	"auto_image/pkg"
	"log"
	"time"
)

func main() {
	doJobs(
		pkg.ReadCsvFile(pkg.EnvVar("SOURCE_FILE")),
		pkg.EnvVar("SAVE_DIR"),
	)
}

func doJobs(jobs []string, saveDir string) {
	maxJobs := 5
	guard := make(chan struct{}, maxJobs)
	jobCount := make(chan int)
	start := time.Now()
	go func() {
		for jc := range jobCount {
			log.Printf("Count of completed jobs: %d.\n", jc)
			log.Printf("%d seconds have passed since start\n", int(time.Since(start).Seconds()))
		}
	}()

	log.Printf("Total count of jobs: %d\n", len(jobs))
	ticker := time.NewTicker(3 * time.Second)
	for i, job := range jobs {
		guard <- struct{}{}
		go func() {
			pic := &pkg.Pic{Name: job, SaveDir: saveDir}
			err := pic.Save()
			if err != nil {
				log.Println(err)
			}
			<-guard
		}()
		go func() {
			select {
			case <-ticker.C:
				jobCount <- i
			}
		}()
	}
}
