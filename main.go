package main

import (
	"log"
	"yukimaterrace/andaman/app"
)

func main() {
	app.CreateApp().Start()

	log.Println("Andaman Started")
}
