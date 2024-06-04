package model

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gdkpixbuf/v2"
	"github.com/mjdiliscia/LemmeRead/utils"
	"go.elara.ws/go-lemmy"
)

type PostModel struct {
	lemmy.PostView
	IsImagePost   bool
	Link          string
	Image         *gdkpixbuf.Pixbuf
	CommunityIcon *gdkpixbuf.Pixbuf
	Comments      []*CommentModel

	commentHolder []CommentModel
}

type PMData struct {
	str    string
	pixbuf *gdkpixbuf.Pixbuf
}

func (pm *PostModel) Init(callback func(error)) {
	var taskSequence *utils.TaskSequence[PMData]
	taskSequence = utils.NewTaskSequence[PMData](func() {
		taskSequence = nil
		callback(nil)
	})

	taskSequence.Add(pm.getMimetypeTask, pm.processMimetypeTask)
	taskSequence.Add(pm.getPostImageTask(), pm.setImageTask)
	taskSequence.Add(pm.getPixbufTask(pm.Community.Icon), pm.setCommunityIconTask)

	taskSequence.Execute()
}

func (pm *PostModel) AddComments(comments []lemmy.CommentView, err error) error {
	if err != nil {
		return err
	}
	slices.SortFunc(comments, func(a lemmy.CommentView, b lemmy.CommentView) int {
		return len(strings.Split(a.Comment.Path, ".")) - len(strings.Split(b.Comment.Path, "."))
	})

	commentMap := make(map[string]*CommentModel, len(comments))
	for _, comment := range comments {
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

func (pm *PostModel) getMimetypeTask() (PMData, error) {
	if pm.Post.URL.IsValid() {
		mimetype, err := utils.GetUrlMimetype(pm.Post.URL.ValueOrZero())
		if err != nil {
			return PMData{}, err
		}
		return PMData{str: mimetype}, err
	} else {
		return PMData{}, nil
	}
}

func (pm *PostModel) processMimetypeTask(data PMData, err error) bool {
	pm.IsImagePost = strings.Split(data.str, "/")[0] == "image"
	if !pm.IsImagePost {
		pm.Link = pm.Post.URL.ValueOrZero()
	}
	return true
}

func (pm *PostModel) getPostImageTask() func() (PMData, error) {
	var url lemmy.Optional[string]
	if pm.IsImagePost {
		url = pm.Post.URL
	} else {
		url = pm.Post.ThumbnailURL
	}

	return pm.getPixbufTask(url)
}

func (pm *PostModel) getPixbufTask(url lemmy.Optional[string]) func() (PMData, error) {
	return func() (PMData, error) {
		if url.IsValid() {
			pixbuf, err := utils.LoadPixmapFromUrl(url.ValueOrZero())
			pmdata := PMData{pixbuf: pixbuf}
			return pmdata, err
		} else {
			return PMData{}, nil
		}
	}
}

func (pm *PostModel) setImageTask(data PMData, err error) bool {
	if err != nil {
		log.Println(err)
	} else {
		pm.Image = data.pixbuf
	}
	return true
}

func (pm *PostModel) setCommunityIconTask(data PMData, err error) bool {
	if err != nil {
		log.Println(err)
	} else {
		pm.CommunityIcon = data.pixbuf
	}
	return true
}
