package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"strings"
)

func init() {
	Register(Checker{
		Name:        "user_get_one_simple",
		Description: "",
		FnCheck:     CheckUserGetOneSimple,
		Deps: []string{
			"user_create_simple",
		},
	})
	Register(Checker{
		Name:        "user_get_one_notfound",
		Description: "",
		FnCheck:     CheckUserGetOneNotFound,
		Deps: []string{
			"user_get_one_simple",
		},
	})
	Register(Checker{
		Name:        "user_get_one_nocase",
		Description: "",
		FnCheck:     CheckUserGetOneNocase,
		Deps: []string{
			"user_get_one_simple",
		},
	})
}

func CheckUserGetOneSimple(c *client.Forum) {
	user := CreateUser(c, nil)
	CheckUser(c, user)
}

func CheckUserGetOneNotFound(c *client.Forum) {
	user := RandomUser()
	_, err := c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(user.Nickname).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewUserGetOneNotFound(), err)
}

func CheckUserGetOneNocase(c *client.Forum) {
	user := CreateUser(c, nil)
	_, err := c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(strings.ToLower(user.Nickname)).
		WithContext(Expected(200, user, nil)))
	CheckNil(err)
}
