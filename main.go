package main

import (
	"log"
	"time"
	"yukimaterrace/andaman/factory"
)

func main() {
	flow := factory.CreateSimulationFlow()

	start := time.Now().UnixNano()

	flow.Start()
	flow.WaitForCompletion()

	end := time.Now().UnixNano()

	log.Printf("elapsed: %f sec\n", float64(end-start)/(1.0e9))
}
