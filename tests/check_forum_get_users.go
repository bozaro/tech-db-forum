package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"math/rand"
	"sort"
	"time"
)

func init() {
	Register(Checker{
		Name:        "forum_get_users_simple",
		Description: "",
		FnCheck:     Modifications(CheckForumGetUsersSimple),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_users_notfound",
		Description: "",
		FnCheck:     Modifications(CheckForumGetUsersNotFound),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_users_empty",
		Description: "",
		FnCheck:     Modifications(CheckForumGetUsersEmpty),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_users_vote",
		Description: "",
		FnCheck:     Modifications(CheckForumGetUsersVote),
		Deps: []string{
			"posts_create_simple",
			"thread_create_vote_simple",
		},
	})
}

func CheckForumGetUsersSimple(c *client.Forum, m *Modify) {
	forum := CreateForum(c, nil, nil)
	threads := []*models.Thread{}
	users := []*models.User{}
	created := time.Now()
	created.Round(time.Millisecond)

	forum_users := map[string]*models.User{}
	// Пост, который не участвует в данном форуме
	CreatePost(c, nil, nil)
	// Пользователи
	for i := 0; i < 8; i++ {
		user := CreateUser(c, nil)
		users = append(users, user)
	}
	// Ветви
	for i := 0; i < 4; i++ {
		user := users[i/2]
		thread := CreateThread(c, nil, forum, user)
		threads = append(threads, thread)
		forum_users[user.Nickname] = user
	}
	// Посты
	for i := 0; i < 10; i++ {
		user := users[i%(len(users)-1)+1]
		thread := threads[rand.Intn(len(threads))]
		post := RandomPost()
		post.Author = user.Nickname
		CreatePost(c, post, thread)
		forum_users[user.Nickname] = user
	}

	// Список пользователей
	all_expected := []models.User{}
	for _, user := range forum_users {
		all_expected = append(all_expected, *user)
	}
	sort.Sort(UserByNickname(all_expected))

	// Desc
	desc := m.NullableBool()
	small := time.Millisecond
	if desc != nil && *desc {
		small = -small
		sort.Sort(sort.Reverse(UserByNickname(all_expected)))
	}

	// Check read all
	c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(forum.Slug).
		WithDesc(desc).
		WithContext(Expected(200, &all_expected, nil)))

	// Check read by 4 records
	limit := int32(4)
	for n := 0; n < len(all_expected); n += int(limit) - 1 {
		m := n + int(limit)
		if m > len(all_expected) {
			m = len(all_expected)
		}
		expected := all_expected[n:m]
		var since *string = nil
		if n > 0 {
			since = &all_expected[n-1].Nickname
		}
		c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
			WithSlug(forum.Slug).
			WithLimit(&limit).
			WithDesc(desc).
			WithSince(since).
			WithContext(Expected(200, &expected, nil)))
	}

	// Check read after all
	c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(forum.Slug).
		WithLimit(&limit).
		WithDesc(desc).
		WithSince(&all_expected[len(all_expected)-1].Nickname).
		WithContext(Expected(200, &[]models.Thread{}, nil)))
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
	author := CreateUser(c, nil)
	thread := CreateThread(c, nil, forum, author)
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
		WithContext(Expected(200, &[]models.User{*author}, nil)))
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
