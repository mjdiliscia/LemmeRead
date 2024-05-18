package model

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/gotk3/gotk3/glib"
	"go.elara.ws/go-lemmy"
)

const MAX_COMMENTS_PER_REQUEST int64 = 40

type AppModel struct {
	KnownPosts    map[int64]PostModel
	NewPosts      func()
	Configuration AppModelConfiguration

	lastAddedPosts     []int64
	nextPageToRetrieve int64
	pendingProcesses   []string
	lemmyClient        *lemmy.Client
	lemmyContext       context.Context
}

func (am *AppModel) Init() {
	am.Configuration = NewAppModelConfiguration("config.json")
	am.CleanModel()
}

func (am *AppModel) InitializeLemmyClient(url string, username string, password string, callback func(error)) {
	var err error
	am.lemmyClient, err = lemmy.New(url)
	if err != nil {
		callback(fmt.Errorf("Couldn't create a Lemmy Client: %s", err))
	}

	am.lemmyContext = context.Background()

	go func() {
		err = am.lemmyClient.ClientLogin(am.lemmyContext, lemmy.Login{
			UsernameOrEmail: username,
			Password:        password,
			TOTP2FAToken:    lemmy.NewOptionalNil[string](),
		})

		callInMain(func() error { return err }, callback)
	}()
}

func (am *AppModel) CleanModel() {
	am.nextPageToRetrieve = 0
	am.KnownPosts = make(map[int64]PostModel)
	am.lastAddedPosts = make([]int64, 0)
	am.pendingProcesses = make([]string, 0)
}

func (am *AppModel) RetrieveMorePosts(callback func(error)) {
	am.RetrievePosts(am.nextPageToRetrieve, func(err error) {
		if err == nil {
			am.nextPageToRetrieve++
		}
		callback(err)
	})
}

func (am *AppModel) RetrievePosts(page int64, callback func(error)) {
	if len(am.pendingProcesses) > 0 {
		callback(fmt.Errorf("Already retrieving posts, ignoring."))
		return
	}

	log.Printf("Retrieving posts from page %d...", page)

	processID := fmt.Sprintf("list%d", page)
	am.pendingProcesses = append(am.pendingProcesses, processID)
	go func() {
		response, err := am.lemmyClient.Posts(am.lemmyContext, lemmy.GetPosts{
			Type: lemmy.NewOptional(am.getCurrentType()),
			Page: lemmy.NewOptional(page + 1),
			Sort: lemmy.NewOptional(am.getCurrentSort()),
		})
		log.Printf("Posts from page %d retrieval completed. Error: %v", page, err)
		callInMain(func() error {
			if slices.Index(am.pendingProcesses, processID) == -1 {
				return fmt.Errorf("Process %s no longer needed", processID)
			}

			return am.addPosts(response.Posts, err)
		}, func(err error) {
			processIndex := slices.Index(am.pendingProcesses, processID)
			if processIndex != -1 {
				am.pendingProcesses = append(am.pendingProcesses[:processIndex], am.pendingProcesses[processIndex+1:]...)
			}

			callback(err)
		})
	}()
}

func (am *AppModel) RetrievePost(postId int64, callback func(error)) {
	go func() {
		response, err := am.lemmyClient.Post(am.lemmyContext, lemmy.GetPost{
			ID: lemmy.NewOptional(postId),
		})
		callInMain(func() error { return am.addPosts([]lemmy.PostView{response.PostView}, err) }, callback)
	}()
}

func (am *AppModel) RetrieveComments(postID int64, callback func(error)) {
	go func() {
		remainingPages := 1 + am.KnownPosts[postID].Counts.Comments/MAX_COMMENTS_PER_REQUEST
		collectedComments := make([]lemmy.CommentView, 0, am.KnownPosts[postID].Counts.Comments)

		for ; remainingPages > 0; remainingPages-- {
			log.Printf("Asking for comments page %d", remainingPages)
			response, err := am.lemmyClient.Comments(am.lemmyContext, lemmy.GetComments{
				PostID: lemmy.NewOptional(am.KnownPosts[postID].Post.ID),
				Limit:  lemmy.NewOptional(MAX_COMMENTS_PER_REQUEST),
				Page:   lemmy.NewOptional(remainingPages),
			})
			if err != nil {
				callInMain(func() error { return err }, callback)
				return
			}
			collectedComments = append(collectedComments, response.Comments...)
		}

		callInMain(func() error {
			if post, ok := am.KnownPosts[postID]; ok {
				err := post.AddComments(collectedComments, nil)
				am.KnownPosts[postID] = post
				return err
			} else {
				keys := make([]int64, len(am.KnownPosts))
				i := 0
				for key := range am.KnownPosts {
					keys[i] = key
					i++
				}
				log.Printf("Known posts: %v", keys)
				return fmt.Errorf("Post %d couldn't be found in local DB", postID)
			}
		}, callback)
	}()
}

