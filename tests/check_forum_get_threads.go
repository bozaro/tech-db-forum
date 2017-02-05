package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/go-openapi/strfmt"
	"sort"
	"time"
)

func init() {
	Register(Checker{
		Name:        "forum_get_threads_simple",
		Description: "",
		FnCheck:     Modifications(CheckForumGetThreadsSimple),
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_threads_notfound",
		Description: "",
		FnCheck:     Modifications(CheckForumGetThreadsNotFound),
		Deps: []string{
			"thread_create_simple",
		},
	})
}

func CheckForumGetThreadsSimple(c *client.Forum, m *Modify) {
	forum := CreateForum(c, nil, nil)
	threads := []models.Thread{}
	created := time.Now()
	created.Round(time.Millisecond)
	for i := 0; i < 10; i++ {
		thread := CreateThread(c, nil, forum, nil)
		threads = append(threads, *thread)
	}
	sort.Sort(ThreadByCreated(threads))

	var desc *bool

	// Desc
	small := time.Millisecond
	switch m.Int(3) {
	case 1:
		v := bool(true)
		small = -small
		desc = &v
		sort.Sort(sort.Reverse(ThreadByCreated(threads)))
	case 2:
		v := bool(false)
		desc = &v
	}

	// Check read all
	c.Operations.ForumGetThreads(operations.NewForumGetThreadsParams().
		WithSlug(forum.Slug).
		WithDesc(desc).
		WithContext(Expected(200, &threads, nil)))

	// Check read by 4 records
	limit := int32(4)
	var since *strfmt.DateTime = nil
	for n := 0; n < len(threads); n += int(limit) - 1 {
		m := n + int(limit)
		if m > len(threads) {
			m = len(threads)
		}
		expected := threads[n:m]
		c.Operations.ForumGetThreads(operations.NewForumGetThreadsParams().
			WithSlug(forum.Slug).
			WithLimit(&limit).
			WithDesc(desc).
			WithSince(since).
			WithContext(Expected(200, &expected, nil)))
		since = &threads[m-1].Created
	}

	// Check read after all
	after_last := strfmt.DateTime(time.Time(threads[len(threads)-1].Created).Add(small))
	c.Operations.ForumGetThreads(operations.NewForumGetThreadsParams().
		WithSlug(forum.Slug).
		WithLimit(&limit).
		WithDesc(desc).
		WithSince(&after_last).
		WithContext(Expected(200, &[]models.Thread{}, nil)))
}

func CheckForumGetThreadsNotFound(c *client.Forum, m *Modify) {
	var limit *int32
	var since *strfmt.DateTime
	var desc *bool

	forum := RandomForum()
	// Limit
	if m.Bool() {
		v := int32(10)
		limit = &v
	}
	// Since
	if m.Bool() {
		v := strfmt.DateTime(time.Now())
		since = &v
	}
	// Desc
	switch m.Int(3) {
	case 1:
		v := bool(true)
		desc = &v
	case 2:
		v := bool(false)
		desc = &v
	}

	// Check
	_, err := c.Operations.ForumGetThreads(operations.NewForumGetThreadsParams().
		WithSlug(forum.Slug).
		WithLimit(limit).
		WithSince(since).
		WithDesc(desc).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewForumGetThreadsNotFound(), err)
}
