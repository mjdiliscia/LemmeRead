package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/ui"
	"go.elara.ws/go-lemmy"
)

const applicationName = "wip.drako.lemmeread"

type Application struct {
	LemmyClient *lemmy.Client
	LemmyContext context.Context
	GtkApplication *gtk.Application
	Window ui.MainWindow
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

	err = app.SetupLemmyClient("https://lemm.ee")
	if err != nil {
		log.Panicf("Couldn't connect to lemmy server: %s", err)
	}
	err = app.LoginLemmyClient("mjdiliscia", "qNZ^jyj2q.0@", "")
	if err != nil {
		log.Panicf("Couldn't login to lemmy: %s", err)
	}

	err = app.Window.SetupMainWindow()
	if err != nil {
		log.Panic(err)
	}
	app.Window.Window.SetApplication(app.GtkApplication)

	var page int64 = 0
	posts, err := app.PostsLemmyClient(page)
	if err != nil {
		log.Panic(err)
	}
	app.Window.PostList.FillPostsData(posts)
	app.Window.OnPostListBottomReached = func() {
		page++
		posts, err := app.PostsLemmyClient(page)
		if err != nil {
			log.Panic(err)
		}
		app.Window.PostList.FillPostsData(posts)
	}
	app.Window.PostList.CommentButtonClicked = func (id int64) {
		post, err := app.PostLemmyClient(id)
		if err != nil {
			log.Panic(err)
		}
		comments, err := app.CommentsLemmyClient(id)
		if err != nil {
			log.Panic(err)
		}
		app.Window.OpenComments(post, comments)
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
