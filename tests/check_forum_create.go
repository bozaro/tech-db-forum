package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/client"
	"github.com/bozaro/tech-db-forum/client/operations"
	"github.com/bozaro/tech-db-forum/models"
	"strings"
)

func init() {
	Register(Checker{
		Name:        "forum_create_simple",
		Description: "",
		FnCheck:     CheckForumCreateSimple,
	})
	Register(Checker{
		Name:        "forum_create_unicode",
		Description: "",
		FnCheck:     CheckForumCreateUnicode,
		Deps: []string{
			"forum_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_create_conflict",
		Description: "",
		FnCheck:     CheckForumCreateConflict,
		Deps: []string{
			"forum_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_create_user_case",
		Description: "",
		FnCheck:     CheckForumCreateUserCase,
		Deps: []string{
			"forum_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_create_user_notfound",
		Description: "",
		FnCheck:     CheckForumCreateUserNotFound,
		Deps: []string{
			"forum_create_simple",
		},
	})
}

func CreateForum(c *client.Forum, forum *models.Forum, owner *models.User) *models.Forum {
	if forum == nil {
		forum = RandomForum()
	}
	if forum.User == "" {
		if owner == nil {
			owner = CreateUser(c, nil)
		}
		forum.User = owner.Nickname
	}

	_, err := c.Operations.ForumCreate(operations.NewForumCreateParams().
		WithForum(forum).
		WithContext(Expected(201, forum, nil)))
	CheckNil(err)

	return forum
}

func CheckForum(c *client.Forum, forum *models.Forum) {
	_, err := c.Operations.ForumGetOne(operations.NewForumGetOneParams().
		WithSlug(forum.Slug).
		WithContext(Expected(200, forum, nil)))
	CheckNil(err)
}

func CheckForumCreateSimple(c *client.Forum) {
	CreateForum(c, nil, nil)
}

func CheckForumCreateUserCase(c *client.Forum) {
	pass := 0
	for true {
		pass++
		Checkpoint(c, fmt.Sprintf("Pass %d", pass))
		user := CreateUser(c, nil)
		forum := RandomForum()

		modify := pass
		// Slug
		switch modify % 3 {
		case 1:
			forum.User = strings.ToLower(user.Nickname)
		case 2:
			forum.User = strings.ToUpper(user.Nickname)
		}
		modify /= 3
		// Done?
		if modify != 0 {
			break
		}

		expected := *forum
		expected.User = user.Nickname
		_, err := c.Operations.ForumCreate(operations.NewForumCreateParams().
			WithForum(forum).
			WithContext(Expected(201, &expected, nil)))
		CheckNil(err)

		CheckForum(c, &expected)
	}
}

func CheckForumCreateUserNotFound(c *client.Forum) {
	user := RandomUser()
	forum := RandomForum()
	forum.User = user.Nickname

	_, err := c.Operations.ForumCreate(operations.NewForumCreateParams().
		WithForum(forum).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(err, operations.NewForumCreateNotFound())
}

func CheckForumCreateUnicode(c *client.Forum) {
	forum := RandomForum()
	forum.Title = "Обсужение Unicode 😋"
	CreateForum(c, forum, nil)
	CheckForum(c, forum)
}

func CheckForumCreateConflict(c *client.Forum) {
	pass := 0
	for true {
		pass++
		Checkpoint(c, fmt.Sprintf("Pass %d", pass))

		forum := CreateForum(c, nil, nil)
		conflict_forum := RandomForum()
		conflict_forum.User = forum.User

		modify := pass
		// Slug
		switch modify % 4 {
		case 1:
			conflict_forum.Slug = forum.Slug
		case 2:
			conflict_forum.Slug = strings.ToLower(forum.Slug)
		case 3:
			conflict_forum.Slug = strings.ToUpper(forum.Slug)
		}
		modify /= 4
		// Done?
		if modify != 0 {
			break
		}
		// Check
		_, err := c.Operations.ForumCreate(operations.NewForumCreateParams().
			WithForum(conflict_forum).
			WithContext(Expected(409, &forum, nil)))
		CheckIsType(operations.NewForumCreateConflict(), err)
	}
}
