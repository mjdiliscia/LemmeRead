package utils

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"regexp"
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/data"
)

func GetUIObject[OType any](builder *gtk.Builder, objectId string) (object *OType, err error) {
	obj, err := builder.GetObject(objectId)
	if err != nil {
		err = fmt.Errorf("Couldn't find object of name '%s' (asked type was %s): %s", objectId, reflect.TypeOf(object).Name(), err)
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
	switch {
	case timestamp.Hours() > 24*365:
		return fmt.Sprintf("%dy ago", int(timestamp.Hours()/24/356))
	case timestamp.Hours() > 24*30:
		return fmt.Sprintf("%dm ago", int(timestamp.Hours()/24/30))
	case timestamp.Hours() > 24*7:
		return fmt.Sprintf("%dw ago", int(timestamp.Hours()/24/7))
	case timestamp.Hours() > 24:
		return fmt.Sprintf("%dd ago", int(timestamp.Hours()/24))
	case timestamp.Hours() > 1:
		return fmt.Sprintf("%dh ago", int(timestamp.Hours()))
	case timestamp.Minutes() > 1:
		return fmt.Sprintf("%dmin ago", int(timestamp.Minutes()))
	default:
		return fmt.Sprintf("%ds", int(timestamp.Seconds()))
	}
}

func MarkdownToLabelMarkup(text string) (markup string) {
    // Bold conversion
    boldRe := regexp.MustCompile(`\*\*(.+?)\*\*`)
    markup = boldRe.ReplaceAllString(text, fmt.Sprintf("<b>%s</b>", "$1"))

    // Italic conversion
    italicRe := regexp.MustCompile(`\_(.+?)\_`)
    markup = italicRe.ReplaceAllString(markup, fmt.Sprintf("<i>%s</i>", "$1"))

    // Strikethrough conversion (basic)
    strikethroughRe := regexp.MustCompile(`~~(.+?)~~`)
    markup = strikethroughRe.ReplaceAllString(markup, fmt.Sprintf("<span style=\"text-decoration: line-through\">%s</span>", "$1"))

	// Link conversion
    linkRe := regexp.MustCompile(`\[(.+?)\]\((.+?)\)`)
	markup = linkRe.ReplaceAllString(markup, fmt.Sprintf("<a href=\"%s\">%s</a>", "$2", "$1"))

	// Textless link conversion
    textlessLinkRe := regexp.MustCompile(`!\[\]\((.+?)\)`)
	markup = textlessLinkRe.ReplaceAllString(markup, fmt.Sprintf("<a href=\"%s\">%s</a>", "$1", "$1"))

	// & correction
	ampersandRe := regexp.MustCompile(`\&`)
	markup = ampersandRe.ReplaceAllString(markup, "&amp;")

	return
}
