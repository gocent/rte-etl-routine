package main

import (
	"log"
	"rte-etl-routine/executor"
	"rte-etl-routine/scheduler"
	"time"
)

func main() {
	sc := scheduler.NewScheduler(10 * time.Second)
	ewExecutor := executor.EcoWatExecutor{}

	err := sc.Add(ewExecutor).Add(executor.EcoWatExecutor{}).Start(true)
	if err != nil {
		log.Fatal("Unable to start scheduler, program stopped. ", err)
	}
}
