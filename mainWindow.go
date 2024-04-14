package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/ui"
	"go.elara.ws/go-lemmy"
)

const (
	applicationTitle = "Lemme Read"
	maxPostImageSize = 400.0
	minWindowWidth   = 500.0
	minWindowHeight  = 600.0
)

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

	win.Window, err = getUIObject[gtk.ApplicationWindow](builder, "window")
	if err != nil {
		err = fmt.Errorf("Couldn't create application window: %s", err)
		return
	}

	win.Window.SetApplication(app.GtkApplication)
	win.Window.SetTitle(applicationTitle)
	win.Window.SetSizeRequest(minWindowWidth, minWindowHeight)

	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		err = fmt.Errorf("Couldn't make a scrolled window: %s", err)
		return
	}
	win.Window.Add(scrolledWindow)

	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		err = fmt.Errorf("Couldn't make a vertical box: %s", err)
	}
	scrolledWindow.Add(vbox)

	postsData, err := app.PostsLemmyClient()
	if err != nil {
		return
	}

	for _, post := range postsData {
		log.Println(post.Post.Name)
		postUI, _ := getPostUI(post)
		vbox.PackStart(postUI, true, true, 0)
	}

	scrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	win.Window.ShowAll()

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

func getPostUI(post lemmy.PostView) (postUI *gtk.Box, err error) {
	var (
		builder *gtk.Builder
	)

	builder, err = gtk.BuilderNewFromString(string(ui.PostUI))
	if err != nil {
		return
	}

	setWidgetProperty(builder, "title", func(label *gtk.Label) { label.SetText(post.Post.Name) })
	setWidgetProperty(builder, "description", func(label *gtk.Label) { label.SetText(post.Post.Body.ValueOrZero()) })
	setWidgetProperty(builder, "communityName", func(label *gtk.Label) { label.SetText(post.Community.Title) })
	setWidgetProperty(builder, "username", func(label *gtk.Label) {
		if post.Creator.DisplayName.IsValid() {
			label.SetText(post.Creator.DisplayName.ValueOrZero())
		} else {
			label.SetText(post.Creator.Name)
		}
	})
	setWidgetProperty(builder, "time", func(label *gtk.Label) {
		label.SetText(fmt.Sprintf("%s ago", time.Since(post.Post.Published).Round(time.Minute).String()))
	})
	setWidgetProperty(builder, "votes", func(spinner *gtk.SpinButton) {
		spinner.SetValue(float64(post.Counts.Score))
		spinner.SetRange(spinner.GetValue()-1, spinner.GetValue()+1)
		spinner.SetIncrements(1, 1)
	})
	setWidgetProperty(builder, "comments", func(button *gtk.Button) { button.SetLabel(fmt.Sprintf("%d comments", post.Counts.Comments)) })

	if post.Post.URL.IsValid() {
		LoadImageFromURL(post.Post.URL.ValueOrZero(), func(pixbuf *gdk.Pixbuf, err error) {
			SetPostImage(builder, pixbuf, err)
		})
	}

	postUI, err = getUIObject[gtk.Box](builder, "post")
	postUI.Unparent()

	return
}

func SetPostImage(builder *gtk.Builder, pixbuf *gdk.Pixbuf, err error) {
	if err != nil {
		return
	}

	imageHeight := float64(pixbuf.GetHeight())
	imageWidth := float64(pixbuf.GetWidth())
	if imageHeight > maxPostImageSize || imageWidth > maxPostImageSize {
		scale := maxPostImageSize / math.Max(imageHeight, imageWidth)
		pixbuf, _ = pixbuf.ScaleSimple(int(imageWidth*scale), int(imageHeight*scale), gdk.INTERP_HYPER)
	}

	image, err := getUIObject[gtk.Image](builder, "image")
	if err != nil {
		log.Println(err)
		return
	}
	image.SetFromPixbuf(pixbuf)
}

func setWidgetProperty[WType any](builder *gtk.Builder, widgetId string, setter func(widget *WType)) (err error) {
	widget, err := getUIObject[WType](builder, widgetId)
	if err != nil {
		log.Printf("Couldn't set property of '%s'", widgetId)
		return
	}
	setter(widget)
	return
}
