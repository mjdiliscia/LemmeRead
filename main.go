package main

import (
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"go.elara.ws/go-lemmy"
)

const applicationName = "wip.drako.lemmeread"

func main() {
	app, err := gtk.ApplicationNew(applicationName, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Panic("Couldn't create Gtk Application.")
	}

	app.Connect("activate", func() { onActivate(app) })

	app.Run(os.Args)
}

func onActivate(app *gtk.Application) {
	_, err := lemmy.New("https://lemm.ee")
	if err != nil {
		log.Panic("Couldn't create a Lemmy Client.")
	}

    appWindow, err := gtk.ApplicationWindowNew(app)
    if err != nil {
        log.Fatal("Couldn't create application window.", err)
    }

    // Set ApplicationWindow Properties
    appWindow.SetTitle("Lemme Read")
    appWindow.SetDefaultSize(400, 400)
    appWindow.Show()
}
