package view

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

type PostView struct {
	Parent *MainView
	CommentViews map[int64]CommentView
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

func (pv *PostView) SetupPostView(post model.PostModel, comments []*model.CommentModel, box *gtk.Box) (err error) {
	_, err = pv.buildAndSetReferences()
	if err != nil {
		return
	}
	pv.parentBox = box

	pv.fillPostData(post, comments == nil)
	pv.buildComments(comments)

	pv.parentBox.PackStart(pv.post, true, false, 0)

	return
}

func (pv *PostView) Destroy() {
	pv.parentBox.Remove(pv.post)
}

func (pv *PostView) buildAndSetReferences() (builder *gtk.Builder, err error) {
	builder, err = gtk.BuilderNewFromString(string(data.PostUI))
	if err != nil {
		return
	}

	pv.post, err = utils.GetUIObject[gtk.Box](builder, "post")
	if err != nil {
		return
	}

	utils.SetWidgetProperty(builder, "card", func(card *gtk.Box) {
		utils.ApplyStyle(&card.Widget)
	})

	pv.title, err = utils.GetUIObject[gtk.Label](builder, "title")
	if err != nil {
		return
	}
	utils.ApplyStyle(&pv.title.Widget)

	pv.communityIcon, err = utils.GetUIObject[gtk.Image](builder, "communityIcon")
	if err != nil {
		return
	}

	pv.communityName, err = utils.GetUIObject[gtk.Label](builder, "communityName")
	if err != nil {
		return
	}
	utils.ApplyStyle(&pv.communityName.Widget)

 	pv.username, err = utils.GetUIObject[gtk.Label](builder, "username")
	if err != nil {
		return
	}

	pv.link, err = utils.GetUIObject[gtk.LinkButton](builder, "linkButton")
	if err != nil {
		return
	}

	pv.timestamp, err = utils.GetUIObject[gtk.Label](builder, "time")
	if err != nil {
		return
	}

	pv.image, err = utils.GetUIObject[gtk.Image](builder, "image")
	if err != nil {
		return
	}

	pv.description, err = utils.GetUIObject[gtk.Label](builder, "description")
	if err != nil {
		return
	}

	pv.votes, err = utils.GetUIObject[gtk.SpinButton](builder, "votes")
	if err != nil {
		return
	}
	pv.votes.SetIncrements(1, 1)

	pv.commentsBox, err = utils.GetUIObject[gtk.Box](builder, "commentsParent")
	if err != nil {
		return
	}

	pv.commentsButton, err = utils.GetUIObject[gtk.Button](builder, "commentsButton")
	if err != nil {
		return
	}

	pv.post.Unparent()

	return
}

func (pv *PostView) fillPostData(post model.PostModel, briefDesc bool) {
	pv.title.SetText(post.Post.Name)

	if post.Post.Body.IsValid() {
		body := post.Post.Body.ValueOrZero()
		if briefDesc && len(body) > MAX_BRIEF_DESC_LEN {
			body = body[:MAX_BRIEF_DESC_LEN] + "..."
		}
		pv.description.SetMarkup(utils.MarkdownToLabelMarkup(body))
	} else {
		pv.description.Hide()
	}

	pv.communityName.SetText(post.Community.Title)
	pv.username.SetText(post.Creator.DisplayName.ValueOr(post.Creator.Name))
	pv.timestamp.SetText(utils.GetNiceDuration(time.Since(post.Post.Published)))

	pv.votes.SetRange(float64(post.Counts.Score)-1, float64(post.Counts.Score)+1)
	pv.votes.SetValue(float64(post.Counts.Score))

	if briefDesc {
		pv.commentsButton.SetLabel(fmt.Sprintf("%d comments", post.Counts.Comments))
		pv.commentsButton.Connect("clicked", func() {
			if pv.CommentsButtonClicked != nil {
				pv.CommentsButtonClicked(post.Post.ID)
			}
		})
	} else {
		pv.commentsButton.Hide()
	}

	if !post.IsImagePost && post.Link != "" {
		pv.link.SetUri(post.Link)
		pv.link.Show()
	}

	if post.Image != nil {
		utils.SetDirectImage(pv.image, post.Image, [2]int{maxPostImageSize, maxPostImageSize}, nil)
	}

	if post.CommunityIcon != nil {
		utils.SetDirectImage(pv.communityIcon, post.CommunityIcon, [2]int{communityIconSize, communityIconSize}, nil)
	}
}

func (pv *PostView) buildComments(inComments []*model.CommentModel) {
	addCommentsTo(pv.commentsBox, inComments)

	newImage, _ := gtk.ImageNew()
	pv.commentsBox.PackStart(newImage, true, true, 0)
	newImage.SetVExpand(true)
}

func addCommentsTo(box *gtk.Box, comments []*model.CommentModel) {
	for _, comment := range(comments) {
		log.Printf("Adding comment %d", comment.Comment.ID)
		commentView, err := NewCommentView(*comment)
		if err != nil {
			log.Printf("Error creating comment UI for %d", comment.Comment.ID)
			return
		}
		box.PackStart(commentView.CommentBox, true, false, 5)
		if len(comment.ChildComments) > 0 {
			log.Printf("%d has %d children", comment.Comment.ID, len(comment.ChildComments))
			addCommentsTo(commentView.childCommentsBox, comment.ChildComments)
		}
	}
}
