package ui

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

type MainWindow struct {
	Window *gtk.ApplicationWindow
	Model *model.AppModel
	PostList PostsUI
	Post *PostUI
	PostListBottomReached func()
	CloseCommentsClicked func()

	stack *gtk.Stack
	postListBox *gtk.Box
	postListScroll *gtk.ScrolledWindow
	postBox *gtk.Box
	postScroll *gtk.ScrolledWindow
}

func (win *MainWindow) SetupMainWindow(appModel *model.AppModel) (err error) {
	win.Model = appModel

	builder, err := win.buildAndSetReferences()
	if err != nil {
		return
	}

	err = win.PostList.SetupPostsUI(win.postListBox)
	if err != nil {
		return
	}

	win.postListScroll.Connect("edge-reached", func(scroll *gtk.ScrolledWindow, position gtk.PositionType) {
		if position == gtk.POS_BOTTOM && win.PostListBottomReached != nil { win.PostListBottomReached() }
	})

	closeCommentsButton, err := utils.GetUIObject[gtk.Button](builder, "closeComments")
	if err != nil {
		return
	}
	closeCommentsButton.Connect("clicked", func() {
		if win.CloseCommentsClicked != nil { win.CloseCommentsClicked() }
	})

	win.Window.Show()

	return nil
}

func (win *MainWindow) buildAndSetReferences() (builder *gtk.Builder, err error) {
	builder, err = gtk.BuilderNewFromString(string(data.MainWindowUI))
	if err != nil {
		return
	}

	win.Window, err = utils.GetUIObject[gtk.ApplicationWindow](builder, "window")
	if err != nil {
		return
	}

	win.stack, err = utils.GetUIObject[gtk.Stack](builder, "stack")
	if err != nil {
		return
	}

	win.postListBox, err = utils.GetUIObject[gtk.Box](builder, "postListBox")
	if err != nil {
		return
	}

	win.postListScroll, err = utils.GetUIObject[gtk.ScrolledWindow](builder, "postListScroll")
	if err != nil {
		return
	}

	win.postBox, err = utils.GetUIObject[gtk.Box](builder, "postBox")
	if err != nil {
		return
	}

	win.postScroll, err = utils.GetUIObject[gtk.ScrolledWindow](builder, "postScroll")
	if err != nil {
		return
	}

	return
}

func (win *MainWindow) OpenComments(postID int64) {
	win.Post = &PostUI{}
	err := win.Post.SetupPostUI(win.Model.KnownPosts[postID], win.Model.KnownPosts[postID].Comments, win.postBox)
	if err != nil {
		log.Println(err)
	}
	win.stack.SetTransitionType(gtk.STACK_TRANSITION_TYPE_SLIDE_LEFT)
	win.stack.SetVisibleChild(&win.postScroll.Container)
}

func (win *MainWindow) CloseComments() {
	win.Post.Destroy()
	win.Post = nil
	win.stack.SetTransitionType(gtk.STACK_TRANSITION_TYPE_SLIDE_RIGHT)
	win.stack.SetVisibleChild(&win.postListScroll.Container)
}

func (win *MainWindow) OnNewPosts() {
	lastAddedPostIDs := win.Model.ConsumeLastAddedPosts()
	log.Printf("Adding %d posts to MainWindow...", len(lastAddedPostIDs))

	posts := make([]model.PostModel, 0, len(lastAddedPostIDs))
	for _, postID := range(lastAddedPostIDs) {
		posts = append(posts, win.Model.KnownPosts[postID])
	}
	win.PostList.FillPostsData(posts)

	log.Println("New posts added to MainWindow.")
}
