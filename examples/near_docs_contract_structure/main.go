package main

import (
	"errors"
	"strconv"

	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
)

// ============================================================================
// Data Models (State Storage)
// ============================================================================

type ContractConfig struct {
	OwnerAccountID string `json:"owner_account_id"`
	Version        string `json:"version"`
}

type User struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	CreatedAt int64  `json:"created_at"`
	Followers uint64 `json:"followers"`
	Following uint64 `json:"following"`
	PostCount uint64 `json:"post_count"`
}

type Post struct {
	ID        uint64 `json:"id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	Likes     uint64 `json:"likes"`
	CreatedAt int64  `json:"created_at"`
}

type Comment struct {
	ID        uint64 `json:"id"`
	PostID    uint64 `json:"post_id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

type UserSettings struct {
	Theme         string `json:"theme"`
	Notifications bool   `json:"notifications"`
	Language      string `json:"language"`
	LastUpdated   int64  `json:"last_updated"`
}

// ============================================================================
// Input Models (DTOs)
// ============================================================================

type UpdateOwnerInput struct {
	NewOwner string `json:"new_owner"`
}

type CreateUserInput struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

type GetUserInput struct {
	Username string `json:"username"`
}

type CreatePostInput struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

type LikePostInput struct {
	Username string `json:"username"`
	PostID   uint64 `json:"post_id"`
}

type FollowUserInput struct {
	Follower  string `json:"follower"`
	Following string `json:"following"`
}

type UpdateSettingsInput struct {
	Username      string `json:"username"`
	Theme         string `json:"theme"`
	Notifications bool   `json:"notifications"`
	Language      string `json:"language"`
}

type UserTargetInput struct {
	Username string `json:"username"`
}

// ============================================================================
// Contract State
// ============================================================================

// @contract:state
type SocialMediaContract struct {
	Config    *collections.UnorderedMap[string, ContractConfig]
	Users     *collections.UnorderedMap[string, User]
	Posts     *collections.UnorderedMap[uint64, Post]
	Comments  *collections.TreeMap[uint64, Comment]
	Followers *collections.LookupSet[string] // Key: "following:follower"
	Following *collections.LookupSet[string] // Key: "follower:following"
	UserPosts *collections.Vector[uint64]    // Global list of Post IDs
	PostLikes *collections.LookupSet[string] // Key: "postID:username"
	Settings  *collections.LookupMap[string, UserSettings]
}

// ============================================================================
// Initialization
// ============================================================================

// @contract:init
func (c *SocialMediaContract) Init() string {
	// Initialize all collections
	c.Config = collections.NewUnorderedMap[string, ContractConfig]("config")
	c.Users = collections.NewUnorderedMap[string, User]("users")
	c.Posts = collections.NewUnorderedMap[uint64, Post]("posts")
	c.Comments = collections.NewTreeMap[uint64, Comment]("comments")
	c.Followers = collections.NewLookupSet[string]("followers")
	c.Following = collections.NewLookupSet[string]("following")
	c.UserPosts = collections.NewVector[uint64]("user_posts")
	c.PostLikes = collections.NewLookupSet[string]("post_likes")
	c.Settings = collections.NewLookupMap[string, UserSettings]("settings")

	// Set initial config
	ownerID, _ := env.GetPredecessorAccountID()
	config := ContractConfig{
		OwnerAccountID: ownerID,
		Version:        "1.0.0",
	}
	c.Config.Insert("main", config)

	env.LogString("Social Media Contract Initialized")
	return "Smart Contract was inited"
}

// ============================================================================
// Config Methods
// ============================================================================

// @contract:view
func (c *SocialMediaContract) GetOwnerAccountID() string {
	config, err := c.Config.Get("main")
	if err != nil {
		env.PanicStr("Contract not initialized")
	}
	return config.OwnerAccountID
}

// @contract:mutating
func (c *SocialMediaContract) UpdateOwnerAccountID(input UpdateOwnerInput) (string, error) {
	config, err := c.Config.Get("main")
	if err != nil {
		return "", errors.New("failed to get contract config")
	}

	caller, _ := env.GetPredecessorAccountID()
	if caller != config.OwnerAccountID {
		return "", errors.New("only owner can update owner account ID")
	}

	config.OwnerAccountID = input.NewOwner
	if err := c.Config.Insert("main", config); err != nil {
		return "", err
	}

	return "success", nil
}

// ============================================================================
// User Methods
// ============================================================================

// @contract:mutating
func (c *SocialMediaContract) CreateUser(input CreateUserInput) (User, error) {
	// Check if user already exists
	if exists, _ := c.Users.Contains(input.Username); exists {
		return User{}, errors.New("username already taken")
	}

	user := User{
		Username:  input.Username,
		Bio:       input.Bio,
		CreatedAt: int64(env.GetBlockTimeMs() / 1_000_000), // ms
		Followers: 0,
		Following: 0,
		PostCount: 0,
	}

	if err := c.Users.Insert(input.Username, user); err != nil {
		return User{}, err
	}

	return user, nil
}

// @contract:view
func (c *SocialMediaContract) GetUser(input GetUserInput) (User, error) {
	return c.Users.Get(input.Username)
}

// ============================================================================
// Post Methods
// ============================================================================

// @contract:mutating
func (c *SocialMediaContract) CreatePost(input CreatePostInput) (Post, error) {
	user, err := c.Users.Get(input.Username)
	if err != nil {
		return Post{}, errors.New("user not found")
	}

	// Generate ID
	newPostId := c.UserPosts.Length() + 1

	post := Post{
		ID:        newPostId,
		Author:    input.Username,
		Content:   input.Content,
		Likes:     0,
		CreatedAt: int64(env.GetBlockTimeMs() / 1_000_000),
	}

	// Update User State
	user.PostCount++
	if err := c.Users.Insert(input.Username, user); err != nil {
		return Post{}, err
	}

	// Save Post
	if err := c.Posts.Insert(post.ID, post); err != nil {
		return Post{}, err
	}

	// Add to Global Feed
	if err := c.UserPosts.Push(post.ID); err != nil {
		return Post{}, err
	}

	return post, nil
}

// @contract:mutating
func (c *SocialMediaContract) LikePost(input LikePostInput) (string, error) {
	// check user exists
	if exists, _ := c.Users.Contains(input.Username); !exists {
		return "", errors.New("user not found")
	}

	// Create composite key: "postID:username"
	likeKey := strconv.FormatUint(input.PostID, 10) + ":" + input.Username

	exists, err := c.PostLikes.Contains(likeKey)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("post already liked")
	}

	// Record Like
	if err := c.PostLikes.Insert(likeKey); err != nil {
		return "", err
	}

	// Update Post Count
	post, err := c.Posts.Get(input.PostID)
	if err != nil {
		return "", errors.New("post not found")
	}
	post.Likes++
	if err := c.Posts.Insert(input.PostID, post); err != nil {
		return "", err
	}

	return "success", nil
}

