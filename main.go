package main

import (
	"log"
	ap "yukimaterrace/andaman/app"
	"yukimaterrace/andaman/config"
)

func main() {
	var app *ap.App
	if config.OandaPractice {
		app = ap.CreatePracticeApp()
	} else {
		app = ap.CreateApp()
	}

	app.Start()

	log.Println("Andaman Started")
}
