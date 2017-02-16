package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/go-openapi/strfmt"
	"time"
)

func init() {
	Register(Checker{
		Name:        "posts_create_simple",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateSimple),
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_empty",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateEmpty),
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_unicode",
		Description: "",
		FnCheck:     CheckPostCreateUnicode,
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_no_thread",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateNoThread),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_no_author",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateNoAuthor),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_with_parent",
		Description: "",
		FnCheck:     CheckPostCreateWithParent,
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_invalid_parent",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateInvalidParent),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_deep_parent",
		Description: "",
		FnCheck:     CheckPostCreateDeepParent,
		Deps: []string{
			"posts_create_simple",
		},
	})
}

func CreatePosts(c *client.Forum, posts []*models.Post, thread *models.Thread) []*models.Post {
	if len(posts) == 0 {
		return []*models.Post{}
	}
	var postsThread int32 = 0
	for _, post := range posts {
		if post.Thread != 0 {
			if postsThread == 0 {
				postsThread = post.Thread
			}
			if postsThread != post.Thread {
				panic("Invalid test data: can't create multiple posts in differ threads")
			}
		}
	}
	slug := ""
	var postsForum string
	if postsThread == 0 {
		if thread == nil {
			thread = CreateThread(c, nil, nil, nil)
		} else {
			slug = thread.Slug
		}
		postsThread = thread.ID
		postsForum = thread.Forum
	}
	if slug == "" {
		slug = fmt.Sprintf("%d", postsThread)
	}
	var author string

	check_forum := postsForum != ""
	if !check_forum {
		postsForum = RandomForum().Slug
	}

	var expected []*models.Post
	if thread != nil {
		author = thread.Author
	}
	for n, post := range posts {
		if post.Author == "" {
			if author == "" {
				author = CreateUser(c, nil).Nickname
			}
			post.Author = author
		}
		post.Thread = 0

		expectedPost := *post
		expectedPost.ID = int64(42 + n)
		expectedPost.Thread = postsThread
		expectedPost.Created = strfmt.DateTime(time.Now())
		expectedPost.Forum = postsForum
		expected = append(expected, &expectedPost)
	}

	result, err := c.Operations.PostsCreate(operations.NewPostsCreateParams().
		WithSlugOrID(slug).
		WithPosts(posts).
		WithContext(Expected(201, &expected, func(data interface{}) interface{} {
			posts := data.(*[]*models.Post)
			for _, post := range *posts {
				post.ID = 0
				if !check_forum {
					post.Forum = ""
				}
				post.Created = strfmt.NewDateTime()
			}
			return data
		})))
	CheckNil(err)

	return result.Payload
}

func CreatePost(c *client.Forum, post *models.Post, thread *models.Thread) *models.Post {
	if post == nil {
		post = RandomPost()
	}
	return CreatePosts(c, []*models.Post{post}, thread)[0]
}

func CheckPost(c *client.Forum, post *models.Post) {
	_, err := c.Operations.PostGetOne(operations.NewPostGetOneParams().
		WithID(post.ID).
		WithContext(Expected(200, &models.PostFull{
			Post: post,
		}, nil)))
	CheckNil(err)
}

func CheckPostCreateSimple(c *client.Forum, m *Modify) {
	thread := CreateThread(c, nil, nil, nil)
	if m.Bool() {
		thread.Slug = ""
	}
	CreatePost(c, nil, thread)
}

func CheckPostCreateEmpty(c *client.Forum, m *Modify) {
	thread := CreateThread(c, nil, nil, nil)
	if m.Bool() {
		thread.Slug = ""
	}
	CreatePosts(c, []*models.Post{}, thread)
}

func CheckPostCreateNoThread(c *client.Forum, m *Modify) {
	posts := []*models.Post{}
	if m.Bool() {
		post := RandomPost()
		post.Author = CreateUser(c, nil).Nickname
		posts = append(posts, post)
	}
	var err error
	_, err = c.Operations.PostsCreate(operations.NewPostsCreateParams().
		WithSlugOrID(THREAD_FAKE_ID).
		WithPosts(posts).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewPostsCreateNotFound(), err)

	_, err = c.Operations.PostsCreate(operations.NewPostsCreateParams().
		WithSlugOrID(RandomThread().Slug).
		WithPosts(posts).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewPostsCreateNotFound(), err)
}

func CheckPostCreateNoAuthor(c *client.Forum, m *Modify) {
	post := RandomPost()
	post.Author = RandomNickname()
	thread := CreateThread(c, nil, nil, nil)

	_, err := c.Operations.PostsCreate(operations.NewPostsCreateParams().
		WithSlugOrID(m.SlugOrId(thread)).
		WithPosts([]*models.Post{post}).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewPostsCreateNotFound(), err)
}

func CheckPostCreateInvalidParent(c *client.Forum, m *Modify) {
	post := RandomPost()
	thread := CreateThread(c, nil, nil, nil)
	post.Author = thread.Author
	parentId := POST_FAKE_ID
	if m.Bool() {
		parentId = CreatePost(c, nil, nil).ID
	}
	post.Parent = parentId
	_, err := c.Operations.PostsCreate(operations.NewPostsCreateParams().
		WithSlugOrID(m.SlugOrId(thread)).
		WithPosts([]*models.Post{post}).
		WithContext(Expected(409, nil, nil)))
	CheckIsType(operations.NewPostsCreateConflict(), err)
}

func CheckPostCreateWithParent(c *client.Forum) {
	post := RandomPost()
	thread := CreateThread(c, nil, nil, nil)
	post.Author = thread.Author
	post.Parent = CreatePost(c, nil, thread).ID
	CreatePost(c, post, thread)
}

func CheckPostCreateDeepParent(c *client.Forum) {
	thread := CreateThread(c, nil, nil, nil)
	var parent int64
	for level := 0; level < 0x10; level++ {
		post := RandomPost()
		post.Author = thread.Author
		post.Parent = parent
		CreatePost(c, post, thread)
		parent = post.ID
	}
}

func CheckPostCreateUnicode(c *client.Forum) {
	post := RandomPost()
	post.Message = "大象销售。不贵。"
	CreatePost(c, post, nil)
}
