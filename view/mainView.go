package view

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/data"
	"github.com/mjdiliscia/LemmeRead/model"
	"github.com/mjdiliscia/LemmeRead/utils"
)

const (
	applicationTitle  = "Lemme Read"
	maxPostImageSize  = 580
	communityIconSize = 30
)

type MainView struct {
	Window                *gtk.ApplicationWindow
	Model                 *model.AppModel
	PostListView          PostListView
	PostView              *PostView
	PostListBottomReached func()
	CloseCommentsClicked  func()

	stack          *gtk.Stack
	postListBox    *gtk.Box
	postListScroll *gtk.ScrolledWindow
	postBox        *gtk.Box
	postScroll     *gtk.ScrolledWindow
	closeComments  *gtk.Button
}

func (mv *MainView) SetupMainView(appModel *model.AppModel) (err error) {
	mv.Model = appModel
	mv.Model.NewPosts = mv.onNewPosts

	_, err = mv.buildAndSetReferences()
	if err != nil {
		return
	}

	err = mv.PostListView.SetupPostListView(mv.postListBox)
	if err != nil {
		return
	}

	mv.postListScroll.Connect("edge-reached", func(scroll *gtk.ScrolledWindow, position gtk.PositionType) {
		if position == gtk.POS_BOTTOM && mv.PostListBottomReached != nil {
			mv.PostListBottomReached()
		}
	})

	mv.closeComments.Connect("clicked", func() {
		if mv.CloseCommentsClicked != nil {
			mv.CloseCommentsClicked()
		}
	})

	mv.Window.Show()

	return nil
}

func (mv *MainView) buildAndSetReferences() (builder *gtk.Builder, err error) {
	builder, err = gtk.BuilderNewFromString(string(data.MainWindowUI))
	if err != nil {
		return
	}

	mv.Window, err = utils.GetUIObject[gtk.ApplicationWindow](builder, "window")
	if err != nil {
		return
	}

	mv.stack, err = utils.GetUIObject[gtk.Stack](builder, "stack")
	if err != nil {
		return
	}

	mv.postListBox, err = utils.GetUIObject[gtk.Box](builder, "postListBox")
	if err != nil {
		return
	}

	mv.postListScroll, err = utils.GetUIObject[gtk.ScrolledWindow](builder, "postListScroll")
	if err != nil {
		return
	}

	mv.postBox, err = utils.GetUIObject[gtk.Box](builder, "postBox")
	if err != nil {
		return
	}

	mv.postScroll, err = utils.GetUIObject[gtk.ScrolledWindow](builder, "postScroll")
	if err != nil {
		return
	}

	mv.closeComments, err = utils.GetUIObject[gtk.Button](builder, "closeComments")
	if err != nil {
		return
	}

	return
}

func (mv *MainView) OpenComments(postID int64) {
	mv.PostView = &PostView{}
	err := mv.PostView.SetupPostView(mv.Model.KnownPosts[postID], mv.Model.KnownPosts[postID].Comments, mv.postBox)
	if err != nil {
		log.Println(err)
	}
	mv.stack.SetTransitionType(gtk.STACK_TRANSITION_TYPE_SLIDE_LEFT)
	mv.stack.SetVisibleChild(&mv.postScroll.Container)

	mv.closeComments.Show()
}

func (mv *MainView) CloseComments() {
	mv.PostView.Destroy()
	mv.PostView = nil
	mv.stack.SetTransitionType(gtk.STACK_TRANSITION_TYPE_SLIDE_RIGHT)
	mv.stack.SetVisibleChild(&mv.postListScroll.Container)

	mv.closeComments.Hide()
}

func (mv *MainView) onNewPosts() {
	lastAddedPostIDs := mv.Model.ConsumeLastAddedPosts()
	log.Printf("Adding %d posts to MainWindow...", len(lastAddedPostIDs))

	posts := make([]model.PostModel, 0, len(lastAddedPostIDs))
	for _, postID := range lastAddedPostIDs {
		posts = append(posts, mv.Model.KnownPosts[postID])
	}
	mv.PostListView.FillPostsData(posts)

	log.Println("New posts added to MainWindow.")
}
