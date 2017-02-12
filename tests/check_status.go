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
