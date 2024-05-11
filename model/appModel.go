package model

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/gotk3/gotk3/glib"
	"go.elara.ws/go-lemmy"
)

type AppModel struct {
	KnownPosts map[int64]PostModel
	NewPosts func()

	lastAddedPosts []int64
	nextPageToRetrieve int64
	pendingProcesses []string
	lemmyClient *lemmy.Client
	lemmyContext context.Context
}

func (am *AppModel) Init() {
	am.nextPageToRetrieve = 0
	am.KnownPosts = make(map[int64]PostModel)
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
			Type: lemmy.NewOptional(lemmy.ListingTypeSubscribed),
			Page: lemmy.NewOptional(page+1),
		})
		log.Printf("Posts from page %d retrieval completed. Error: %v", page, err)
		callInMain(func() error { return am.addPosts(response.Posts, err) }, func(err error) {
			processIndex := slices.Index(am.pendingProcesses, processID)
			am.pendingProcesses = append(am.pendingProcesses[:processIndex], am.pendingProcesses[processIndex+1:]...)

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
		response, err := am.lemmyClient.Comments(am.lemmyContext, lemmy.GetComments{
			PostID: lemmy.NewOptional(am.KnownPosts[postID].Post.ID),
			Limit: lemmy.NewOptional(am.KnownPosts[postID].Counts.Comments),
		})
		if err != nil {
			callInMain(func() error {return err}, callback)
			return
		}
		callInMain(func() error {
			if post, ok := am.KnownPosts[postID]; ok {
				err := post.AddComments(response.Comments, err)
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
	defer func(){ am.lastAddedPosts = make([]int64, 0) }()
	return am.lastAddedPosts
}

func (am *AppModel) addPosts(posts []lemmy.PostView, err error) error {
	if err != nil {
		log.Println("addPost called with errors, ignoring call.")
		return err
	}

	log.Printf("Adding %d new posts to local DB.", len(posts))
	for _, post := range(posts) {
		if _, ok := am.KnownPosts[post.Post.ID]; !ok {
			postModel := PostModel{PostView: post}
			postID := post.Post.ID

			processID := fmt.Sprintf("post%d", postID)
			am.pendingProcesses = append(am.pendingProcesses, processID)
			postModel.Init(func (err error) {
				processIndex := slices.Index(am.pendingProcesses, processID)
				am.pendingProcesses = append(am.pendingProcesses[:processIndex], am.pendingProcesses[processIndex+1:]...)

				if err != nil {
					log.Printf("Something went wrong with post %d, skipping: %s", postID, err)
					return
				}
				am.KnownPosts[postID] = postModel
				am.lastAddedPosts = append(am.lastAddedPosts, postID)
				log.Printf("Added new post %d to %p DB with %d posts.", postID, &am.KnownPosts, len(am.KnownPosts))
				if am.NewPosts != nil {
					am.NewPosts()
				}
			})
		}
	}
	return err
}

func callInMain(function func() error, callback func(error)) {
	glib.IdleAdd(func() bool {
		err := function()
		callback(err)
		return false
	})
}
