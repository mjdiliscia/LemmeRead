package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/controller"
	"github.com/mjdiliscia/LemmeRead/model"
	"github.com/mjdiliscia/LemmeRead/ui"
	"go.elara.ws/go-lemmy"
)

const applicationName = "wip.drako.lemmeread"

type Application struct {
	LemmyClient *lemmy.Client
	LemmyContext context.Context
	GtkApplication *gtk.Application
	Window ui.MainWindow
	Model model.AppModel
	Controller controller.PostsController
}

func NewApplication() (app Application, err error) {
	app.GtkApplication, err = gtk.ApplicationNew(applicationName, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return Application{}, fmt.Errorf("Couldn't create Gtk Application: %s", err)
	}

	app.GtkApplication.Connect("activate", func() { app.onActivate() })

	return app, nil
}

func (app* Application) onActivate() {
	app.initMainWindow()
	app.setupControllers()
	app.lemmyStartup()
}

func (app *Application) initMainWindow() {
	log.Println("About to setup MainWindow...")
	err := app.Window.SetupMainWindow(&app.Model)
	if err != nil {
		log.Panic(err)
	}
	app.Window.Window.SetApplication(app.GtkApplication)
	log.Println("MainWindow setup finished.")
}

func (app *Application) setupControllers() {
	app.Controller.Init(&app.Window, &app.Model)
}

func (app *Application) lemmyStartup() {
	app.Model.Init(&app.Window)
	log.Println("About to initialize and login to Lemmy...")
	app.Model.InitializeLemmyClient("https://lemm.ee", "mjdiliscia", "qNZ^jyj2q.0@", func(err error) {
		if err != nil {
			log.Panic(err)
		}
		log.Println("Initialization finished.")
		log.Println("About to retrieve first page of posts...")
		app.Model.RetrieveMorePosts(func(err error) {
			if err != nil {
				log.Panic(err)
			}
			log.Println("Inital posts retrieval finished.")
		})
	})
}
