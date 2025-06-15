package main

import (
	"errors"
	"sync"

	"github.com/vlmoon99/near-sdk-go/collections"
	contractBuilder "github.com/vlmoon99/near-sdk-go/contract"
	"github.com/vlmoon99/near-sdk-go/env"
)

var (
	contractInstance interface{}
	contractOnce     sync.Once
)

type ContractConfig struct {
	OwnerAccountID string
	Version        string
}

type User struct {
	Username  string
	Bio       string
	CreatedAt int64
	Followers uint64
	Following uint64
	PostCount uint64
}

type Post struct {
	ID        uint64
	Author    string
	Content   string
	Likes     uint64
	CreatedAt int64
}

type Comment struct {
	ID        uint64
	PostID    uint64
	Author    string
	Content   string
	CreatedAt int64
}

type UserSettings struct {
	Theme         string
	Notifications bool
	Language      string
	LastUpdated   int64
}

type StateManager struct {
	config    *collections.UnorderedMap[string, ContractConfig]
	users     *collections.UnorderedMap[string, User]
	posts     *collections.UnorderedMap[uint64, Post]
	comments  *collections.TreeMap[uint64, Comment]
	followers *collections.LookupSet[string]
	following *collections.LookupSet[string]
	userPosts *collections.Vector[uint64]
	postLikes *collections.LookupSet[string]
	settings  *collections.LookupMap[string, UserSettings]
}

func NewStateManager() *StateManager {
	sm := &StateManager{
		config:    collections.NewUnorderedMap[string, ContractConfig]("config"),
		users:     collections.NewUnorderedMap[string, User]("users"),
		posts:     collections.NewUnorderedMap[uint64, Post]("posts"),
		comments:  collections.NewTreeMap[uint64, Comment]("comments"),
		followers: collections.NewLookupSet[string]("followers"),
		following: collections.NewLookupSet[string]("following"),
		userPosts: collections.NewVector[uint64]("user_posts"),
		postLikes: collections.NewLookupSet[string]("post_likes"),
		settings:  collections.NewLookupMap[string, UserSettings]("user_settings"),
	}

	_, err := sm.config.Get("main")
	accountId, _ := env.GetPredecessorAccountID()
	if err != nil {
		config := ContractConfig{
			OwnerAccountID: accountId,
			Version:        "1.0.0",
		}
		sm.config.Insert("main", config)
	}

	return sm
}

func (sm *StateManager) GetOwnerAccountID() string {
	config, err := sm.config.Get("main")
	if err != nil {
		env.PanicStr("failed to get contract config")
	}
	return config.OwnerAccountID
}

func (sm *StateManager) UpdateOwnerAccountID(newOwner string) {
	config, err := sm.config.Get("main")
	if err != nil {
		env.PanicStr("failed to get contract config")
	}
	config.OwnerAccountID = newOwner
	if err := sm.config.Insert("main", config); err != nil {
		env.PanicStr("failed to update contract config")
	}
}

type SocialMediaContract struct {
	state *StateManager
}

func NewSocialMediaContract() *SocialMediaContract {
	return &SocialMediaContract{
		state: NewStateManager(),
	}
}

func GetContract() interface{} {
	contractOnce.Do(func() {
		if contractInstance == nil {
			contractInstance = NewSocialMediaContract()
		}
	})
	return contractInstance
}

//go:export GetOwnerAccountID
func GetOwnerAccountID() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		contract := GetContract().(*SocialMediaContract)
		contractBuilder.ReturnValue(contract.state.GetOwnerAccountID())
		return nil
	})
}

//go:export UpdateOwnerAccountID
func UpdateOwnerAccountID() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		newOwner, err := input.JSON.GetString("new_owner")
		if err != nil {
			return err
		}

		contract := GetContract().(*SocialMediaContract)
		accountId, _ := env.GetPredecessorAccountID()

		if accountId != contract.state.GetOwnerAccountID() {
			return errors.New("only owner can update owner account ID")
		}

		contract.state.UpdateOwnerAccountID(newOwner)
		contractBuilder.ReturnValue("success")
		return nil
	})
}

//go:export CreateUser
func CreateUser() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		username, err := input.JSON.GetString("username")
		if err != nil {
			return err
		}

		bio, err := input.JSON.GetString("bio")
		if err != nil {
			return err
		}

		contract := GetContract().(*SocialMediaContract)

		user := User{
			Username:  username,
			Bio:       bio,
			CreatedAt: int64(env.GetBlockTimeMs()),
			Followers: 0,
			Following: 0,
			PostCount: 0,
		}

		if err := contract.state.users.Insert(username, user); err != nil {
			return err
		}

		contractBuilder.ReturnValue(user)
		return nil
	})
}

