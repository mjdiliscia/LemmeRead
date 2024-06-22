package model

import (
	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	"go.elara.ws/go-lemmy"
)

type CommentModel struct {
	lemmy.CommentView
	UserIcon      gdkpixbuf.Pixbuf
	ChildComments []*CommentModel
}
