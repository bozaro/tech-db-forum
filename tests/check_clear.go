package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
)

func init() {
	Register(Checker{
		Name:        "clear",
		Description: "",
		FnCheck:     CheckClear,
		Deps: []string{
			"status",
		},
	})
}

func CheckClear(c *client.Forum, f *Factory) {
	f.CreatePost(c, nil, nil)
	var err error
	_, err = c.Operations.Clear(operations.NewClearParams().
		WithContext(Expected(200, nil, nil)))
	CheckNil(err)
	_, err = c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(200, &models.Status{}, nil)))
	CheckNil(err)
}
