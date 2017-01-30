package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/go-openapi/strfmt"
	"strings"
)

func init() {
	Register(Checker{
		Name:        "user_update_simple",
		Description: "",
		FnCheck:     CheckUserUpdateSimple,
		Deps: []string{
			"user_create_simple",
		},
	})
	Register(Checker{
		Name:        "user_update_empty",
		Description: "",
		FnCheck:     CheckUserUpdateEmpty,
		Deps: []string{
			"user_update_simple",
		},
	})
	Register(Checker{
		Name:        "user_update_part",
		Description: "",
		FnCheck:     CheckUserUpdatePart,
		Deps: []string{
			"user_update_simple",
		},
	})
	Register(Checker{
		Name:        "user_update_notfound",
		Description: "",
		FnCheck:     CheckUserUpdateNotFound,
		Deps: []string{
			"user_update_simple",
		},
	})
	Register(Checker{
		Name:        "user_update_conflict",
		Description: "",
		FnCheck:     CheckUserUpdateConflict,
		Deps: []string{
			"user_update_simple",
		},
	})
}

func CheckUserUpdateSimple(c *client.Forum) {
	user := CreateUser(c, nil)

	update := RandomUser()
	update.Nickname = ""

	expected := *update
	expected.Nickname = user.Nickname

	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user.Nickname).
		WithProfile(update).
		WithContext(Expected(200, &expected, nil)))

	CheckUser(c, &expected)
}

func CheckUserUpdateEmpty(c *client.Forum) {
	user := CreateUser(c, nil)

	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user.Nickname).
		WithProfile(&models.User{}).
		WithContext(Expected(200, user, nil)))

	CheckUser(c, user)
}

func CheckUserUpdatePart(c *client.Forum) {
	pass := 0
	for true {
		pass++
		Checkpoint(c, fmt.Sprintf("Pass %d", pass))

		fake := RandomUser()
		expected := CreateUser(c, nil)
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
			WithContext(Expected(200, &expected, nil)))

		CheckUser(c, expected)
	}
}

func CheckUserUpdateNotFound(c *client.Forum) {
	user := RandomUser()
	_, err := c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user.Nickname).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewUserUpdateNotFound(), err)
}

func CheckUserUpdateConflict(c *client.Forum) {
	user1 := CreateUser(c, nil)
	user2 := CreateUser(c, nil)

	update := &models.User{
		Email: user1.Email,
	}

	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user2.Nickname).
		WithProfile(update).
		WithContext(Expected(409, nil, nil)))

	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(strings.ToLower(user2.Nickname)).
		WithProfile(update).
		WithContext(Expected(409, nil, nil)))

	update.Email = strfmt.Email(strings.ToLower(update.Email.String()))
	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user2.Nickname).
		WithProfile(update).
		WithContext(Expected(409, nil, nil)))

	CheckUser(c, user1)
	CheckUser(c, user2)
}
