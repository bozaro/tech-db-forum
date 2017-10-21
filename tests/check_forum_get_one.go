package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
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

func CheckForumGetOneSimple(c *client.Forum, f *Factory) {
	forum := f.CreateForum(c, nil, nil)
	CheckForum(c, forum)
}

func CheckForumGetOneNotFound(c *client.Forum, f *Factory) {
	forum := f.RandomForum()
	_, err := c.Operations.ForumGetOne(operations.NewForumGetOneParams().
		WithSlug(forum.Slug).
		WithContext(ExpectedError(404, "Can't find forum with slug: %s", forum.Slug)))
	CheckIsType(operations.NewForumGetOneNotFound(), err)
}

func CheckForumGetOneNocase(c *client.Forum, f *Factory, m *Modify) {
	forum := f.CreateForum(c, nil, nil)
	slug := m.Case(forum.Slug)
	_, err := c.Operations.ForumGetOne(operations.NewForumGetOneParams().
		WithSlug(slug).
		WithContext(Expected(200, forum, nil)))
	CheckNil(err)
}

func CheckForumGetOneCounter(c *client.Forum, f *Factory) {
	f1 := f.CreateForum(c, nil, nil)
	f2 := f.CreateForum(c, nil, nil)

	t1 := f.CreateThread(c, nil, f1, nil)
	f.CreatePosts(c, f.RandomPosts(3), t1)
	t2 := f.CreateThread(c, nil, f1, nil)
	f.CreatePosts(c, f.RandomPosts(5), t2)
	f.CreatePosts(c, f.RandomPosts(2), t1)
	t3 := f.CreateThread(c, nil, f2, nil)
	f.CreatePosts(c, f.RandomPosts(4), t3)

	f1.Threads = 2
	f1.Posts = 10
	CheckForum(c, f1)

	f2.Threads = 1
	f2.Posts = 4
	CheckForum(c, f2)
}

func (self *PForum) Validate(v PerfValidator, forum *models.Forum, version PVersion) {
	v.CheckInt64(self.Posts, forum.Posts, "Posts")
	v.CheckInt32(self.Threads, forum.Threads, "Threads")
	v.Finish(version, self.Version)
}

func PerfForumGetOneSuccess(p *Perf, f *Factory) {
	forum := p.data.GetForum(-1)
	version := forum.Version
	s := p.Session()
	result, err := p.c.Operations.ForumGetOne(operations.NewForumGetOneParams().
		WithSlug(forum.Slug).
		WithContext(s.Expected(200)))
	CheckNil(err)

	s.Validate(func(v PerfValidator) {
		forum.Validate(v, result.Payload, version)
	})
}

func PerfForumGetOneNotFound(p *Perf, f *Factory) {
	slug := f.RandomSlug()
	_, err := p.c.Operations.ForumGetOne(operations.NewForumGetOneParams().
		WithSlug(slug).
		WithContext(ExpectedError(404, "Can't find forum with slug: %s", slug)))
	CheckIsType(operations.NewForumGetOneNotFound(), err)
}
