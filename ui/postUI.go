package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/data"
	"github.com/mjdiliscia/LemmeRead/model"
	"github.com/mjdiliscia/LemmeRead/utils"
)

const MAX_BRIEF_DESC_LEN int = 500;

type PostUI struct {
	Parent *MainWindow
	CommentUIs map[int64]CommentUI
	CommentsButtonClicked func(int64)

	parentBox *gtk.Box
	post *gtk.Box
	title *gtk.Label
	communityIcon *gtk.Image
	communityName *gtk.Label
	username *gtk.Label
	timestamp *gtk.Label
	link *gtk.LinkButton
	image *gtk.Image
	description *gtk.Label
	votes *gtk.SpinButton
	commentsBox *gtk.Box
	commentsButton *gtk.Button
}

func (pui *PostUI) SetupPostUI(post model.PostModel, comments []*model.CommentModel, box *gtk.Box) (err error) {
	_, err = pui.buildAndSetReferences()
	if err != nil {
		return
	}
	pui.parentBox = box

	pui.fillPostData(post, comments == nil)
	pui.buildComments(comments)

	pui.parentBox.PackStart(pui.post, true, false, 0)

	return
}

func (pui *PostUI) Destroy() {
	pui.parentBox.Remove(pui.post)
}

func (pui *PostUI) buildAndSetReferences() (builder *gtk.Builder, err error) {
	builder, err = gtk.BuilderNewFromString(string(data.PostUI))
	if err != nil {
		return
	}

	pui.post, err = utils.GetUIObject[gtk.Box](builder, "post")
	if err != nil {
		return
	}

	utils.SetWidgetProperty(builder, "card", func(card *gtk.Box) {
		utils.ApplyStyle(&card.Widget)
	})

	pui.title, err = utils.GetUIObject[gtk.Label](builder, "title")
	if err != nil {
		return
	}
	utils.ApplyStyle(&pui.title.Widget)

	pui.communityIcon, err = utils.GetUIObject[gtk.Image](builder, "communityIcon")
	if err != nil {
		return
	}

	pui.communityName, err = utils.GetUIObject[gtk.Label](builder, "communityName")
	if err != nil {
		return
	}
	utils.ApplyStyle(&pui.communityName.Widget)

 	pui.username, err = utils.GetUIObject[gtk.Label](builder, "username")
	if err != nil {
		return
	}

	pui.link, err = utils.GetUIObject[gtk.LinkButton](builder, "linkButton")
	if err != nil {
		return
	}

	pui.timestamp, err = utils.GetUIObject[gtk.Label](builder, "time")
	if err != nil {
		return
	}

	pui.image, err = utils.GetUIObject[gtk.Image](builder, "image")
	if err != nil {
		return
	}

	pui.description, err = utils.GetUIObject[gtk.Label](builder, "description")
	if err != nil {
		return
	}

	pui.votes, err = utils.GetUIObject[gtk.SpinButton](builder, "votes")
	if err != nil {
		return
	}
	pui.votes.SetIncrements(1, 1)

	pui.commentsBox, err = utils.GetUIObject[gtk.Box](builder, "commentsParent")
	if err != nil {
		return
	}

	pui.commentsButton, err = utils.GetUIObject[gtk.Button](builder, "commentsButton")
	if err != nil {
		return
	}

	pui.post.Unparent()

	return
}

func (pui *PostUI) fillPostData(post model.PostModel, briefDesc bool) {
	pui.title.SetText(post.Post.Name)

	if post.Post.Body.IsValid() {
		body := post.Post.Body.ValueOrZero()
		if briefDesc && len(body) > MAX_BRIEF_DESC_LEN {
			body = body[:MAX_BRIEF_DESC_LEN] + "..."
		}
		pui.description.SetMarkup(utils.MarkdownToLabelMarkup(body))
	} else {
		pui.description.Hide()
	}

	pui.communityName.SetText(post.Community.Title)
	pui.username.SetText(post.Creator.DisplayName.ValueOr(post.Creator.Name))
	pui.timestamp.SetText(utils.GetNiceDuration(time.Since(post.Post.Published)))

	pui.votes.SetRange(float64(post.Counts.Score)-1, float64(post.Counts.Score)+1)
	pui.votes.SetValue(float64(post.Counts.Score))

	if briefDesc {
		pui.commentsButton.SetLabel(fmt.Sprintf("%d comments", post.Counts.Comments))
		pui.commentsButton.Connect("clicked", func() {
			if pui.CommentsButtonClicked != nil {
				pui.CommentsButtonClicked(post.Post.ID)
			}
		})
	} else {
		pui.commentsButton.Hide()
	}

	urlIsThumbURL := post.Post.ThumbnailURL.IsValid() && post.Post.URL.ValueOrZero() == post.Post.ThumbnailURL.ValueOrZero()
	if post.Post.URL.IsValid() && !urlIsThumbURL {
		pui.link.SetUri(post.Post.URL.ValueOrZero())
		pui.link.Show()
	}

	if post.Image != nil {
		utils.SetDirectImage(pui.image, post.Image, [2]int{maxPostImageSize, maxPostImageSize}, nil)
	}

	if post.CommunityIcon != nil {
		utils.SetDirectImage(pui.communityIcon, post.CommunityIcon, [2]int{communityIconSize, communityIconSize}, nil)
	}
}

func (pui *PostUI) buildComments(inComments []*model.CommentModel) {
	addCommentsTo(pui.commentsBox, inComments)

	newImage, _ := gtk.ImageNew()
	pui.commentsBox.PackStart(newImage, true, true, 0)
	newImage.SetVExpand(true)
}

func addCommentsTo(box *gtk.Box, comments []*model.CommentModel) {
	for _, comment := range(comments) {
		log.Printf("Adding comment %d", comment.Comment.ID)
		commentUI, err := NewCommentUI(*comment)
		if err != nil {
			log.Printf("Error creating comment UI for %d", comment.Comment.ID)
			return
		}
		box.PackStart(commentUI.CommentBox, true, false, 5)
		if len(comment.ChildComments) > 0 {
			log.Printf("%d has %d children", comment.Comment.ID, len(comment.ChildComments))
			addCommentsTo(commentUI.childCommentsBox, comment.ChildComments)
		}
	}
}
