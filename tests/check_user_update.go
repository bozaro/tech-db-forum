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

func CheckUserUpdateSimple(c *client.Forum, f *Factory) {
	user := f.CreateUser(c, nil)

	expected := f.RandomUser()
	expected.Nickname = user.Nickname
	update := models.UserUpdate{
		About:    expected.About,
		Email:    expected.Email,
		Fullname: expected.Fullname,
	}

	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user.Nickname).
		WithProfile(&update).
		WithContext(Expected(200, expected, nil)))

	CheckUser(c, expected)
}

func CheckUserUpdateEmpty(c *client.Forum, f *Factory) {
	user := f.CreateUser(c, nil)

	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user.Nickname).
		WithProfile(&models.UserUpdate{}).
		WithContext(Expected(200, user, nil)))

	CheckUser(c, user)
}

func CheckUserUpdatePart(c *client.Forum, f *Factory, m *Modify) {
	fake := f.RandomUser()
	expected := f.CreateUser(c, nil)
	update := &models.UserUpdate{}

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
		WithContext(Expected(200, expected, nil)))

	CheckUser(c, expected)
}

func CheckUserUpdateNotFound(c *client.Forum, f *Factory) {
	user := f.RandomUser()
	update := &models.UserUpdate{}
	update.Fullname = user.Fullname

	_, err := c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(user.Nickname).
		WithProfile(update).
		WithContext(ExpectedError(404, "Can't find user by nickname: %s", user.Nickname)))
	CheckIsType(operations.NewUserUpdateNotFound(), err)
}

func CheckUserUpdateConflict(c *client.Forum, f *Factory, m *Modify) {
	user1 := f.CreateUser(c, nil)
	user2 := f.CreateUser(c, nil)

	update := &models.UserUpdate{
		Email: strfmt.Email(m.Case(user1.Email.String())),
	}

	c.Operations.UserUpdate(operations.NewUserUpdateParams().
		WithNickname(m.Case(user2.Nickname)).
		WithProfile(update).
		WithContext(ExpectedError(409, "This email is already registered by user: %s", user1.Nickname)))

	CheckUser(c, user1)
	CheckUser(c, user2)
}
