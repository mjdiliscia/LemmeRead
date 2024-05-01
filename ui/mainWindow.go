package ui

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/mjdiliscia/LemmeRead/data"
	"github.com/mjdiliscia/LemmeRead/utils"
	"go.elara.ws/go-lemmy"
)

const (
	applicationTitle  = "Lemme Read"
	maxPostImageSize  = 580
	communityIconSize = 30
)

type MainWindow struct {
	Window *gtk.ApplicationWindow
	PostList PostsUI
	Post *PostUI
	OnPostListBottomReached func()

	stack *gtk.Stack
	postListBox *gtk.Box
	postListScroll *gtk.ScrolledWindow
	postBox *gtk.Box
	postScroll *gtk.ScrolledWindow
}

func (win *MainWindow) SetupMainWindow() (err error) {
	builder, err := win.buildAndSetReferences()
	if err != nil {
		return
	}

	err = win.PostList.SetupPostsUI(win.postListBox)
	if err != nil {
		return
	}

	win.postListScroll.Connect("edge-reached", func(scroll *gtk.ScrolledWindow, position gtk.PositionType) {
		if position == gtk.POS_BOTTOM && win.OnPostListBottomReached != nil { win.OnPostListBottomReached() }
	})

	closeCommentsButton, err := utils.GetUIObject[gtk.Button](builder, "closeComments")
	if err != nil {
		return
	}
	closeCommentsButton.Connect("clicked", func() { win.CloseComments() })

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

func (win *MainWindow) OpenComments(post lemmy.PostView, comments []lemmy.CommentView) {
	win.Post = &PostUI{}
	err := win.Post.SetupPostUI(post, comments, win.postBox)
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
