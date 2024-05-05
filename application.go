package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
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
	var err error

	log.Println("About to setup MainWindow...")
	err = app.Window.SetupMainWindow(&app.Model)
	if err != nil {
		log.Panic(err)
	}
	app.Window.Window.SetApplication(app.GtkApplication)
	log.Println("MainWindow setup finished.")

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

	app.Window.OnPostListBottomReached = func() {
		app.Model.RetrieveMorePosts(func(err error) {
			if err != nil {
				log.Println(err)
			}
		})
	}
	app.Window.PostList.CommentButtonClicked = func (id int64) {
		app.Model.RetrieveComments(id, func(err error) {
			if err != nil {
				log.Println(err)
				return
			}
			app.Window.OpenComments(id)
		})
	}
}

func (app *Application) SetupLemmyClient(url string) (err error) {
	app.LemmyClient, err = lemmy.New(url)
	if err != nil {
		return fmt.Errorf("Couldn't create a Lemmy Client: %s", err)
	}

	return nil
}

func (app *Application) LoginLemmyClient(user, pass, totp string) (err error) {
	app.LemmyContext = context.Background()

	totpToken := lemmy.NewOptionalNil[string]()
	if len(totp) != 0 {
		totpToken = lemmy.NewOptional(totp)
	}

	err = app.LemmyClient.ClientLogin(app.LemmyContext, lemmy.Login{
		UsernameOrEmail: user,
		Password: pass,
		TOTP2FAToken: totpToken,
	})

	return
}

func (app *Application) PostsLemmyClient(page int64) (posts []lemmy.PostView, err error) {
	response, err := app.LemmyClient.Posts(app.LemmyContext, lemmy.GetPosts{
		Type: lemmy.NewOptional(lemmy.ListingTypeSubscribed),
		Page: lemmy.NewOptional(page+1),
	})
	if err != nil {
		return
	}

	posts = response.Posts
	return
}

func (app *Application) PostLemmyClient(postId int64) (post lemmy.PostView, err error) {
	response, err := app.LemmyClient.Post(app.LemmyContext, lemmy.GetPost{
		ID: lemmy.NewOptional(postId),
	})
	if err != nil {
		return
	}

	post = response.PostView
	return
}

func (app *Application) CommentsLemmyClient(postId int64) (comments []lemmy.CommentView, err error) {
	response, err := app.LemmyClient.Comments(app.LemmyContext, lemmy.GetComments{
		PostID: lemmy.NewOptional(postId),
		Limit: lemmy.NewOptional(int64(50)),
	})
	if err != nil {
		return
	}

	comments = response.Comments
	return
}
