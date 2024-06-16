package view

import (
	"log"
	"strconv"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/mjdiliscia/LemmeRead/data"
	"github.com/mjdiliscia/LemmeRead/model"
	"github.com/mjdiliscia/LemmeRead/utils"
)

const (
	applicationTitle = "Lemme Read"
	maxPostImageSize = 580
)

type MainView struct {
	Window                *adw.ApplicationWindow
	Model                 *model.AppModel
	PostListView          PostListView
	PostView              *PostView
	PostListBottomReached func()
	CommentsClosed        func()
	OrderChanged          func(int)
	FilterChanged         func(int)

	stack          *adw.NavigationView
	postListPage   *adw.NavigationPage
	postListBox    *gtk.Box
	postListScroll *gtk.ScrolledWindow
	postPage       *adw.NavigationPage
	postBox        *gtk.Box
	postScroll     *gtk.ScrolledWindow
	search         *gtk.Button
	sortButton     *gtk.MenuButton
	sortItems      map[int]*gtk.CheckButton
	filterItems    map[int]*adw.ActionRow
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
		if position == gtk.PosBottom && mv.PostListBottomReached != nil {
			mv.PostListBottomReached()
		}
	})

	mv.stack.Connect("popped", func() {
		if mv.CommentsClosed != nil {
			mv.CommentsClosed()
		}

	})

	mv.setupSort()
	mv.setupFilter()

	mv.Window.Show()

	return nil
}

func (mv *MainView) setupSort() {
	for index, sortItem := range mv.sortItems {
		sortItem.SetActive(index == int(mv.Model.Configuration.GetOrder()))

		idx, item := index, sortItem
		sortItem.Connect("toggled", func() {
			if item.Active() {
				mv.sortButton.Popdown()
				if mv.OrderChanged != nil {
					mv.OrderChanged(idx)
				}
			}
		})
	}
}

func (mv *MainView) setupFilter() {
	selectedFilter := int(mv.Model.Configuration.GetFilter())
	filterItem := mv.filterItems[selectedFilter]
	listBox := filterItem.Parent().(*gtk.ListBox)
	listBox.Connect("row-selected", func(self *gtk.ListBox, row *gtk.ListBoxRow) {
		item := row.Cast().(*adw.ActionRow)
		mv.postListPage.SetTitle(item.Title())
		if mv.FilterChanged != nil {
			var itemId int
			for id, filter := range mv.filterItems {
				if filter.Title() == item.Title() {
					itemId = id
				}
			}
			mv.FilterChanged(itemId)
		}
	})
	listBox.SelectRow(&filterItem.ListBoxRow)
}

func (mv *MainView) CleanView() {
	mv.PostListView.CleanView()
}

func (mv *MainView) buildAndSetReferences() (builder *gtk.Builder, err error) {
	builder = gtk.NewBuilderFromString(string(data.MainWindowUI), -1)

	mv.Window, err = utils.GetUIObject[*adw.ApplicationWindow](builder, "window")
	if err != nil {
		return
	}

	mv.stack, err = utils.GetUIObject[*adw.NavigationView](builder, "stack")
	if err != nil {
		return
	}

	mv.postListPage, err = utils.GetUIObject[*adw.NavigationPage](builder, "postListPage")
	if err != nil {
		return
	}

	mv.postListBox, err = utils.GetUIObject[*gtk.Box](builder, "postListBox")
	if err != nil {
		return
	}

	mv.postListScroll, err = utils.GetUIObject[*gtk.ScrolledWindow](builder, "postListScroll")
	if err != nil {
		return
	}

	mv.postPage, err = utils.GetUIObject[*adw.NavigationPage](builder, "postPage")
	if err != nil {
		return
	}

	mv.postBox, err = utils.GetUIObject[*gtk.Box](builder, "postBox")
	if err != nil {
		return
	}

	mv.postScroll, err = utils.GetUIObject[*gtk.ScrolledWindow](builder, "postScroll")
	if err != nil {
		return
	}

	if err != nil {
		return
	}

	mv.search, err = utils.GetUIObject[*gtk.Button](builder, "search")
	if err != nil {
		return
	}

	mv.sortButton, err = utils.GetUIObject[*gtk.MenuButton](builder, "sort")
	if err != nil {
		return
	}

	mv.sortItems = make(map[int]*gtk.CheckButton)
	for i := 0; i < 8; i++ {
		mv.sortItems[i], err = utils.GetUIObject[*gtk.CheckButton](builder, "sort"+strconv.Itoa(i))
		if err != nil {
			return
		}
	}

	mv.filterItems = make(map[int]*adw.ActionRow)
	for i := 0; i < 3; i++ {
		mv.filterItems[i], err = utils.GetUIObject[*adw.ActionRow](builder, "filter"+strconv.Itoa(i))
		if err != nil {
			return
		}
	}

	return
}

func (mv *MainView) OpenComments(postID int64) {
	mv.PostView = &PostView{}
	err := mv.PostView.SetupPostView(mv.Model.KnownPosts[postID], mv.Model.KnownPosts[postID].Comments, mv.postBox)
	if err != nil {
		log.Println(err)
	}
	mv.stack.Push(mv.postPage)
}

func (mv *MainView) CloseComments() {
	mv.PostView.Destroy()
	mv.PostView = nil
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
