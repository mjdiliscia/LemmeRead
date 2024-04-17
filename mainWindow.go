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
	stack *gtk.Stack
	postView *gtk.ScrolledWindow
	mainView *gtk.ScrolledWindow
	commentsView CommentsView
	currentPage int
}

type CommentsView struct {
	title *gtk.Label
	communityIcon *gtk.Image
	communityName *gtk.Label
	username *gtk.Label
	timestamp *gtk.Label
	link *gtk.LinkButton
	image *gtk.Image
	description *gtk.TextView
	votes *gtk.SpinButton
	comments *gtk.TextView
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

	win.stack, err = getUIObject[gtk.Stack](builder, "stack")
	if err != nil {
		err = fmt.Errorf("Couldn't find stack: %s", err)
		return
	}

	win.postView, err = getUIObject[gtk.ScrolledWindow](builder, "postView")
	if err != nil {
		err = fmt.Errorf("Couldn't find Post View: %s", err)
		return
	}

	win.mainView, err = getUIObject[gtk.ScrolledWindow](builder, "mainView")
	if err != nil {
		err = fmt.Errorf("Couldn't find Main View: %s", err)
		return
	}

	win.commentsView, err = NewCommentsView(builder)
	if err != nil {
		return
	}

	moreButton, err := getUIObject[gtk.Button](builder, "more")
	if err != nil {
		err = fmt.Errorf("Couldn't find More button: %s", err)
		return
	}
	moreButton.Connect("clicked", func() { win.onMoreClicked() })

	closeCommentsButton, err := getUIObject[gtk.Button](builder, "closeComments")
	if err != nil {
		err = fmt.Errorf("Couldn't find Close Comments button: %s", err)
		return
	}
	closeCommentsButton.Connect("clicked", func() { win.onCloseComments() })

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
		postUI, _ := win.getPostUI(post)
		convertToCard(postUI)
		win.postsContainer.PackStart(postUI, false, false, 0)
	}
	win.postsContainer.PackStart(win.toolbar, false, false, 0)
}

func (win *MainWindow)getPostUI(post lemmy.PostView) (postUI *gtk.Box, err error) {
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
	setWidgetProperty(builder, "comments", func(button *gtk.Button) {
		button.SetLabel(fmt.Sprintf("%d comments", post.Counts.Comments))
		button.Connect("clicked", func() { win.onOpenComments(post.Post.ID) })
	})

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

func setDirectImage(image *gtk.Image, pixbuf *gdk.Pixbuf, maxSize [2]int, err error) {
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

	image.SetFromPixbuf(pixbuf)
	image.Show()
}

func setImage(builder *gtk.Builder, pixbuf *gdk.Pixbuf, imageId string, maxSize [2]int, err error) {
	if err != nil {
		return
	}

	image, err := getUIObject[gtk.Image](builder, imageId)
	if err != nil {
		log.Println(err)
		return
	}

	setDirectImage(image, pixbuf, maxSize, err)
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

func (win *MainWindow)onOpenComments(postId int64) {
	win.fillComments(postId)
	win.stack.SetTransitionType(gtk.STACK_TRANSITION_TYPE_SLIDE_LEFT)
	win.stack.SetVisibleChild(&win.postView.Container)
}

func (win *MainWindow)onCloseComments() {
	win.stack.SetTransitionType(gtk.STACK_TRANSITION_TYPE_SLIDE_RIGHT)
	win.stack.SetVisibleChild(&win.mainView.Container)
}

func (win *MainWindow)fillComments(postId int64) {
	post, _ := win.app.PostLemmyClient(postId)

	win.commentsView.title.SetText(post.Post.Name)
	win.commentsView.communityName.SetText(fmt.Sprintf("<span size=\"large\">%s</span>", post.Community.Title))
	win.commentsView.communityName.SetUseMarkup(true)
	win.commentsView.timestamp.SetText(fmt.Sprintf("%s ago", time.Since(post.Post.Published).Round(time.Minute).String()))

	if post.Post.Body.IsValid() {
		buffer, _ := win.commentsView.description.GetBuffer()
		buffer.SetText(post.Post.Body.ValueOrZero())
		win.commentsView.description.Show()
	} else {
		win.commentsView.description.Hide()
	}

	win.commentsView.votes.SetRange(float64(post.Counts.Score)-1, float64(post.Counts.Score)+1)
	win.commentsView.votes.SetValue(float64(post.Counts.Score))

	if post.Creator.DisplayName.IsValid() {
		win.commentsView.username.SetText(post.Creator.DisplayName.ValueOrZero())
	} else {
		win.commentsView.username.SetText(post.Creator.Name)
	}

	win.commentsView.communityIcon.Clear()
	if post.Community.Icon.IsValid() {
		LoadImageFromURL(post.Community.Icon.ValueOrZero(), func(pixbuf *gdk.Pixbuf, err error) {
			setDirectImage(win.commentsView.communityIcon, pixbuf, [2]int{communityIconSize, communityIconSize}, err)
		})
	}

	win.commentsView.image.Clear()
	if post.Post.ThumbnailURL.IsValid() {
		LoadImageFromURL(post.Post.ThumbnailURL.ValueOrZero(), func(pixbuf *gdk.Pixbuf, err error) {
			setDirectImage(win.commentsView.image, pixbuf, [2]int{maxPostImageSize, maxPostImageSize}, err)
		})
	}

	if post.Post.URL.IsValid() && (!post.Post.ThumbnailURL.IsValid() || post.Post.URL.ValueOrZero() != post.Post.ThumbnailURL.ValueOrZero()) {
		win.commentsView.link.SetUri(post.Post.URL.ValueOrZero())
		win.commentsView.link.Show()
	} else {
		win.commentsView.link.Hide()
	}
}

func NewCommentsView(builder *gtk.Builder) (commentsView CommentsView, err error) {
	card, err := getUIObject[gtk.Box](builder, "commentsContainer")
	if err != nil {
		return
	}
	convertToCard(card)

	commentsView.title, err = getUIObject[gtk.Label](builder, "title")
	if err != nil {
		return
	}

	commentsView.communityIcon, err = getUIObject[gtk.Image](builder, "communityIcon")
	if err != nil {
		return
	}

	commentsView.communityName, err = getUIObject[gtk.Label](builder, "communityName")
	if err != nil {
		return
	}

	commentsView.username, err = getUIObject[gtk.Label](builder, "username")
	if err != nil {
		return
	}

	commentsView.link, err = getUIObject[gtk.LinkButton](builder, "linkButton")
	if err != nil {
		return
	}

	commentsView.timestamp, err = getUIObject[gtk.Label](builder, "time")
	if err != nil {
		return
	}

	commentsView.image, err = getUIObject[gtk.Image](builder, "image")
	if err != nil {
		return
	}

	commentsView.description, err = getUIObject[gtk.TextView](builder, "description")
	if err != nil {
		return
	}

	commentsView.votes, err = getUIObject[gtk.SpinButton](builder, "votes")
	if err != nil {
		return
	}
	commentsView.votes.SetIncrements(1, 1)

	commentsView.comments, err = getUIObject[gtk.TextView](builder, "commentsText")
	if err != nil {
		return
	}

	return
}
