package model

import (
	"github.com/gotk3/gotk3/gdk"
	"go.elara.ws/go-lemmy"
)

type CommentModel struct {
	lemmy.CommentView
	UserIcon      gdk.Pixbuf
	ChildComments []*CommentModel
}
