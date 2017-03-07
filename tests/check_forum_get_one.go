package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
)

func init() {
	Register(Checker{
		Name:        "forum_get_one_simple",
		Description: "",
		FnCheck:     CheckForumGetOneSimple,
		Deps: []string{
			"forum_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_one_notfound",
		Description: "",
		FnCheck:     CheckForumGetOneNotFound,
		Deps: []string{
			"forum_get_one_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_one_nocase",
		Description: "",
		FnCheck:     Modifications(CheckForumGetOneNocase),
		Deps: []string{
			"forum_get_one_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_one_counter",
		Description: "",
		FnCheck:     CheckForumGetOneCounter,
		Deps: []string{
			"posts_create_simple",
		},
	})
}

func CheckForumGetOneSimple(c *client.Forum) {
	forum := CreateForum(c, nil, nil)
	CheckForum(c, forum)
}

func CheckForumGetOneNotFound(c *client.Forum) {
	forum := RandomForum()
	_, err := c.Operations.ForumGetOne(operations.NewForumGetOneParams().
		WithSlug(forum.Slug).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewForumGetOneNotFound(), err)
}

func CheckForumGetOneNocase(c *client.Forum, m *Modify) {
	forum := CreateForum(c, nil, nil)
	slug := m.Case(forum.Slug)
	_, err := c.Operations.ForumGetOne(operations.NewForumGetOneParams().
		WithSlug(slug).
		WithContext(Expected(200, forum, nil)))
	CheckNil(err)
}

func CheckForumGetOneCounter(c *client.Forum) {
	forum := CreateForum(c, nil, nil)

	t1 := CreateThread(c, nil, forum, nil)
	CreatePosts(c, RandomPosts(3), t1)
	t2 := CreateThread(c, nil, forum, nil)
	CreatePosts(c, RandomPosts(5), t2)
	CreatePosts(c, RandomPosts(2), t1)

	forum.Threads = 2
	forum.Posts = 10
	CheckForum(c, forum)
}
