package model

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/mjdiliscia/LemmeRead/utils"
	"go.elara.ws/go-lemmy"
)

type PostModel struct {
	lemmy.PostView
	Image *gdk.Pixbuf
	CommunityIcon *gdk.Pixbuf
	Comments []CommentModel
}

func (pm *PostModel) Init(callback func(error)) {
	pm.loadPixmapAndContinue(pm.Post.ThumbnailURL, func(pixbuf *gdk.Pixbuf) {
		pm.Image = pixbuf
	}, func() {
		pm.loadPixmapAndContinue(pm.Community.Icon, func(pixbuf *gdk.Pixbuf) {
			pm.CommunityIcon = pixbuf
		}, func() {
			callback(nil)
		})
	})
}

func (pm *PostModel) AddComments(comments []lemmy.CommentView, err error) error {
	if err != nil {
		return err
	}
	slices.SortFunc(comments, func(a lemmy.CommentView, b lemmy.CommentView) int {
		return len(strings.Split(a.Comment.Path, ".")) - len(strings.Split(b.Comment.Path, "."))
	})

	commentMap := make(map[string]*CommentModel, len(comments))
	for _, comment := range(comments) {
		parent := strings.Replace(comment.Comment.Path, fmt.Sprintf(".%d", comment.Comment.ID), "", 1)
		if parent == "0" {
			pm.Comments = append(pm.Comments, CommentModel{CommentView: comment})
			commentMap[comment.Comment.Path] = &pm.Comments[len(pm.Comments)-1]
		} else if parentComment, ok := commentMap[parent]; ok {
			parentComment.ChildComments = append(parentComment.ChildComments, CommentModel{CommentView: comment})
		} else {
			log.Printf("Couldn't find %s", parent)
		}
	}

	return nil
}

func (pm *PostModel) loadPixmapAndContinue(url lemmy.Optional[string], apply func(pixbuf *gdk.Pixbuf), next func()) {
	if url.IsValid() {
		utils.LoadPixmapFromURL(url.ValueOrZero(), func(pb *gdk.Pixbuf, err error) {
			if err != nil {
				log.Println(err)
			} else {
				apply(pb)
			}
			next()
		})
	} else {
		next()
	}
}
