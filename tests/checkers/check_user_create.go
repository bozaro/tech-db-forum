package checkers

import (
	"github.com/bozaro/tech-db-forum/tests/client"
	"github.com/bozaro/tech-db-forum/tests/client/operations"
	"github.com/bozaro/tech-db-forum/tests/models"
)

func CreateUser(c *client.Forum, user *models.User) *models.User {
	if user == nil {
		user = RandomUser()
	}

	request := *user
	request.Nickname = ""

	_, err := c.Operations.UserCreate(operations.NewUserCreateParams().
		WithNickname(user.Nickname).
		WithProfile(&request).
		WithContext(Expected(201, user, nil)))
	CheckNil(err)

	return user
}

func CheckUser(c *client.Forum, user *models.User) {
	_, err := c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(user.Nickname).
		WithContext(Expected(200, user, nil)))
	CheckNil(err)
}