// ============================================================================
// Social Graph Methods
// ============================================================================

// @contract:mutating
func (c *SocialMediaContract) FollowUser(input FollowUserInput) (string, error) {
	if input.Follower == input.Following {
		return "", errors.New("cannot follow yourself")
	}

	// Check validity
	followerUser, err := c.Users.Get(input.Follower)
	if err != nil {
		return "", errors.New("follower user does not exist")
	}
	followingUser, err := c.Users.Get(input.Following)
	if err != nil {
		return "", errors.New("following user does not exist")
	}

	// Insert into Sets
	followerKey := input.Following + ":" + input.Follower
	if err := c.Followers.Insert(followerKey); err != nil {
		return "", err
	}

	followingKey := input.Follower + ":" + input.Following
	if err := c.Following.Insert(followingKey); err != nil {
		return "", err
	}

	// Update Counts
	followerUser.Following++
	if err := c.Users.Insert(input.Follower, followerUser); err != nil {
		return "", err
	}

	followingUser.Followers++
	if err := c.Users.Insert(input.Following, followingUser); err != nil {
		return "", err
	}

	return "success", nil
}

// ============================================================================
// User Settings Methods
// ============================================================================

// @contract:mutating
func (c *SocialMediaContract) UpdateUserSettings(input UpdateSettingsInput) (UserSettings, error) {
	currentSettings, err := c.Settings.Get(input.Username)
	if err != nil {
		// Create new if not exists
		currentSettings = UserSettings{
			Theme:         input.Theme,
			Notifications: input.Notifications,
			Language:      input.Language,
			LastUpdated:   int64(env.GetBlockTimeMs() / 1_000_000),
		}
	} else {
		// Update existing
		currentSettings.Theme = input.Theme
		currentSettings.Notifications = input.Notifications
		currentSettings.Language = input.Language
		currentSettings.LastUpdated = int64(env.GetBlockTimeMs() / 1_000_000)
	}

	if err := c.Settings.Insert(input.Username, currentSettings); err != nil {
		return UserSettings{}, err
	}

	return currentSettings, nil
}

// @contract:view
func (c *SocialMediaContract) GetUserSettings(input UserTargetInput) (UserSettings, error) {
	return c.Settings.Get(input.Username)
}

// @contract:mutating
func (c *SocialMediaContract) DeleteUserSettings(input UserTargetInput) (string, error) {
	if err := c.Settings.Remove(input.Username); err != nil {
		return "", err
	}
	return "success", nil
}