//go:export CreatePost
func CreatePost() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		username, err := input.JSON.GetString("username")
		if err != nil {
			return err
		}

		content, err := input.JSON.GetString("content")
		if err != nil {
			return err
		}

		contract := GetContract().(*SocialMediaContract)

		user, err := contract.state.users.Get(username)
		if err != nil {
			return err
		}

		post := Post{
			ID:        user.PostCount + 1,
			Author:    username,
			Content:   content,
			Likes:     0,
			CreatedAt: int64(env.GetBlockTimeMs()),
		}

		user.PostCount++
		if err := contract.state.users.Insert(username, user); err != nil {
			return err
		}

		if err := contract.state.posts.Insert(post.ID, post); err != nil {
			return err
		}

		if err := contract.state.userPosts.Push(post.ID); err != nil {
			return err
		}

		contractBuilder.ReturnValue(post)
		return nil
	})
}

//go:export FollowUser
func FollowUser() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		follower, err := input.JSON.GetString("follower")
		if err != nil {
			return err
		}

		following, err := input.JSON.GetString("following")
		if err != nil {
			return err
		}

		contract := GetContract().(*SocialMediaContract)

		followerKey := following + ":" + follower
		if err := contract.state.followers.Insert(followerKey); err != nil {
			return err
		}

		followingKey := follower + ":" + following
		if err := contract.state.following.Insert(followingKey); err != nil {
			return err
		}

		followerUser, err := contract.state.users.Get(follower)
		if err != nil {
			return err
		}
		followerUser.Following++
		if err := contract.state.users.Insert(follower, followerUser); err != nil {
			return err
		}

		followingUser, err := contract.state.users.Get(following)
		if err != nil {
			return err
		}
		followingUser.Followers++
		if err := contract.state.users.Insert(following, followingUser); err != nil {
			return err
		}

		contractBuilder.ReturnValue("success")
		return nil
	})
}

//go:export LikePost
func LikePost() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		username, err := input.JSON.GetString("username")
		if err != nil {
			return err
		}

		postID, err := input.JSON.GetInt("post_id")
		if err != nil {
			return err
		}

		contract := GetContract().(*SocialMediaContract)

		likeKey := string(postID) + ":" + username
		exists, err := contract.state.postLikes.Contains(likeKey)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("post already liked")
		}

		if err := contract.state.postLikes.Insert(likeKey); err != nil {
			return err
		}

		post, err := contract.state.posts.Get(uint64(postID))
		if err != nil {
			return err
		}
		post.Likes++
		if err := contract.state.posts.Insert(uint64(postID), post); err != nil {
			return err
		}

		contractBuilder.ReturnValue("success")
		return nil
	})
}

//go:export UpdateUserSettings
func UpdateUserSettings() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		username, err := input.JSON.GetString("username")
		if err != nil {
			return err
		}

		theme, err := input.JSON.GetString("theme")
		if err != nil {
			return err
		}

		notificationsStr, err := input.JSON.GetString("notifications")
		if err != nil {
			return err
		}
		notifications := notificationsStr == "true"

		language, err := input.JSON.GetString("language")
		if err != nil {
			return err
		}

		contract := GetContract().(*SocialMediaContract)

		settings, err := contract.state.settings.Get(username)
		if err != nil {
			settings = UserSettings{
				Theme:         theme,
				Notifications: notifications,
				Language:      language,
				LastUpdated:   int64(env.GetBlockTimeMs()),
			}
		} else {
			settings.Theme = theme
			settings.Notifications = notifications
			settings.Language = language
			settings.LastUpdated = int64(env.GetBlockTimeMs())
		}

		if err := contract.state.settings.Insert(username, settings); err != nil {
			return err
		}

		contractBuilder.ReturnValue(settings)
		return nil
	})
}

//go:export GetUserSettings
func GetUserSettings() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		username, err := input.JSON.GetString("username")
		if err != nil {
			return err
		}

		contract := GetContract().(*SocialMediaContract)
		settings, err := contract.state.settings.Get(username)
		if err != nil {
			return err
		}

		contractBuilder.ReturnValue(settings)
		return nil
	})
}

//go:export DeleteUserSettings
func DeleteUserSettings() {
	contractBuilder.HandleClientJSONInput(func(input *contractBuilder.ContractInput) error {
		username, err := input.JSON.GetString("username")
		if err != nil {
			return err
		}

		contract := GetContract().(*SocialMediaContract)
		if err := contract.state.settings.Remove(username); err != nil {
			return err
		}

		contractBuilder.ReturnValue("success")
		return nil
	})
}
