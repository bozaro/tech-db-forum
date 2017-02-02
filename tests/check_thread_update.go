package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
)

func init() {
	Register(Checker{
		Name:        "thread_update_simple",
		Description: "",
		FnCheck:     CheckThreadUpdateSimple,
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
		FnCheck:     CheckThreadUpdatePart,
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

func CheckThreadUpdateSimple(c *client.Forum) {
	for pass := 0; pass < 2; pass++ {
		Checkpoint(c, fmt.Sprintf("Pass %d", pass+1))
		thread := CreateThread(c, nil, nil, nil)

		temp := RandomThread()
		update := models.ThreadUpdate{}
		update.Title = temp.Title
		update.Message = temp.Message

		expected := *thread
		expected.Title = update.Title
		expected.Message = update.Message

		id := thread.Slug
		if pass == 1 {
			id = fmt.Sprintf("%d", thread.ID)
		}

		c.Operations.ThreadUpdate(operations.NewThreadUpdateParams().
			WithSlugOrID(id).
			WithThread(&update).
			WithContext(Expected(200, &expected, nil)))

		CheckThread(c, &expected)
	}
}

func CheckThreadUpdateEmpty(c *client.Forum) {
	thread := CreateThread(c, nil, nil, nil)

	c.Operations.ThreadUpdate(operations.NewThreadUpdateParams().
		WithSlugOrID(thread.Slug).
		WithThread(&models.ThreadUpdate{}).
		WithContext(Expected(200, thread, nil)))

	CheckThread(c, thread)
}

func CheckThreadUpdatePart(c *client.Forum) {
	pass := 0
	for true {
		pass++
		Checkpoint(c, fmt.Sprintf("Pass %d", pass))

		fake := RandomThread()
		expected := CreateThread(c, nil, nil, nil)
		update := &models.ThreadUpdate{}

		modify := pass
		// Slug or ID
		id := expected.Slug
		if (modify & 1) == 1 {
			id = fmt.Sprintf("%d", expected.ID)
		}
		modify >>= 1
		// Title
		if (modify & 1) == 1 {
			expected.Title = fake.Title
			update.Title = fake.Title
		}
		modify >>= 1
		// Message
		if (modify & 1) == 1 {
			expected.Message = fake.Message
			update.Message = fake.Message
		}
		modify >>= 1
		// Done?
		if modify != 0 {
			break
		}
		// Check
		c.Operations.ThreadUpdate(operations.NewThreadUpdateParams().
			WithSlugOrID(id).
			WithThread(update).
			WithContext(Expected(200, &expected, nil)))

		CheckThread(c, expected)
	}
}

func CheckThreadUpdateNotFound(c *client.Forum) {
	thread := RandomThread()
	_, err := c.Operations.ThreadUpdate(operations.NewThreadUpdateParams().
		WithSlugOrID(thread.Slug).
		WithThread(&models.ThreadUpdate{
			Title:   thread.Title,
			Message: thread.Message,
		}).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewThreadUpdateNotFound(), err)

	_, err = c.Operations.ThreadUpdate(operations.NewThreadUpdateParams().
		WithSlugOrID(THREAD_FAKE_ID).
		WithThread(&models.ThreadUpdate{
			Title:   thread.Title,
			Message: thread.Message,
		}).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewThreadUpdateNotFound(), err)
}
