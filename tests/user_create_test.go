package main

import (
	"github.com/bozaro/tech-db-forum/tests/client/operations"
	"github.com/bozaro/tech-db-forum/tests/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func CreateUser(t *testing.T, user *models.User) *models.User {
	if user == nil {
		user = RandomUser()
	}

	request := *user
	request.Nickname = ""

	_, err := c.Operations.UserCreate(operations.NewUserCreateParams().
		WithNickname(user.Nickname).
		WithProfile(&request).
		WithContext(Expected(t, 201, user, nil)))
	assert.Nil(t, err)

	return user
}

func CheckUser(t *testing.T, user *models.User) {
	_, err := c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(user.Nickname).
		WithContext(Expected(t, 200, user, nil)))
	assert.Nil(t, err)
}
