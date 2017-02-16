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
		FnCheck:     Modifications(CheckThreadCreateSimple),
		Deps: []string{
			"forum_get_one_simple",
		},
	})
	Register(Checker{
		Name:        "thread_create_nocase",
		Description: "",
		FnCheck:     Modifications(CheckThreadCreateNoCase),
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "thread_create_noforum",
		Description: "",
		FnCheck:     CheckThreadCreateNoForum,
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "thread_create_noauthor",
		Description: "",
		FnCheck:     CheckThreadCreateNoAuthor,
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "thread_create_unicode",
		Description: "",
		FnCheck:     CheckThreadCreateUnicode,
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "thread_create_conflict",
		Description: "",
		FnCheck:     Modifications(CheckThreadCreateConflict),
		Deps: []string{
			"thread_create_simple",
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
	if forum != nil {
		expected.Forum = forum.Slug
	}
	check_create := !time.Time(expected.Created).IsZero()
	result, err := c.Operations.ThreadCreate(operations.NewThreadCreateParams().
		WithSlug(thread.Forum).
		WithThread(thread).
		WithContext(Expected(201, &expected, func(data interface{}) interface{} {
			thread := data.(*models.Thread)
			thread.ID = 0
			if !check_create {
				thread.Created = strfmt.NewDateTime()
			} else {
				thread.Created = strfmt.DateTime(time.Time(thread.Created).UTC())
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

func CheckThreadCreateSimple(c *client.Forum, m *Modify) {
	thread := RandomThread()
	if thread.Slug == "" || time.Time(thread.Created).IsZero() {
		panic("Incorrect test login")
	}

	// Slug
	if m.Bool() {
		thread.Slug = ""
	}
	// Created
	if m.Bool() {
		thread.Created = strfmt.NewDateTime()
	}

	// Check
	CreateThread(c, thread, nil, nil)
}

func CheckThreadCreateNoCase(c *client.Forum, m *Modify) {
	forum := CreateForum(c, nil, nil)
	thread := RandomThread()

	// Slug
	thread.Forum = m.Case(forum.Slug)

	// Check
	CreateThread(c, thread, forum, nil)
}

func CheckThreadCreateUnicode(c *client.Forum) {
	thread := RandomThread()
	thread.Title = "松尾芭蕉"
	thread.Message = "かれ朶に烏の\nとまりけり\n秋の暮"
	CreateThread(c, thread, nil, nil)
}

func CheckThreadCreateNoAuthor(c *client.Forum) {
	thread := RandomThread()
	forum := CreateForum(c, nil, nil)
	thread.Author = RandomNickname()
	_, err := c.Operations.ThreadCreate(operations.NewThreadCreateParams().
		WithSlug(forum.Slug).
		WithThread(thread).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewThreadCreateNotFound(), err)
}

func CheckThreadCreateNoForum(c *client.Forum) {
	thread := RandomThread()
	forum := RandomForum()
	_, err := c.Operations.ThreadCreate(operations.NewThreadCreateParams().
		WithSlug(forum.Slug).
		WithThread(thread).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewThreadCreateNotFound(), err)
}

func CheckThreadCreateConflict(c *client.Forum, m *Modify) {
	forum := CreateForum(c, nil, nil)
	thread := CreateThread(c, nil, nil, nil)

	conflict := RandomThread()
	conflict.Author = forum.User
	conflict.Slug = thread.Slug

	// Slug
	conflict.Slug = m.Case(thread.Slug)

	// Check
	_, err := c.Operations.ThreadCreate(operations.NewThreadCreateParams().
		WithSlug(forum.Slug).
		WithThread(conflict).
		WithContext(Expected(409, thread, nil)))
	CheckIsType(operations.NewThreadCreateConflict(), err)
}
