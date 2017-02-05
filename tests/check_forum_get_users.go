package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
)

func init() {
	Register(Checker{
		Name:        "forum_get_users_simple",
		Description: "",
		FnCheck:     Modifications(CheckForumGetUsersSimple),
		Deps: []string{
			"post_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_users_notfound",
		Description: "",
		FnCheck:     Modifications(CheckForumGetUsersNotFound),
		Deps: []string{
			"post_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_users_empty",
		Description: "",
		FnCheck:     Modifications(CheckForumGetUsersEmpty),
		Deps: []string{
			"post_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_users_vote",
		Description: "",
		FnCheck:     Modifications(CheckForumGetUsersVote),
		Deps: []string{
			"post_create_simple",
			"thread_create_vote_simple",
		},
	})
}

func CheckForumGetUsersSimple(c *client.Forum, m *Modify) {
	/*forum := CreateForum(c, nil, nil)
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
	c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
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
		c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
			WithSlug(forum.Slug).
			WithLimit(&limit).
			WithDesc(desc).
			WithSince(since).
			WithContext(Expected(200, &expected, nil)))
		since = &threads[m - 1].Created
	}

	// Check read after all
	after_last := strfmt.DateTime(time.Time(threads[len(threads) - 1].Created).Add(small))
	c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(forum.Slug).
		WithLimit(&limit).
		WithDesc(desc).
		WithSince(&after_last).
		WithContext(Expected(200, &[]models.Thread{}, nil)))*/
}

func CheckForumGetUsersEmpty(c *client.Forum, m *Modify) {
	var limit *int32
	var since *string
	var desc *bool

	forum := CreateForum(c, nil, nil)
	// Limit
	if m.Bool() {
		v := int32(10)
		limit = &v
	}
	// Since
	if m.Bool() {
		v := RandomNickname()
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
	c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(forum.Slug).
		WithLimit(limit).
		WithSince(since).
		WithDesc(desc).
		WithContext(Expected(200, &[]models.User{}, nil)))
}

func CheckForumGetUsersVote(c *client.Forum, m *Modify) {
	var limit *int32
	var since *string
	var desc *bool

	forum := CreateForum(c, nil, nil)

	user := CreateUser(c, nil)
	thread := CreateThread(c, nil, nil, nil)
	var err error
	// Like user1
	thread.Votes = 1
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(m.SlugOrId(thread)).
		WithVote(&models.Vote{
			Nickname: user.Nickname,
			Voice:    1,
		}).
		WithContext(Expected(200, &thread, nil)))
	CheckNil(err)

	thread.Votes = 1
	CheckThread(c, thread)

	// Check
	c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(forum.Slug).
		WithLimit(limit).
		WithSince(since).
		WithDesc(desc).
		WithContext(Expected(200, &[]models.User{}, nil)))
}

func CheckForumGetUsersNotFound(c *client.Forum, m *Modify) {
	var limit *int32
	var since *string
	var desc *bool

	forum := RandomForum()
	// Limit
	if m.Bool() {
		v := int32(10)
		limit = &v
	}
	// Since
	if m.Bool() {
		v := RandomNickname()
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
	_, err := c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(forum.Slug).
		WithLimit(limit).
		WithSince(since).
		WithDesc(desc).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewForumGetUsersNotFound(), err)
}
