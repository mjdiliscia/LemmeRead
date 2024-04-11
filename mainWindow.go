package main

import (
	"fmt"
	"math"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/ui"
	"go.elara.ws/go-lemmy"
)

const applicationTitle = "Lemme Read"

type MainWindow struct {
	Window *gtk.ApplicationWindow
}

func NewMainWindow(app *Application) (win MainWindow, err error) {
	if app == nil {
		err = fmt.Errorf("Received app is nil")
		return
	}

	builder, err := gtk.BuilderNewFromString(ui.MainWindowUI)
	if err != nil {
		err = fmt.Errorf("Couldn't make the main window builder: %s", err)
		return
	}

	//win.Window, err = gtk.ApplicationWindowNew(app.GtkApplication)
	win.Window, err = getUIObject[gtk.ApplicationWindow](builder, "window")
    if err != nil {
		err = fmt.Errorf("Couldn't create application window: %s", err)
		return
    }

	win.Window.SetApplication(app.GtkApplication)
    win.Window.SetTitle(applicationTitle)

	builder, err = gtk.BuilderNewFromString(ui.PostUI)
	if err != nil {
		return
	}

	postsData, err := app.PostsLemmyClient()
	if err != nil {
		return
	}
	postUI, err := getPostUI(postsData[0])
	win.Window.Add(postUI)

	return win, nil
}

func getUIObject[OType any](builder *gtk.Builder, objectId string) (object *OType, err error) {
	obj, err := builder.GetObject(objectId)
	if err != nil {
		return
	}
	object, ok := any(obj).(*OType)
	if ok {
		return object, nil
	} else {
		return object, fmt.Errorf("Object '%s' can't be correctly casted.", objectId)
	}
}

func getPostUI(post lemmy.PostView) (postUI gtk.IWidget, err error) {
	var (
		builder *gtk.Builder
	)

	builder, err = gtk.BuilderNewFromString(ui.PostUI)
	if err != nil {
		return
	}

	postUI, err = getUIObject[gtk.Box](builder, "post")
	if err != nil {
		return
	}

	title, err := getUIObject[gtk.Label](builder, "title")
	if err != nil {
		return
	}
	title.SetText(post.Post.Name)

	description, err := getUIObject[gtk.Label](builder, "description")
	if err != nil {
		return
	}
	description.SetText(post.Post.Body.ValueOrZero())

	community, err := getUIObject[gtk.Label](builder, "communityName")
	if err != nil {
		return
	}
	community.SetText(post.Community.Name)

	username, err := getUIObject[gtk.Label](builder, "username")
	if err != nil {
		return
	}
	username.SetText(post.Creator.Name)

	timeUI, err := getUIObject[gtk.Label](builder, "time")
	if err != nil {
		return
	}
	timeUI.SetText(fmt.Sprintf("%s ago", time.Since(post.Post.Published).Round(time.Minute).String()))

	votes, err := getUIObject[gtk.SpinButton](builder, "votes")
	if err != nil {
		return
	}
	votes.SetValue(float64(post.Counts.Score))
	votes.SetRange(math.Max(votes.GetValue()-1, 0), votes.GetValue()+1)
	votes.SetIncrements(1, 1)

	comments, err := getUIObject[gtk.Button](builder, "comments")
	if err != nil {
		return
	}
	comments.SetLabel(fmt.Sprintf("%d comments", post.Counts.Comments))

	return
}
