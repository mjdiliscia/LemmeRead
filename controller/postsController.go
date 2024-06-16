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

func (pc *PostsController) Init(mv *view.MainView, am *model.AppModel) {
	pc.mainView = mv
	pc.appModel = am

	mv.PostListBottomReached = pc.onPostListBottomReached
	mv.PostListView.CommentClicked = pc.onCommentsClicked
	mv.CommentsClosed = pc.onCloseCommentsClicked
	mv.OrderChanged = pc.onOrderChanged
	mv.FilterChanged = pc.onFilterChanged
}

func (pc *PostsController) onPostListBottomReached() {
	pc.appModel.RetrieveMorePosts(func(err error) {
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

func (pc *PostsController) onOrderChanged(newOrder int) {
	pc.appModel.Configuration.SetOrder(model.PostsOrder(newOrder))
	pc.mainView.CleanView()
	pc.appModel.CleanModel()
	pc.appModel.RetrieveMorePosts(func(err error) {
		if err != nil {
			log.Println(err)
		}
	})
}

func (pc *PostsController) onFilterChanged(newFilter int) {
	pc.appModel.Configuration.SetFilter(model.PostsFilter(newFilter))
	pc.mainView.CleanView()
	pc.appModel.CleanModel()
	pc.appModel.RetrieveMorePosts(func(err error) {
		if err != nil {
			log.Println(err)
		}
	})
}
