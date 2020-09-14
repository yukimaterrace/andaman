package main

import (
	"yukimaterrace/andaman/factory"
)

func main() {
	flow := factory.CreateSimulationApp()

	flow.Start()
	flow.WaitForCompletion()
}
