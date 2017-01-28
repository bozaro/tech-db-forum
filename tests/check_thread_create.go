package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/client"
	"github.com/bozaro/tech-db-forum/client/operations"
	"github.com/bozaro/tech-db-forum/models"
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
	_, err := c.Operations.ThreadCreate(operations.NewThreadCreateParams().
		WithSlug(thread.Forum).
		WithThread(thread).
		WithContext(Expected(201, &expected, func(data interface{}) interface{} {
			data.(*models.Thread).ID = 0
			return data
		})))
	CheckNil(err)

	return thread
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
