package utils

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
)

func LoadPixmapFromURL(url string, callback func(*gdk.Pixbuf, error)) {
	callbackInMainThread := func(pixbuf *gdk.Pixbuf, err error) {
		glib.IdleAdd(func() bool {
			callback(pixbuf, err)
			return false
		})
	}

	go func() {
		response, err := http.Get(url)
		if err != nil {
			callbackInMainThread(&gdk.Pixbuf{}, err)
			return
		}

		if response.StatusCode != 200 {
			callbackInMainThread(&gdk.Pixbuf{}, fmt.Errorf("%s", response.Status))
			return
		}

		data, err := io.ReadAll(response.Body)
		if err != nil {
			callbackInMainThread(&gdk.Pixbuf{}, err)
			return
		}

		loader, err := gdk.PixbufLoaderNew()
		if err != nil {
			callbackInMainThread(&gdk.Pixbuf{}, err)
			return
		}

		pixbuf, err := loader.WriteAndReturnPixbuf(data)
		callbackInMainThread(pixbuf, err)
	}()
}
