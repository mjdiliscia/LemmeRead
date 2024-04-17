package main

import (
	"fmt"
	"log"
	"math"

	//"os"
	//"strings"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/ui"
	"go.elara.ws/go-lemmy"
)

const (
	applicationTitle  = "Lemme Read"
	maxPostImageSize  = 580
	communityIconSize = 30
)

type MainWindow struct {
	app *Application
	Window *gtk.ApplicationWindow
	toolbar *gtk.Box
	postsContainer *gtk.Box
	currentPage int
}

func NewMainWindow(app *Application) (win MainWindow, err error) {
	if app == nil {
		err = fmt.Errorf("Received app is nil")
		return
	}
	win.app = app

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

	win.postsContainer, err = getUIObject[gtk.Box](builder, "postContainer")
	if err != nil {
		err = fmt.Errorf("Couldn't find the vertical box: %s", err)
		return
	}

	win.toolbar, err = getUIObject[gtk.Box](builder, "toolbar")
	if err != nil {
		err = fmt.Errorf("Couldn't find bottom toolbar: %s", err)
		return
	}

	moreButton, err := getUIObject[gtk.Button](builder, "more")
	if err != nil {
		err = fmt.Errorf("Couldn't find More button: %s", err)
		return
	}
	moreButton.Connect("clicked", func() { win.onMoreClicked() })

	postsData, err := app.PostsLemmyClient(int64(win.currentPage))
	if err != nil {
		return
	}
	win.fillPosts(postsData)

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

func (win *MainWindow)fillPosts(postsData []lemmy.PostView) {
	win.postsContainer.Remove(win.toolbar)

	for _, post := range postsData {
		postUI, _ := getPostUI(post)
		convertToCard(postUI)
		win.postsContainer.PackStart(postUI, false, false, 0)
	}
	win.postsContainer.PackStart(win.toolbar, false, false, 0)
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
		}
	})
	setWidgetProperty(builder, "descriptionScroll", func(scroll *gtk.ScrolledWindow) {
		if !post.Post.Body.IsValid() {
			scroll.Hide()
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

	if post.Post.URL.IsValid() && (!post.Post.ThumbnailURL.IsValid() || post.Post.URL.ValueOrZero() != post.Post.ThumbnailURL.ValueOrZero()) {
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

	if imageWidthScale > 1.0 || imageHeightScale > 1.0 {
		scale := math.Max(imageWidthScale, imageHeightScale)
		pixbuf, _ = pixbuf.ScaleSimple(int(imageWidth/scale), int(imageHeight/scale), gdk.INTERP_HYPER)
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

func convertToCard(box *gtk.Box) {
	box.SetName("card")

	cssProvider, _ := gtk.CssProviderNew()
	cssProvider.LoadFromData(string(ui.StyleCSS))
	context, _ := box.GetStyleContext()
	context.AddProvider(cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func (win *MainWindow)onMoreClicked() {
	win.currentPage++
	postsData, err := win.app.PostsLemmyClient(int64(win.currentPage))
	if err != nil {
		log.Printf("Error getting posts from page %d: %s", win.currentPage, err)
		return
	}

	win.fillPosts(postsData)
}
