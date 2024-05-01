package ui

import (
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/data"
	"github.com/mjdiliscia/LemmeRead/utils"
	"go.elara.ws/go-lemmy"
)

type CommentUI struct {
	CommentBox *gtk.Box
	VotesChanged func(int64, int64)

	commentID int64
	username *gtk.Label
	timestamp *gtk.Label
	commentText *gtk.Label
	votes *gtk.SpinButton
	userImage *gtk.Image
	foldButton *gtk.Button
	unfoldButton *gtk.Button
	childCommentsBox *gtk.Box
}

func NewCommentUI(comment lemmy.CommentView) (cui CommentUI, err error) {
	_, err = cui.buildAndSetReferences()
	if err != nil {
		return
	}

	cui.fillCommentData(comment)

	return
}

func (cui *CommentUI) AddChildComment(commentUI CommentUI) {
	cui.childCommentsBox.PackStart(commentUI.CommentBox, true, false, 0)
}

func (cui *CommentUI) buildAndSetReferences() (commentBox *gtk.Box, err error) {
	builder, err := gtk.BuilderNewFromString(string(data.CommentUI))

	if err != nil {
		return
	}

	cui.username, err = utils.GetUIObject[gtk.Label](builder, "username")
	if err != nil {
		return
	}

	cui.timestamp, err = utils.GetUIObject[gtk.Label](builder, "timestamp")
	if err != nil {
		return
	}

	cui.commentText, err = utils.GetUIObject[gtk.Label](builder, "commentText")
	if err != nil {
		return
	}

	cui.votes, err = utils.GetUIObject[gtk.SpinButton](builder, "votes")
	if err != nil {
		return
	}
	cui.votes.SetIncrements(1, 1)
	cui.votes.Connect("value-changed", func() {
		if cui.VotesChanged != nil {
			cui.VotesChanged(cui.commentID, int64(cui.votes.GetValue()))
		}
	})

	cui.userImage, err = utils.GetUIObject[gtk.Image](builder, "userImage")
	if err != nil {
		return
	}

	cui.foldButton, err = utils.GetUIObject[gtk.Button](builder, "fold")
	if err != nil {
		return
	}
	cui.foldButton.Connect("clicked", func() {
		cui.foldButton.Hide()
		cui.unfoldButton.Show()
		cui.commentText.Hide()
		cui.votes.Hide()
		cui.childCommentsBox.Hide()
	})

	cui.unfoldButton, err = utils.GetUIObject[gtk.Button](builder, "unfold")
	if err != nil {
		return
	}
	cui.unfoldButton.Connect("clicked", func() {
		cui.foldButton.Show()
		cui.unfoldButton.Hide()
		cui.commentText.Show()
		cui.votes.Show()
		cui.childCommentsBox.Show()
	})

	cui.childCommentsBox, err = utils.GetUIObject[gtk.Box](builder, "children")

	cui.CommentBox, err = utils.GetUIObject[gtk.Box](builder, "commentBox")
	if err != nil {
		return
	}
	cui.CommentBox.Unparent()

	return
}

func (cui *CommentUI) fillCommentData(comment lemmy.CommentView) {
	cui.commentID = comment.Comment.ID
	cui.username.SetText(comment.Creator.DisplayName.ValueOr(comment.Creator.Name))
	cui.timestamp.SetText(utils.GetNiceDuration(time.Since(comment.Comment.Published)))

	cui.commentText.SetMarkup(utils.MarkdownToLabelMarkup(comment.Comment.Content))

	cui.votes.SetRange(float64(comment.Counts.Score)-1, float64(comment.Counts.Score)+1)
	cui.votes.SetValue(float64(comment.Counts.Score))

	if comment.Creator.Avatar.IsValid() {
		utils.LoadPixmapFromURL(comment.Creator.Avatar.ValueOrZero(), func(pixbuf *gdk.Pixbuf, err error) {
			utils.SetDirectImage(cui.userImage, pixbuf, [2]int{communityIconSize, communityIconSize}, err)
		})
	}
}
