package view

import (
	"log"
	"slices"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/model"
)

type PostListView struct {
	CommentClicked func(int64)

	postsBox   *gtk.Box
	postViews  []PostView
	shownPosts []int64
}

func (plv *PostListView) SetupPostListView(box *gtk.Box) (err error) {
	plv.postsBox = box

	return
}

func (plv *PostListView) FillPostsData(posts []model.PostModel) {
	for _, post := range posts {
		if slices.Index(plv.shownPosts, post.Post.ID) != -1 {
			log.Printf("Post %d already being shown, skipping.", post.Post.ID)
			continue
		}

		log.Printf("Adding post %d to PostsUI...", post.Post.ID)
		plv.shownPosts = append(plv.shownPosts, post.Post.ID)
		plv.postViews = append(plv.postViews, PostView{})
		postView := plv.postViews[len(plv.postViews)-1]
		err := postView.SetupPostView(post, nil, plv.postsBox)
		if err != nil {
			log.Println(err)
		}

		postView.CommentsButtonClicked = func(id int64) {
			if plv.CommentClicked != nil {
				plv.CommentClicked(id)
			}
		}
		log.Printf("Added post %d to PostUI.", post.Post.ID)
	}
}
