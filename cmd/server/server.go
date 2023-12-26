package main

import (
	"log"
	"web-service/internal/pkg/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Println("app can not be created", err)
		return
	}
	app.Run()
}
