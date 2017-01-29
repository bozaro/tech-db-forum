package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/client"
	"github.com/bozaro/tech-db-forum/client/operations"
	"github.com/bozaro/tech-db-forum/models"
	"github.com/go-openapi/strfmt"
	"time"
)

func init() {
	Register(Checker{
		Name:        "post_create_simple",
		Description: "",
		FnCheck:     CheckPostCreateSimple,
		Deps: []string{
			"thread_create_simple",
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

func CheckPost(c *client.Forum, thread *models.Thread) {
	_, err := c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(fmt.Sprintf("%d", thread.ID)).
		WithContext(Expected(200, thread, nil)))
	CheckNil(err)
}

func CheckPostCreateSimple(c *client.Forum) {
	for pass := 1; pass <= 2; pass++ {
		Checkpoint(c, fmt.Sprintf("Pass %d", pass))
		thread := CreateThread(c, nil, nil, nil)
		if pass == 2 {
			thread.Slug = ""
		}
		CreatePost(c, nil, thread)
	}
}
