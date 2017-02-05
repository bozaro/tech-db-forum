package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/go-openapi/strfmt"
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
		FnCheck:     Modifications(CheckUserUpdatePart),
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
		FnCheck:     Modifications(CheckUserUpdateConflict),
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

func CheckUserUpdatePart(c *client.Forum, m *Modify) {
	fake := RandomUser()
	expected := CreateUser(c, nil)
	update := &models.User{}

	// Email
	if m.Bool() {
		expected.Email = fake.Email
		update.Email = fake.Email
	}
	// About
	if m.Bool() {
		expected.About = fake.About
		update.About = fake.About
	}
	// Fullname
	if m.Bool() {
		expected.Fullname = fake.Fullname
		update.Fullname = fake.Fullname
	}

	// Check
	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(expected.Nickname).
		WithProfile(update).
		WithContext(Expected(200, &expected, nil)))

	CheckUser(c, expected)
}

func CheckUserUpdateNotFound(c *client.Forum) {
	user := RandomUser()
	_, err := c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user.Nickname).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewUserUpdateNotFound(), err)
}

func CheckUserUpdateConflict(c *client.Forum, m *Modify) {
	user1 := CreateUser(c, nil)
	user2 := CreateUser(c, nil)

	update := &models.User{
		Email: strfmt.Email(m.Case(user1.Email.String())),
	}

	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(m.Case(user2.Nickname)).
		WithProfile(update).
		WithContext(Expected(409, nil, nil)))

	CheckUser(c, user1)
	CheckUser(c, user2)
}
