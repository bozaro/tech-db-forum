package main

import (
	"github.com/bozaro/tech-db-forum/tests/client/operations"
	"github.com/bozaro/tech-db-forum/tests/models"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUserUpdateSimple(t *testing.T) {
	user := CreateUser(t, nil)

	update := RandomUser()
	update.Nickname = ""

	expected := *update
	expected.Nickname = user.Nickname

	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user.Nickname).
		WithProfile(update).
		WithContext(Expected(t, 200, &expected, nil)))

	CheckUser(t, &expected)
}

func TestUserUpdateEmpty(t *testing.T) {
	user := CreateUser(t, nil)

	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user.Nickname).
		WithProfile(&models.User{}).
		WithContext(Expected(t, 200, user, nil)))

	CheckUser(t, user)
}

func TestUserUpdatePart(t *testing.T) {
	pass := 0
	for true {
		pass++

		fake := RandomUser()
		expected := CreateUser(t, nil)
		update := &models.User{}

		modify := pass
		// Email
		if (modify & 1) == 1 {
			expected.Email = fake.Email
			update.Email = fake.Email
		}
		modify >>= 1
		// About
		if (modify & 1) == 1 {
			expected.About = fake.About
			update.About = fake.About
		}
		modify >>= 1
		// Fullname
		if (modify & 1) == 1 {
			expected.Fullname = fake.Fullname
			update.Fullname = fake.Fullname
		}
		modify >>= 1
		// Done?
		if modify != 0 {
			break
		}
		// Check
		c.Operations.UserUpdate(operations.NewUserUpdateParams().
			WithNickname(expected.Nickname).
			WithProfile(update).
			WithContext(Expected(t, 200, &expected, nil)))

		CheckUser(t, expected)
	}
}

func TestUserUpdateNotFound(t *testing.T) {
	user := RandomUser()
	_, err := c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user.Nickname).
		WithContext(Expected(t, 404, nil, nil)))
	assert.IsType(t, operations.NewUserUpdateNotFound(), err)
}

func TestUserUpdateConflict(t *testing.T) {
	user1 := CreateUser(t, nil)
	user2 := CreateUser(t, nil)

	update := &models.User{
		Email: user1.Email,
	}

	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user2.Nickname).
		WithProfile(update).
		WithContext(Expected(t, 409, nil, nil)))

	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(strings.ToLower(user2.Nickname)).
		WithProfile(update).
		WithContext(Expected(t, 409, nil, nil)))

	update.Email = strfmt.Email(strings.ToLower(update.Email.String()))
	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user2.Nickname).
		WithProfile(update).
		WithContext(Expected(t, 409, nil, nil)))

	CheckUser(t, user1)
	CheckUser(t, user2)
}
