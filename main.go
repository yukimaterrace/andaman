package main

import (
	"yukimaterrace/andaman/controller"
)

func main() {
	// flow := factory.CreateGridSearchFlow()

	// start := time.Now().UnixNano()

	// flow.Start()
	// flow.WaitForCompletion()

	// end := time.Now().UnixNano()

	// log.Printf("elapsed: %f sec\n", float64(end-start)/(1.0e9))

	controller := controller.CreateController()
	controller.Logger.Fatal(controller.Start(":1323"))
}
