package ui

import (
	"log"
	"slices"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/model"
)

type PostsUI struct {
	CommentButtonClicked func(int64)

	postsBox *gtk.Box
	posts []PostUI
	shownPosts []int64
}

func (pui *PostsUI) SetupPostsUI(box *gtk.Box) (err error) {
	pui.postsBox = box

	return
}

func (pui *PostsUI) FillPostsData(posts []model.PostModel) {
	for _, post := range posts {
		if slices.Index(pui.shownPosts, post.Post.ID) != -1 {
			log.Printf("Post %d already being shown, skipping.", post.Post.ID)
			continue
		}

		log.Printf("Adding post %d to PostsUI...", post.Post.ID)
		pui.shownPosts = append(pui.shownPosts, post.Post.ID)
		pui.posts = append(pui.posts, PostUI{})
		postUI := pui.posts[len(pui.posts)-1]
		err := postUI.SetupPostUI(post, nil, pui.postsBox)
		if err != nil {
			log.Println(err)
		}

		postUI.CommentsButtonClicked = func(id int64) {
			if pui.CommentButtonClicked != nil {
				pui.CommentButtonClicked(id)
			}
		}
		log.Printf("Added post %d to PostUI.", post.Post.ID)
	}
}
