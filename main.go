package main

import (
	"log"
	"os"
)

func main() {
	app, err := NewApplication()
	if err != nil {
		log.Panic(err)
	}

	app.GtkApplication.Run(os.Args)
}
