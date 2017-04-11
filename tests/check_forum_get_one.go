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
	PerfRegister(PerfTest{
		Name:   "forum_get_one_success",
		Mode:   ModeRead,
		Weight: WeightNormal,
		FnPerf: PerfForumGetOneSuccess,
	})
	PerfRegister(PerfTest{
		Name:   "forum_get_one_not_found",
		Mode:   ModeRead,
		Weight: WeightRare,
		FnPerf: PerfForumGetOneNotFound,
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
	f1 := CreateForum(c, nil, nil)
	f2 := CreateForum(c, nil, nil)

	t1 := CreateThread(c, nil, f1, nil)
	CreatePosts(c, RandomPosts(3), t1)
	t2 := CreateThread(c, nil, f1, nil)
	CreatePosts(c, RandomPosts(5), t2)
	CreatePosts(c, RandomPosts(2), t1)
	t3 := CreateThread(c, nil, f2, nil)
	CreatePosts(c, RandomPosts(4), t3)

	f1.Threads = 2
	f1.Posts = 10
	CheckForum(c, f1)

	f2.Threads = 1
	f2.Posts = 4
	CheckForum(c, f2)
}

func PerfForumGetOneSuccess(p *Perf) {
	forum := p.data.GetForum(-1)
	version := forum.Version
	result, err := p.c.Operations.ForumGetOne(operations.NewForumGetOneParams().
		WithSlug(forum.Slug).
		WithContext(Expected(200, nil, nil)))
	CheckNil(err)

	p.Validate(func(v PerfValidator) {
		payload := result.Payload
		v.CheckInt(forum.Posts, int(payload.Posts), "Incorrect Posts count")
		v.CheckInt(forum.Threads, int(payload.Threads), "Incorrect Threads count")
		v.Finish(version, forum.Version)
	})
}

func PerfForumGetOneNotFound(p *Perf) {
	slug := RandomForum().Slug
	_, err := p.c.Operations.ForumGetOne(operations.NewForumGetOneParams().
		WithSlug(slug).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewForumGetOneNotFound(), err)
}
