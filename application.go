package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/controller"
	"github.com/mjdiliscia/LemmeRead/model"
	"github.com/mjdiliscia/LemmeRead/view"
	"go.elara.ws/go-lemmy"
)

const applicationName = "io.github.mjdiliscia.lemmeread"

type Application struct {
	LemmyClient    *lemmy.Client
	LemmyContext   context.Context
	GtkApplication *gtk.Application
	View           view.MainView
	Model          model.AppModel
	Controller     controller.PostsController
}

func NewApplication() (app Application, err error) {
	app.GtkApplication, err = gtk.ApplicationNew(applicationName, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return Application{}, fmt.Errorf("Couldn't create Gtk Application: %s", err)
	}

	app.GtkApplication.Connect("activate", func() { app.onActivate() })

	return app, nil
}

func (app *Application) onActivate() {
	app.initAppModel()
	if app.Model.Configuration.HaveLemmyData() {
		app.initMainView()
		app.setupControllers()
		app.lemmyStartup()
	} else {
		app.initLoginView()
	}
}

func (app *Application) initAppModel() {
	app.Model.Init()
}

func (app *Application) initMainView() {
	log.Println("About to setup MainWindow...")
	err := app.View.SetupMainView(&app.Model)
	if err != nil {
		log.Panic(err)
	}
	app.View.Window.SetApplication(app.GtkApplication)
	log.Println("MainWindow setup finished.")
}

func (app *Application) setupControllers() {
	app.Controller.Init(&app.View, &app.Model)
}

func (app *Application) lemmyStartup() {
	log.Println("About to initialize and login to Lemmy...")
	err := app.Model.InitializeLemmyClient()
	app.onLemmyStarted(err)
}

func (app *Application) initLoginView() {
	var loginView view.LoginView
	err := loginView.SetupLoginView()
	if err != nil {
		log.Panic(err)
	}
	loginView.Window.SetApplication(app.GtkApplication)
	loginView.LoginClicked = func(server string, username string, password string) {
		loginView.DestroyWindow()
		app.initMainView()
		app.setupControllers()
		app.Model.InitializeLemmyClientWithLogin(server, username, password, app.onLemmyStarted)
	}
}

func (app *Application) onLemmyStarted(err error) {
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
}
