package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"math/rand"
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
	PerfRegister(PerfTest{
		Name:   "thread_get_one_success",
		Mode:   ModeRead,
		Weight: WeightNormal,
		FnPerf: PerfThreadGetOneSuccess,
	})
	PerfRegister(PerfTest{
		Name:   "thread_get_one_not_found",
		Mode:   ModeRead,
		Weight: WeightRare,
		FnPerf: PerfThreadGetOneNotFound,
	})
}

func CheckThreadGetOneSimple(c *client.Forum, f *Factory, m *Modify) {
	expected := f.CreateThread(c, nil, nil, nil)

	// Slug or ID
	id := m.SlugOrId(expected)

	// Check
	c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(id).
		WithContext(Expected(200, expected, filterThread)))

	CheckThread(c, expected)
}

func CheckThreadGetOneNotFound(c *client.Forum, f *Factory) {
	thread := f.RandomThread()
	_, err := c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(thread.Slug).
		WithContext(ExpectedError(404, "Can't find thread by slug: %s", thread.Slug)))
	CheckIsType(operations.NewThreadGetOneNotFound(), err)

	_, err = c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(THREAD_FAKE_ID).
		WithContext(ExpectedError(404, "Can't find thread by id: %d", THREAD_FAKE_ID)))
	CheckIsType(operations.NewThreadGetOneNotFound(), err)
}

func (self *PThread) Validate(v PerfValidator, thread *models.Thread, version PVersion) {
	v.CheckInt32(self.ID, thread.ID, "ID")
	v.CheckStr(self.Forum.Slug, thread.Forum, "Forum")
	v.CheckStr(self.Slug, thread.Slug, "Slug")
	v.CheckStr(self.Author.Nickname, thread.Author, "Author")
	v.CheckHash(self.MessageHash, thread.Message, "Message")
	v.CheckHash(self.TitleHash, thread.Title, "Title")
	v.CheckDate(&self.Created, thread.Created, "Created")
	v.CheckInt32(self.Votes, thread.Votes, "Votes")
	v.Finish(version, self.Version)
}

func PerfThreadGetOneSuccess(p *Perf, f *Factory) {
	thread := p.data.GetThread(-1, POST_POWER, 0.5)
	version := thread.Version
	slugOrId := GetSlugOrId(thread.Slug, int64(thread.ID))
	s := p.Session()
	result, err := p.c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(slugOrId).
		WithContext(s.Expected(200)))
	CheckNil(err)

	s.Validate(func(v PerfValidator) {
		thread.Validate(v, result.Payload, version)
	})
}

func PerfThreadGetOneNotFound(p *Perf, f *Factory) {
	var id int32
	slug := f.RandomSlug()
	for {
		id = rand.Int31n(100000000)
		if p.data.GetThreadById(id) == nil {
			break
		}
	}
	slugOrId := GetSlugOrId(slug, int64(id))
	_, err := p.c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(slugOrId).
		WithContext(ExpectedError(404, "Can't find thread by slug or id: %s", slugOrId)))
	CheckIsType(operations.NewThreadGetOneNotFound(), err)
}
