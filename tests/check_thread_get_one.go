package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
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

func CheckThreadGetOneSimple(c *client.Forum, m *Modify) {
	expected := CreateThread(c, nil, nil, nil)

	// Slug or ID
	id := m.SlugOrId(expected)

	// Check
	c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(id).
		WithContext(Expected(200, expected, filterThread)))

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

func PerfThreadGetOneSuccess(p *Perf) {
	thread := p.data.GetThread(-1)
	version := thread.Version
	slugOrId := GetSlugOrId(thread.Slug, int64(thread.ID))
	result, err := p.c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(slugOrId).
		WithContext(Expected(200, nil, nil)))
	CheckNil(err)

	p.Validate(func(v PerfValidator) {
		payload := result.Payload
		v.CheckInt(int(thread.ID), int(payload.ID), "Incorrect ID")
		v.CheckStr(thread.Forum.Slug, payload.Forum, "Incorrect Forum")
		v.CheckStr(thread.Slug, payload.Slug, "Incorrect Slug")
		v.CheckStr(thread.Author.Nickname, payload.Author, "Incorrect Author")
		v.CheckHash(thread.MessageHash, payload.Message, "Incorrect Message")
		v.CheckHash(thread.TitleHash, payload.Title, "Incorrect Title")
		v.CheckDate(&thread.Created, payload.Created, "Incorrect Created")
		v.CheckInt(int(thread.Votes), int(payload.Votes), "Incorrect Votes")
		v.Finish(version, thread.Version)
	})
}

func PerfThreadGetOneNotFound(p *Perf) {
	thread := RandomThread()
	for {
		thread.ID = rand.Int31n(100000000)
		if p.data.GetThreadById(thread.ID) == nil {
			break
		}
	}
	slugOrId := GetSlugOrId(thread.Slug, int64(thread.ID))
	_, err := p.c.Operations.ThreadGetOne(operations.NewThreadGetOneParams().
		WithSlugOrID(slugOrId).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewThreadGetOneNotFound(), err)
}
