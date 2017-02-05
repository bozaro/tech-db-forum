package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
)

func init() {
	Register(Checker{
		Name:        "thread_get_one_simple",
		Description: "",
		FnCheck:     Modifications(CheckThreadGetOneSimple),
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "thread_get_one_notfound",
		Description: "",
		FnCheck:     CheckThreadGetOneNotFound,
		Deps: []string{
			"thread_get_one_simple",
		},
	})
}

func CheckThreadGetOneSimple(c *client.Forum, m *Modify) {
	expected := CreateThread(c, nil, nil, nil)

	// Slug or ID
	id := m.SlugOrId(expected)

	// Check
	c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(id).
		WithContext(Expected(200, &expected, nil)))

	CheckThread(c, expected)
}

func CheckThreadGetOneNotFound(c *client.Forum) {
	thread := RandomThread()
	_, err := c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(thread.Slug).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewThreadGetOneNotFound(), err)

	_, err = c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(THREAD_FAKE_ID).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewThreadGetOneNotFound(), err)
}
