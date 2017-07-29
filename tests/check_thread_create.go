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

func (f *Factory) CreateThread(c *client.Forum, thread *models.Thread, forum *models.Forum, author *models.User) *models.Thread {
	if thread == nil {
		thread = f.RandomThread()
	}
	if thread.Forum == "" {
		if forum == nil {
			forum = f.CreateForum(c, nil, author)
		}
		thread.Forum = forum.Slug
	}
	if thread.Author == "" {
		if author == nil {
			author = f.CreateUser(c, nil)
		}
		thread.Author = author.Nickname
	}

	expected := *thread
	expected.ID = 42
	if forum != nil {
		expected.Forum = forum.Slug
	}
	check_create := expected.Created != nil
	result, err := c.Operations.ThreadCreate(operations.NewThreadCreateParams().
		WithSlug(thread.Forum).
		WithThread(thread).
		WithContext(Expected(201, &expected, func(data interface{}) interface{} {
			thread := data.(*models.Thread)
			if thread.ID != 0 {
				thread.ID = expected.ID
			}
			if !check_create {
				thread.Created = nil
			}
			if thread.Created != nil {
				created := strfmt.DateTime(time.Time(*thread.Created).UTC())
				thread.Created = &created
			}
			return thread
		})))
	CheckNil(err)
	if thread.Slug != result.Payload.Slug {
		log.Errorf("Unexpected created thread slug: %s -> %s", thread.Slug, result.Payload.Slug)
	}

	return result.Payload
}

func CheckThread(c *client.Forum, thread *models.Thread) {
	_, err := c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(fmt.Sprintf("%d", thread.ID)).
		WithContext(Expected(200, thread, filterThread)))
	CheckNil(err)
}

func CheckThreadCreateSimple(c *client.Forum, f *Factory, m *Modify) {
	thread := f.RandomThread()
	if thread.Slug == "" || thread.Created == nil {
		panic("Incorrect thread data")
	}

	// Slug
	if m.Bool() {
		thread.Slug = ""
	}
	// Created
	if m.Bool() {
		thread.Created = nil
	}

	// Check
	f.CreateThread(c, thread, nil, nil)
}

func CheckThreadCreateNoCase(c *client.Forum, f *Factory, m *Modify) {
	forum := f.CreateForum(c, nil, nil)
	thread := f.RandomThread()

	// Slug
	thread.Forum = m.Case(forum.Slug)

	// Check
	f.CreateThread(c, thread, forum, nil)
}

func CheckThreadCreateUnicode(c *client.Forum, f *Factory) {
	thread := f.RandomThread()
	thread.Title = "松尾芭蕉"
	thread.Message = "かれ朶に烏の\nとまりけり\n秋の暮"
	f.CreateThread(c, thread, nil, nil)
}

func CheckThreadCreateNoAuthor(c *client.Forum, f *Factory) {
	thread := f.RandomThread()
	forum := f.CreateForum(c, nil, nil)
	thread.Author = f.RandomNickname()
	_, err := c.Operations.ThreadCreate(operations.NewThreadCreateParams().
		WithSlug(forum.Slug).
		WithThread(thread).
		WithContext(ExpectedError(404, "Can't find thread author by nickname: %s", thread.Author)))
	CheckIsType(operations.NewThreadCreateNotFound(), err)
}

func CheckThreadCreateNoForum(c *client.Forum, f *Factory) {
	thread := f.RandomThread()
	user := f.CreateUser(c, nil)
	thread.Author = user.Nickname
	forum := f.RandomForum()
	_, err := c.Operations.ThreadCreate(operations.NewThreadCreateParams().
		WithSlug(forum.Slug).
		WithThread(thread).
		WithContext(ExpectedError(404, "Can't find thread forum by slug: %s", forum.Slug)))
	CheckIsType(operations.NewThreadCreateNotFound(), err)
}

func CheckThreadCreateConflict(c *client.Forum, f *Factory, m *Modify) {
	forum := f.CreateForum(c, nil, nil)
	thread := f.CreateThread(c, nil, nil, nil)

	conflict := f.RandomThread()
	conflict.Author = forum.User
	conflict.Slug = thread.Slug

	// Slug
	conflict.Slug = m.Case(thread.Slug)

	// Check
	_, err := c.Operations.ThreadCreate(operations.NewThreadCreateParams().
		WithSlug(forum.Slug).
		WithThread(conflict).
		WithContext(Expected(409, thread, filterThread)))
	CheckIsType(operations.NewThreadCreateConflict(), err)
}
