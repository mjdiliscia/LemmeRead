package main

import (
	"fmt"
	"log"
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
	setWidgetProperty(builder, "time", func(label *gtk.Label) { label.SetText(fmt.Sprintf("%s ago", time.Since(post.Post.Published).Round(time.Minute).String())) })
	setWidgetProperty(builder, "votes", func(spinner *gtk.SpinButton) {
		spinner.SetValue(float64(post.Counts.Score))
		spinner.SetRange(spinner.GetValue()-1, spinner.GetValue()+1)
		spinner.SetIncrements(1, 1)
	})
	setWidgetProperty(builder, "comments", func(button *gtk.Button) { button.SetLabel(fmt.Sprintf("%d comments", post.Counts.Comments)) })

	postUI, err = getUIObject[gtk.Box](builder, "post")
	if err != nil {
		return
	}

	/*if post.Post.URL.IsValid() {
		stringURL := post.Post.URL.ValueOrZero()
		res, err := http.Get(stringURL)
		if err != nil {
			log.Printf("Error downloading image '%s': %s", stringURL, err)
		} else {
			defer res.Body.Close()
			log.Printf("%s\n%s:%d", post.Post.Name, stringURL, res.ContentLength)

			image, err := getUIObject[gtk.Image](builder, "image")
			if err != nil {
				log.Println(err)
			} else {
				partialData := make([]byte, 128)
				data := make([]byte, 0)
				for ; err == nil; {
					_, err = res.Body.Read(partialData)
					data = append(data, partialData...)
				}
				if err != nil && err.Error() != "EOF" {
					log.Panicf("Error reading response body: %s", err)
				}
				loader, err := gdk.PixbufLoaderNew()
				if err != nil {
					log.Panic(err)
				}
				pixbuf, err := loader.WriteAndReturnPixbuf(data)
				if err != nil {
					log.Panicf("Couldn't load image: %s", err)
				}
				image.SetFromPixbuf(pixbuf)
			}
		}
	}*/
	return
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
