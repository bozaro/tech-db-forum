package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"strings"
)

func init() {
	Register(Checker{
		Name:        "thread_get_one_simple",
		Description: "",
		FnCheck:     CheckThreadGetOneSimple,
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

func CheckThreadGetOneSimple(c *client.Forum) {
	pass := 0
	for true {
		pass++
		Checkpoint(c, fmt.Sprintf("Pass %d", pass))

		expected := CreateThread(c, nil, nil, nil)

		modify := pass
		// Slug or ID
		id := expected.Slug
		switch modify & 2 {
		case 1:
			id = fmt.Sprintf("%d", expected.ID)
		case 2:
			id = strings.ToLower(expected.Slug)
		case 3:
			id = strings.ToUpper(expected.Slug)
		}
		modify >>= 2
		// Done?
		if modify != 0 {
			break
		}
		// Check
		c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
			WithSlugOrID(id).
			WithContext(Expected(200, &expected, nil)))

		CheckThread(c, expected)
	}
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