func (am *AppModel) ConsumeLastAddedPosts() []int64 {
	var (
		beginReady int = -1
		endReady   int = -1
	)

	for idx, postId := range am.lastAddedPosts {
		if postId == 0 && beginReady == -1 {
			return make([]int64, 0)
		}
		if postId > 0 && beginReady == -1 {
			beginReady = idx
		}
		if postId == 0 && endReady == -1 {
			endReady = idx
		}
	}

	if endReady == -1 {
		endReady = len(am.lastAddedPosts)
	}

	defer func() {
		for idx := beginReady; idx < endReady; idx++ {
			am.lastAddedPosts[idx] = -1
		}
	}()

	response := append(make([]int64, 0), am.lastAddedPosts[beginReady:endReady]...)
	return response
}

func (am *AppModel) addPosts(posts []lemmy.PostView, err error) error {
	if err != nil {
		log.Println("addPost called with errors, ignoring call.")
		return err
	}

	log.Printf("Adding %d new posts to local DB.", len(posts))
	am.lastAddedPosts = make([]int64, len(posts))
	for idx, post := range posts {
		if _, ok := am.KnownPosts[post.Post.ID]; !ok {
			postModel := PostModel{PostView: post}
			postID := post.Post.ID
			postIdx := idx

			processID := fmt.Sprintf("post%d", postID)
			am.pendingProcesses = append(am.pendingProcesses, processID)
			postModel.Init(func(err error) {
				processIndex := slices.Index(am.pendingProcesses, processID)
				if processIndex == -1 {
					log.Printf("Process for post %d not needed anymore, skipping: %s", postID, err)
					return
				}
				am.pendingProcesses = append(am.pendingProcesses[:processIndex], am.pendingProcesses[processIndex+1:]...)

				if err != nil {
					log.Printf("Something went wrong with post %d, skipping: %s", postID, err)
					return
				}
				am.KnownPosts[postID] = postModel
				am.lastAddedPosts[postIdx] = postID
				log.Printf("Added new post %d to %p DB with %d posts.", postID, &am.KnownPosts, len(am.KnownPosts))
				if am.NewPosts != nil {
					am.NewPosts()
				}
			})
		}
	}
	return err
}

func (am *AppModel) getCurrentSort() lemmy.SortType {
	order := am.Configuration.GetOrder()
	switch order {
	case PostOrderActive:
		return lemmy.SortTypeActive
	case PostOrderHot:
		return lemmy.SortTypeHot
	case PostOrderScaled:
		return lemmy.SortTypeScaled
	case PostOrderControversial:
		return lemmy.SortTypeControversial
	case PostOrderNew:
		return lemmy.SortTypeNew
	case PostOrderOld:
		return lemmy.SortTypeOld
	case PostOderMostComments:
		return lemmy.SortTypeMostComments
	case PostOrderNewComments:
		return lemmy.SortTypeNewComments
	}
	return lemmy.SortTypeActive
}

func (am *AppModel) getCurrentType() lemmy.ListingType {
	filter := am.Configuration.GetFilter()
	switch filter {
	case PostFilterSubscribed:
		return lemmy.ListingTypeSubscribed
	case PostFilterLocal:
		return lemmy.ListingTypeLocal
	case PostFilterAll:
		return lemmy.ListingTypeAll
	}
	return lemmy.ListingTypeSubscribed
}

func callInMain(function func() error, callback func(error)) {
	glib.IdleAdd(func() bool {
		err := function()
		callback(err)
		return false
	})
}
