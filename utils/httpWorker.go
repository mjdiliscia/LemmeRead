package utils

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gotk3/gotk3/gdk"
)

func LoadPixmapFromUrl(url string) (pixbuf *gdk.Pixbuf, err error) {
	response, err := http.Get(url)
	if err != nil {
		return
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("%s", response.Status)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}

	loader, err := gdk.PixbufLoaderNew()
	if err != nil {
		return
	}

	pixbuf, err = loader.WriteAndReturnPixbuf(data)
	return
}

func GetUrlMimetype(url string) (string, error) {
	response, err := http.Head(url)

	if response.StatusCode != 200 {
		return "", fmt.Errorf("%s", response.Status)
	}

	return response.Header.Get("Content-Type"), err
}
