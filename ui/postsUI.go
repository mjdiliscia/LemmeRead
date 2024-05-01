package ui

import (
	"log"
	"slices"

	"github.com/gotk3/gotk3/gtk"
	"go.elara.ws/go-lemmy"
)

type PostsUI struct {
	CommentButtonClicked func(int64)

	postsBox *gtk.Box
	posts []PostUI
	shownPosts []int64
}

func NewPostsUI(box *gtk.Box) (pui PostsUI, err error) {
	pui.postsBox = box

	return
}

func (pui *PostsUI) FillPostsData(posts []lemmy.PostView) {
	for _, post := range posts {
		if slices.Index(pui.shownPosts, post.Post.ID) != -1 {
			continue
		}

		pui.shownPosts = append(pui.shownPosts, post.Post.ID)
		postUI, err := NewPostUI(post, nil, pui.postsBox)
		if err != nil {
			log.Println(err)
		}
		pui.posts = append(pui.posts, postUI)

		postUI.CommentsButtonClicked = func(id int64) {
			if pui.CommentButtonClicked != nil {
				pui.CommentButtonClicked(id)
			}
		}
	}
}
