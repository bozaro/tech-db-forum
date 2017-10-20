package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
)

func init() {
	Register(Checker{
		Name:        "status",
		Description: "",
		FnCheck:     CheckStatus,
		Deps: []string{
			"posts_create_simple",
		},
	})
	PerfRegister(PerfTest{
		Name:   "status",
		Mode:   ModeRead,
		Weight: WeightNever,
		FnPerf: PerfStatus,
	})
}

func CheckStatus(c *client.Forum, f *Factory) {
	result, err := c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(200, nil, nil)))
	CheckNil(err)

	status := result.Payload
	// Add single user
	user := f.CreateUser(c, nil)
	status.User++
	_, err = c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(200, status, nil)))
	CheckNil(err)

	// Add forum
	forum := f.CreateForum(c, nil, user)
	status.Forum++
	_, err = c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(200, status, nil)))
	CheckNil(err)

	// Add thread
	thread := f.CreateThread(c, nil, forum, user)
	status.Thread++
	_, err = c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(200, status, nil)))
	CheckNil(err)

	// Add post
	f.CreatePost(c, nil, thread)
	status.Post++
	_, err = c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(200, status, nil)))
	CheckNil(err)
}

func PerfStatus(p *Perf, f *Factory) {
	status := p.data.Status
	version := status.Version
	s := p.Session()
	result, err := p.c.Operations.Status(operations.NewStatusParams().
		WithContext(s.Expected(200)))
	CheckNil(err)

	s.Validate(func(v PerfValidator) {
		payload := result.Payload
		v.CheckInt32(status.Forum, payload.Forum, "Incorrect Forum count")
		v.CheckInt64(status.Post, payload.Post, "Incorrect Post count")
		v.CheckInt32(status.User, payload.User, "Incorrect User count")
		v.CheckInt32(status.Thread, payload.Thread, "Incorrect Thread count")
		v.Finish(version, status.Version)
	})
}
