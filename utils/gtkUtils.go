package utils

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/data"
)

func GetUIObject[OType any](builder *gtk.Builder, objectId string) (object *OType, err error) {
	obj, err := builder.GetObject(objectId)
	if err != nil {
		fmt.Errorf("Couldn't find object of name '%s' (asked type was %s): %s", objectId, reflect.TypeOf(object).Name(), err)
		return
	}
	object, ok := any(obj).(*OType)
	if ok {
		return object, nil
	} else {
		return object, fmt.Errorf("Object '%s' can't be correctly casted to %s", objectId, reflect.TypeOf(object).Name())
	}
}

func SetDirectImage(image *gtk.Image, pixbuf *gdk.Pixbuf, maxSize [2]int, err error) {
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

func SetImage(builder *gtk.Builder, pixbuf *gdk.Pixbuf, imageId string, maxSize [2]int, err error) {
	if err != nil {
		return
	}

	image, err := GetUIObject[gtk.Image](builder, imageId)
	if err != nil {
		log.Println(err)
		return
	}

	SetDirectImage(image, pixbuf, maxSize, err)
}

func SetWidgetProperty[WType any](builder *gtk.Builder, widgetId string, setter func(widget *WType)) (err error) {
	widget, err := GetUIObject[WType](builder, widgetId)
	if err != nil {
		return
	}
	setter(widget)
	return
}

func ApplyStyle(widget *gtk.Widget) {
	cssProvider, _ := gtk.CssProviderNew()
	cssProvider.LoadFromData(string(data.StyleCSS))
	context, _ := widget.GetStyleContext()
	context.AddProvider(cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func GetNiceDuration(timestamp time.Duration) string {
	return fmt.Sprintf("%s ago", timestamp.String())
}
