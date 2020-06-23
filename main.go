package main

import (
	"log"
	"yukimaterrace/andaman/app"
	"yukimaterrace/andaman/config"
)

func main() {
	var ap *app.App
	if config.OandaPractice {
		ap = app.CreatePracticeApp()
	} else {
		ap = app.CreateApp()
	}

	ap.Start()
	log.Println("Andaman Started")
}
