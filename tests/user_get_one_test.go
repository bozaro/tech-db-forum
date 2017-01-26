package main

import (
	"github.com/bozaro/tech-db-forum/tests/client/operations"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserGetOneSimple(t *testing.T) {
	user := CreateUser(t, nil)
	CheckUser(t, user)
}

func TestUserGetOneNotFound(t *testing.T) {
	user := RandomUser()
	_, err := c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(user.Nickname).
		WithContext(Expected(t, 404, nil, nil)))
	assert.IsType(t, operations.NewUserGetOneNotFound(), err)
}
