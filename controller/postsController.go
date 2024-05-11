package controller

import (
	"log"

	"github.com/mjdiliscia/LemmeRead/model"
	"github.com/mjdiliscia/LemmeRead/view"
)

type PostsController struct {
	mainView *view.MainView
	appModel *model.AppModel
}

func (pc *PostsController) Init(mw *view.MainView, am *model.AppModel) {
	pc.mainView = mw
	pc.appModel = am

	mw.PostListBottomReached = pc.onPostListBottomReached
	mw.PostListView.CommentClicked = pc.onCommentsClicked
	mw.CloseCommentsClicked = pc.onCloseCommentsClicked
}

func (pc *PostsController) onPostListBottomReached() {
	pc.appModel.RetrieveMorePosts(func (err error) {
		if err != nil {
			log.Println(err)
		}
	})
}

func (pc *PostsController) onCommentsClicked(id int64) {
	pc.appModel.RetrieveComments(id, func(err error) {
		if err != nil {
			log.Println(err)
			return
		}
		pc.mainView.OpenComments(id)
	})
}

func (pc *PostsController) onCloseCommentsClicked() {
	pc.mainView.CloseComments()
}
