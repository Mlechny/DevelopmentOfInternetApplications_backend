package main

import (
	"awesomeProject/internal/pkg/app"
	"log"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Println("App can not be created", err)
		return
	}
	app.Run()
}
