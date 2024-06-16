package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
)

var httpCache map[string]*gdkpixbuf.Pixbuf

func LoadPixmapFromUrl(url string) (pixbuf *gdkpixbuf.Pixbuf, err error) {
	if httpCache == nil {
		httpCache = make(map[string]*gdkpixbuf.Pixbuf)
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

		var loader *gdkpixbuf.PixbufLoader
		loader = gdkpixbuf.NewPixbufLoader()

		err = loader.Write(data)
		if err != nil {
			return
		}

		err = loader.Close()
		if err != nil {
			return
		}

		pixbuf = loader.Pixbuf()
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
