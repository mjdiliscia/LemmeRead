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
	maxPostImageSize = 400
	communityIconSize = 30
)

type MainWindow struct {
	Window *gtk.ApplicationWindow
}

func NewMainWindow(app *Application) (win MainWindow, err error) {
	if app == nil {
		err = fmt.Errorf("Received app is nil")
		return
	}

	builder, err := gtk.BuilderNewFromString(string(ui.MainWindowUI))
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

	vbox, err := getUIObject[gtk.Box](builder, "postContainer")
	if err != nil {
		err = fmt.Errorf("Couldn't find the vertical box: %s", err)
	}

	postsData, err := app.PostsLemmyClient()
	if err != nil {
		return
	}

	for _, post := range postsData {
		postUI, _ := getPostUI(post)
		vbox.PackStart(postUI, true, true, 0)
	}

	win.Window.Show()
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
	setWidgetProperty(builder, "description", func(textView *gtk.TextView) {
		if post.Post.Body.IsValid() {
			buffer, _ := textView.GetBuffer()
			buffer.SetText(post.Post.Body.ValueOrZero())
		} else {
			textView.Hide()
		}
	})
	setWidgetProperty(builder, "communityName", func(label *gtk.Label) {
		label.SetText(fmt.Sprintf("<span size=\"large\">%s</span>", post.Community.Title))
		label.SetUseMarkup(true)
	})
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
		spinner.SetRange(float64(post.Counts.Score)-1, float64(post.Counts.Score)+1)
		spinner.SetIncrements(1, 1)
		spinner.SetValue(float64(post.Counts.Score))
	})
	setWidgetProperty(builder, "comments", func(button *gtk.Button) { button.SetLabel(fmt.Sprintf("%d comments", post.Counts.Comments)) })

	if post.Post.ThumbnailURL.IsValid() {
		LoadImageFromURL(post.Post.ThumbnailURL.ValueOrZero(), func(pixbuf *gdk.Pixbuf, err error) {
			setImage(builder, pixbuf, "image", [2]int{maxPostImageSize, maxPostImageSize}, err)
		})
	}

	if post.Post.URL.IsValid() && (!post.Post.ThumbnailURL.IsValid() || post.Post.URL.ValueOrZero() != post.Post.ThumbnailURL.ValueOrZero()){
		setWidgetProperty(builder, "linkButton", func(link *gtk.LinkButton) {
			link.SetUri(post.Post.URL.ValueOrZero())
			link.Show()
		})
	}

	if post.Community.Icon.IsValid() {
		LoadImageFromURL(post.Community.Icon.ValueOrZero(), func(pixbuf *gdk.Pixbuf, err error) {
			setImage(builder, pixbuf, "communityIcon", [2]int{communityIconSize, communityIconSize}, err)
		})
	}

	postUI, err = getUIObject[gtk.Box](builder, "post")
	postUI.Unparent()

	return
}

func setImage(builder *gtk.Builder, pixbuf *gdk.Pixbuf, imageId string, maxSize [2]int, err error) {
	if err != nil {
		return
	}

	imageWidth := float64(pixbuf.GetWidth())
	imageWidthScale := imageWidth / float64(maxSize[0])
	imageHeight := float64(pixbuf.GetHeight())
	imageHeightScale := imageHeight / float64(maxSize[1])

	if  imageWidthScale > 1.0 || imageHeightScale > 1.0 {
		scale := math.Max(imageWidthScale, imageHeightScale)
		pixbuf, _ = pixbuf.ScaleSimple(int(imageWidth / scale), int(imageHeight / scale), gdk.INTERP_HYPER)
	}

	image, err := getUIObject[gtk.Image](builder, imageId)
	if err != nil {
		log.Println(err)
		return
	}

	image.SetFromPixbuf(pixbuf)
	image.Show()
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
