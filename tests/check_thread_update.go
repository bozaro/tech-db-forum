package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
)

func init() {
	Register(Checker{
		Name:        "thread_update_simple",
		Description: "",
		FnCheck:     Modifications(CheckThreadUpdateSimple),
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "thread_update_empty",
		Description: "",
		FnCheck:     CheckThreadUpdateEmpty,
		Deps: []string{
			"thread_update_simple",
		},
	})
	Register(Checker{
		Name:        "thread_update_part",
		Description: "",
		FnCheck:     Modifications(CheckThreadUpdatePart),
		Deps: []string{
			"thread_update_simple",
		},
	})
	Register(Checker{
		Name:        "thread_update_notfound",
		Description: "",
		FnCheck:     CheckThreadUpdateNotFound,
		Deps: []string{
			"thread_update_simple",
		},
	})
}

func CheckThreadUpdateSimple(c *client.Forum, f *Factory, m *Modify) {
	thread := f.CreateThread(c, nil, nil, nil)

	temp := f.RandomThread()
	update := models.ThreadUpdate{}
	update.Title = temp.Title
	update.Message = temp.Message

	expected := *thread
	expected.Title = update.Title
	expected.Message = update.Message

	id := m.SlugOrId(thread)

	c.Operations.ThreadUpdate(operations.NewThreadUpdateParams().
		WithSlugOrID(id).
		WithThread(&update).
		WithContext(Expected(200, &expected, filterThread)))

	CheckThread(c, &expected)
}

func CheckThreadUpdateEmpty(c *client.Forum, f *Factory) {
	thread := f.CreateThread(c, nil, nil, nil)

	c.Operations.ThreadUpdate(operations.NewThreadUpdateParams().
		WithSlugOrID(thread.Slug).
		WithThread(&models.ThreadUpdate{}).
		WithContext(Expected(200, thread, filterThread)))

	CheckThread(c, thread)
}

func CheckThreadUpdatePart(c *client.Forum, f *Factory, m *Modify) {
	fake := f.RandomThread()
	expected := f.CreateThread(c, nil, nil, nil)
	update := &models.ThreadUpdate{}

	// Slug or ID
	id := m.SlugOrId(expected)
	// Title
	if m.Bool() {
		expected.Title = fake.Title
		update.Title = fake.Title
	}
	// Message
	if m.Bool() {
		expected.Message = fake.Message
		update.Message = fake.Message
	}

	// Check
	c.Operations.ThreadUpdate(operations.NewThreadUpdateParams().
		WithSlugOrID(id).
		WithThread(update).
		WithContext(Expected(200, expected, filterThread)))

	CheckThread(c, expected)
}

func CheckThreadUpdateNotFound(c *client.Forum, f *Factory) {
	thread := f.RandomThread()
	_, err := c.Operations.ThreadUpdate(operations.NewThreadUpdateParams().
		WithSlugOrID(thread.Slug).
		WithThread(&models.ThreadUpdate{
			Title:   thread.Title,
			Message: thread.Message,
		}).
		WithContext(ExpectedError(404, "Can't find thread by slug: %s", thread.Slug)))
	CheckIsType(operations.NewThreadUpdateNotFound(), err)

	_, err = c.Operations.ThreadUpdate(operations.NewThreadUpdateParams().
		WithSlugOrID(THREAD_FAKE_ID).
		WithThread(&models.ThreadUpdate{
			Title:   thread.Title,
			Message: thread.Message,
		}).
		WithContext(ExpectedError(404, "Can't find forum by id: %s", THREAD_FAKE_ID)))
	CheckIsType(operations.NewThreadUpdateNotFound(), err)
}
