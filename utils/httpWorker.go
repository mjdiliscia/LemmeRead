package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gotk3/gotk3/gdk"
)

var httpCache map[string]*gdk.Pixbuf

func LoadPixmapFromUrl(url string) (pixbuf *gdk.Pixbuf, err error) {
	if httpCache == nil {
		httpCache = make(map[string]*gdk.Pixbuf)
	}

	timestamp := time.Now()
	pixbuf, ok := httpCache[url]
	if !ok {
		var response *http.Response
		response, err = http.Get(url)
		if err != nil {
			return
		}

		if response.StatusCode != 200 {
			return nil, fmt.Errorf("%s", response.Status)
		}

		var data []byte
		data, err = io.ReadAll(response.Body)
		if err != nil {
			return
		}

		var loader *gdk.PixbufLoader
		loader, err = gdk.PixbufLoaderNew()
		if err != nil {
			return
		}

		pixbuf, err = loader.WriteAndReturnPixbuf(data)
		log.Printf("GET time for '%s': %d", url, time.Now().Sub(timestamp).Milliseconds())

		httpCache[url] = pixbuf
	} else {
		log.Printf("CACHE time for '%s': %d", url, time.Now().Sub(timestamp).Milliseconds())
	}

	return
}

func GetUrlMimetype(url string) (string, error) {
	response, err := http.Head(url)

	if response.StatusCode != 200 {
		return "", fmt.Errorf("%s", response.Status)
	}

	return response.Header.Get("Content-Type"), err
}
