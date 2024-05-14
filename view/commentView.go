package view

import (
	"time"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/data"
	"github.com/mjdiliscia/LemmeRead/model"
	"github.com/mjdiliscia/LemmeRead/utils"
)

type CommentView struct {
	CommentBox   *gtk.Box
	VotesChanged func(int64, int64)

	commentID        int64
	username         *gtk.Label
	timestamp        *gtk.Label
	commentText      *gtk.Label
	votes            *gtk.SpinButton
	userImage        *gtk.Image
	foldButton       *gtk.Button
	unfoldButton     *gtk.Button
	childCommentsBox *gtk.Box
}

func NewCommentView(comment model.CommentModel) (cv CommentView, err error) {
	_, err = cv.buildAndSetReferences()
	if err != nil {
		return
	}

	cv.fillCommentData(comment)

	return
}

func (cv *CommentView) AddChildComment(commentView CommentView) {
	cv.childCommentsBox.PackStart(commentView.CommentBox, true, false, 0)
}

func (cv *CommentView) buildAndSetReferences() (commentBox *gtk.Box, err error) {
	builder, err := gtk.BuilderNewFromString(string(data.CommentUI))

	if err != nil {
		return
	}

	cv.username, err = utils.GetUIObject[gtk.Label](builder, "username")
	if err != nil {
		return
	}

	cv.timestamp, err = utils.GetUIObject[gtk.Label](builder, "timestamp")
	if err != nil {
		return
	}

	cv.commentText, err = utils.GetUIObject[gtk.Label](builder, "commentText")
	if err != nil {
		return
	}

	cv.votes, err = utils.GetUIObject[gtk.SpinButton](builder, "votes")
	if err != nil {
		return
	}
	cv.votes.SetIncrements(1, 1)
	cv.votes.Connect("value-changed", func() {
		if cv.VotesChanged != nil {
			cv.VotesChanged(cv.commentID, int64(cv.votes.GetValue()))
		}
	})

	cv.userImage, err = utils.GetUIObject[gtk.Image](builder, "userImage")
	if err != nil {
		return
	}

	cv.foldButton, err = utils.GetUIObject[gtk.Button](builder, "fold")
	if err != nil {
		return
	}
	cv.foldButton.Connect("clicked", func() {
		cv.foldButton.Hide()
		cv.unfoldButton.Show()
		cv.commentText.Hide()
		cv.votes.Hide()
		cv.childCommentsBox.Hide()
	})

	cv.unfoldButton, err = utils.GetUIObject[gtk.Button](builder, "unfold")
	if err != nil {
		return
	}
	cv.unfoldButton.Connect("clicked", func() {
		cv.foldButton.Show()
		cv.unfoldButton.Hide()
		cv.commentText.Show()
		cv.votes.Show()
		cv.childCommentsBox.Show()
	})

	cv.childCommentsBox, err = utils.GetUIObject[gtk.Box](builder, "children")

	cv.CommentBox, err = utils.GetUIObject[gtk.Box](builder, "commentBox")
	if err != nil {
		return
	}
	cv.CommentBox.Unparent()

	return
}

func (cv *CommentView) fillCommentData(comment model.CommentModel) {
	cv.commentID = comment.Comment.ID
	cv.username.SetText(comment.Creator.DisplayName.ValueOr(comment.Creator.Name))
	cv.timestamp.SetText(utils.GetNiceDuration(time.Since(comment.Comment.Published)))

	cv.commentText.SetMarkup(utils.MarkdownToLabelMarkup(comment.Comment.Content))

	cv.votes.SetRange(float64(comment.Counts.Score)-1, float64(comment.Counts.Score)+1)
	cv.votes.SetValue(float64(comment.Counts.Score))

	if comment.Creator.Avatar.IsValid() {
		var taskSequence *utils.TaskSequence[*gdk.Pixbuf]
		taskSequence = utils.NewTaskSequence[*gdk.Pixbuf](func() {
			taskSequence = nil
		})

		taskSequence.Add(func() (*gdk.Pixbuf, error) {
			return utils.LoadPixmapFromUrl(comment.Creator.Avatar.ValueOrZero())
		}, func(pixbuf *gdk.Pixbuf, err error) bool {
			utils.SetDirectImage(cv.userImage, pixbuf, [2]int{communityIconSize, communityIconSize}, err)
			return true
		})
	}
}
