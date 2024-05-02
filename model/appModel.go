package model

import (
	"context"
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"go.elara.ws/go-lemmy"
)

type AppModel struct {
	KnownPosts map[int64]PostModel

	lemmyClient *lemmy.Client
	lemmyContext context.Context
}

func (am *AppModel) Init() {
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

func (am *AppModel) RetrievePosts(page int64, callback func(error)) {
	go func() {
		response, err := am.lemmyClient.Posts(am.lemmyContext, lemmy.GetPosts{
			Type: lemmy.NewOptional(lemmy.ListingTypeSubscribed),
			Page: lemmy.NewOptional(page+1),
		})
		callInMain(func() error { return am.addPosts(response.Posts, err) }, callback)
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

func (am *AppModel) RetrieveComments(post *PostModel, callback func(error)) {
	go func() {
		response, err := am.lemmyClient.Comments(am.lemmyContext, lemmy.GetComments{
			PostID: lemmy.NewOptional(post.Post.ID),
			Limit: lemmy.NewOptional(post.Counts.Comments),
		})
		callInMain(func() error { return post.AddComments(response.Comments, err) }, callback)
	}()
}

func (am *AppModel) addPosts(posts []lemmy.PostView, err error) error {
	if err != nil {
		return err
	}

	for _, post := range(posts) {
		if _, ok := am.KnownPosts[post.Post.ID]; !ok {
			am.KnownPosts[post.Post.ID] = PostModel{PostView: post}
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
