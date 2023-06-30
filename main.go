package main

import (
	"auto_image/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

var (
	jobsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "auto_image_jobs_processed_counter",
		Help: "The total count of processed jobs",
	})
	jobsDuration = promauto.NewSummary(prometheus.SummaryOpts{
		Name:       "auto_image_jobs_duration",
		Help:       "Jobs duration in milliseconds",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
)

func main() {

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

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
			err := pic.Save(jobsDuration)
			if err != nil {
				log.Println(err)
			}
			jobsCounter.Inc()
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
