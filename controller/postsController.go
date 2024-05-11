package controller

import (
	"log"

	"github.com/mjdiliscia/LemmeRead/model"
	"github.com/mjdiliscia/LemmeRead/ui"
)

type PostsController struct {
	mainWindow *ui.MainWindow
	appModel *model.AppModel
}

func (pc *PostsController) Init(mw *ui.MainWindow, am *model.AppModel) {
	pc.mainWindow = mw
	pc.appModel = am

	mw.PostListBottomReached = pc.onPostListBottomReached
	mw.PostList.CommentClicked = pc.onCommentsClicked
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
		pc.mainWindow.OpenComments(id)
	})
}

func (pc *PostsController) onCloseCommentsClicked() {
	pc.mainWindow.CloseComments()
}
