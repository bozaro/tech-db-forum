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
		Weight: WeightRare,
		FnPerf: PerfStatus,
	})
}

func CheckStatus(c *client.Forum) {
	result, err := c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(200, nil, nil)))
	CheckNil(err)

	status := result.Payload
	// Add single user
	user := CreateUser(c, nil)
	status.User++
	_, err = c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(200, &status, nil)))
	CheckNil(err)

	// Add forum
	forum := CreateForum(c, nil, user)
	status.Forum++
	_, err = c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(200, &status, nil)))
	CheckNil(err)

	// Add thread
	thread := CreateThread(c, nil, forum, user)
	status.Thread++
	_, err = c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(200, &status, nil)))
	CheckNil(err)

	// Add post
	CreatePost(c, nil, thread)
	status.Post++
	_, err = c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(200, &status, nil)))
	CheckNil(err)
}

func PerfStatus(p *Perf) {
	global_old := p.data.Status(false)
	result, err := p.c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(200, nil, nil)))
	CheckNil(err)

	p.Validate(func(v PerfValidator) {
		global_new := p.data.Status(true)
		payload := result.Payload
		v.CheckBetween(int(global_old.Forum), int(payload.Forum), int(global_new.Forum), "Incorrect Forum count")
		v.CheckBetween(int(global_old.Post), int(payload.Post), int(global_new.Post), "Incorrect Post count")
		v.CheckBetween(int(global_old.User), int(payload.User), int(global_new.User), "Incorrect User count")
		v.CheckBetween(int(global_old.Thread), int(payload.Thread), int(global_new.Thread), "Incorrect Thread count")
	})
}
