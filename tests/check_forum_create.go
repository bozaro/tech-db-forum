package tests

import (
	"github.com/bozaro/tech-db-forum/client"
	"github.com/bozaro/tech-db-forum/client/operations"
	"github.com/bozaro/tech-db-forum/models"
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
	CreateUser(c, nil)
}

func CheckForumCreateUnicode(c *client.Forum) {
	forum := RandomForum()
	forum.Title = "Обсужение Unicode 😋"
	CreateForum(c, forum, nil)
	CheckForum(c, forum)
}
