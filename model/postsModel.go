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
	Comments []*CommentModel

	commentHolder []CommentModel
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
		log.Printf("Processing comment with path %s...", comment.Comment.Path)
		if slices.IndexFunc(pm.commentHolder, func(c CommentModel) bool { return comment.Comment.ID == c.Comment.ID }) >= 0 {
			log.Printf("Comment %d already known, skipping.", comment.Comment.ID)
			continue
		}

		pm.commentHolder = append(pm.commentHolder, CommentModel{CommentView: comment})
		commentPtr := &pm.commentHolder[len(pm.commentHolder)-1]
		commentMap[comment.Comment.Path] = commentPtr

		parent := strings.Replace(comment.Comment.Path, fmt.Sprintf(".%d", comment.Comment.ID), "", 1)
		if parent == "0" {
			log.Println("is root comment.")
			pm.Comments = append(pm.Comments, commentPtr)
		} else if parentComment, ok := commentMap[parent]; ok {
			log.Println("is child comment.")
			parentComment.ChildComments = append(parentComment.ChildComments, commentPtr)
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
