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
		Name:        "post_create_simple",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateSimple),
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "post_create_unicode",
		Description: "",
		FnCheck:     CheckPostCreateUnicode,
		Deps: []string{
			"post_create_simple",
		},
	})
	Register(Checker{
		Name:        "post_create_no_thread",
		Description: "",
		FnCheck:     CheckPostCreateNoThread,
		Deps: []string{
			"post_create_simple",
		},
	})
	Register(Checker{
		Name:        "post_create_no_author",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateNoAuthor),
		Deps: []string{
			"post_create_simple",
		},
	})
	Register(Checker{
		Name:        "post_create_with_parent",
		Description: "",
		FnCheck:     CheckPostCreateWithParent,
		Deps: []string{
			"post_create_simple",
		},
	})
	Register(Checker{
		Name:        "post_create_invalid_parent",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateInvalidParent),
		Deps: []string{
			"post_create_simple",
		},
	})
	Register(Checker{
		Name:        "post_create_deep_parent",
		Description: "",
		FnCheck:     CheckPostCreateDeepParent,
		Deps: []string{
			"post_create_simple",
		},
	})
}

func CreatePost(c *client.Forum, post *models.Post, thread *models.Thread) *models.Post {
	if post == nil {
		post = RandomPost()
	}
	slug := ""
	if post.Thread == 0 {
		if thread == nil {
			thread = CreateThread(c, nil, nil, nil)
		} else {
			slug = thread.Slug
		}
		post.Thread = thread.ID
		post.Forum = thread.Forum
	} else {
		slug = fmt.Sprintf("%d", post.Thread)
	}
	if slug == "" {
		slug = fmt.Sprintf("%d", thread.ID)
	}
	if post.Author == "" {
		if thread == nil {
			post.Author = CreateUser(c, nil).Nickname
		} else {
			post.Author = thread.Author
		}
	}

	expected := *post
	expected.ID = 42
	expected.Thread = post.Thread
	expected.Created = strfmt.DateTime(time.Now())

	var parent *int64
	if post.Parent != 0 {
		parent = &post.Parent
	}
	post.Thread = 0
	check_forum := post.Forum != ""

	result, err := c.Operations.PostCreate(operations.NewPostCreateParams().
		WithSlugOrID(slug).
		WithParent(parent).
		WithPost(post).
		WithContext(Expected(201, &expected, func(data interface{}) interface{} {
			post := data.(*models.Post)
			post.ID = 0
			if !check_forum {
				post.Forum = ""
			}
			post.Created = strfmt.NewDateTime()
			return data
		})))
	CheckNil(err)

	return result.Payload
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

func CheckPostCreateNoThread(c *client.Forum) {
	post := RandomPost()
	post.Author = CreateUser(c, nil).Nickname

	var err error
	_, err = c.Operations.PostCreate(operations.NewPostCreateParams().
		WithSlugOrID(THREAD_FAKE_ID).
		WithPost(post).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewPostCreateNotFound(), err)

	_, err = c.Operations.PostCreate(operations.NewPostCreateParams().
		WithSlugOrID(RandomThread().Slug).
		WithPost(post).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewPostCreateNotFound(), err)
}

func CheckPostCreateNoAuthor(c *client.Forum, m *Modify) {
	post := RandomPost()
	post.Author = RandomNickname()
	thread := CreateThread(c, nil, nil, nil)

	_, err := c.Operations.PostCreate(operations.NewPostCreateParams().
		WithSlugOrID(m.SlugOrId(thread)).
		WithPost(post).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewPostCreateNotFound(), err)
}

func CheckPostCreateInvalidParent(c *client.Forum, m *Modify) {
	post := RandomPost()
	thread := CreateThread(c, nil, nil, nil)
	post.Author = thread.Author
	parentId := POST_FAKE_ID
	if m.Bool() {
		parentId = CreatePost(c, nil, nil).ID
	}
	_, err := c.Operations.PostCreate(operations.NewPostCreateParams().
		WithSlugOrID(m.SlugOrId(thread)).
		WithParent(&parentId).
		WithPost(post).
		WithContext(Expected(409, nil, nil)))
	CheckIsType(operations.NewPostCreateConflict(), err)
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
