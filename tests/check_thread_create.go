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
		Name:        "thread_create_simple",
		Description: "",
		FnCheck:     CheckThreadCreateSimple,
		Deps: []string{
			"forum_get_one_simple",
		},
	})
	Register(Checker{
		Name:        "thread_create_noslug",
		Description: "",
		FnCheck:     CheckThreadCreateNoSlug,
		Deps: []string{
			"forum_get_one_simple",
		},
	})
}

func CreateThread(c *client.Forum, thread *models.Thread, forum *models.Forum, author *models.User) *models.Thread {
	if thread == nil {
		thread = RandomThread()
	}
	if thread.Forum == "" {
		if forum == nil {
			forum = CreateForum(c, nil, author)
		}
		thread.Forum = forum.Slug
	}
	if thread.Author == "" {
		if author == nil {
			author = CreateUser(c, nil)
		}
		thread.Author = author.Nickname
	}

	expected := *thread
	expected.ID = 42
	check_create := !time.Time(expected.Created).IsZero()
	result, err := c.Operations.ThreadCreate(operations.NewThreadCreateParams().
		WithSlug(thread.Forum).
		WithThread(thread).
		WithContext(Expected(201, &expected, func(data interface{}) interface{} {
			thread := data.(*models.Thread)
			thread.ID = 0
			if !check_create {
				thread.Created = strfmt.NewDateTime()
			}
			return thread
		})))
	CheckNil(err)

	return result.Payload
}

func CheckThread(c *client.Forum, thread *models.Thread) {
	_, err := c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(fmt.Sprintf("%d", thread.ID)).
		WithContext(Expected(200, thread, nil)))
	CheckNil(err)
}

func CheckThreadCreateSimple(c *client.Forum) {
	CreateThread(c, nil, nil, nil)
}

func CheckThreadCreateNoSlug(c *client.Forum) {
	thread := RandomThread()
	thread.Slug = ""
	CreateThread(c, thread, nil, nil)
	CreateThread(c, thread, nil, nil)
}
