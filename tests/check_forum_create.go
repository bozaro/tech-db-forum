package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
)

func init() {
	Register(Checker{
		Name:        "forum_create_simple",
		Description: "",
		FnCheck:     CheckForumCreateSimple,
		Deps: []string{
			"user_get_one_simple",
		},
	})
	Register(Checker{
		Name:        "forum_create_unicode",
		Description: "",
		FnCheck:     CheckForumCreateUnicode,
		Deps: []string{
			"forum_get_one_simple",
		},
	})
	Register(Checker{
		Name:        "forum_create_conflict",
		Description: "",
		FnCheck:     Modifications(CheckForumCreateConflict),
		Deps: []string{
			"forum_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_create_user_case",
		Description: "",
		FnCheck:     Modifications(CheckForumCreateUserCase),
		Deps: []string{
			"forum_get_one_simple",
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

func (f *Factory) CreateForum(c *client.Forum, forum *models.Forum, owner *models.User) *models.Forum {
	if forum == nil {
		forum = f.RandomForum()
	}
	if forum.User == "" {
		if owner == nil {
			owner = f.CreateUser(c, nil)
		}
		forum.User = owner.Nickname
	}

	result, err := c.Operations.ForumCreate(operations.NewForumCreateParams().
		WithForum(forum).
		WithContext(Expected(201, forum, nil)))
	CheckNil(err)
	if forum.Slug != result.Payload.Slug {
		log.Errorf("Unexpected created forum slug: %s -> %s", forum.Slug, result.Payload.Slug)
	}

	return result.Payload
}

func CheckForum(c *client.Forum, forum *models.Forum) {
	_, err := c.Operations.ForumGetOne(operations.NewForumGetOneParams().
		WithSlug(forum.Slug).
		WithContext(Expected(200, forum, nil)))
	CheckNil(err)
}

func CheckForumCreateSimple(c *client.Forum, f *Factory) {
	f.CreateForum(c, nil, nil)
}

func CheckForumCreateUserCase(c *client.Forum, f *Factory, m *Modify) {
	user := f.CreateUser(c, nil)
	forum := f.RandomForum()

	// Slug
	forum.User = m.Case(user.Nickname)

	// Check
	expected := *forum
	expected.User = user.Nickname
	_, err := c.Operations.ForumCreate(operations.NewForumCreateParams().
		WithForum(forum).
		WithContext(Expected(201, &expected, nil)))
	CheckNil(err)

	CheckForum(c, &expected)
}

func CheckForumCreateUserNotFound(c *client.Forum, f *Factory) {
	user := f.RandomUser()
	forum := f.RandomForum()
	forum.User = user.Nickname

	_, err := c.Operations.ForumCreate(operations.NewForumCreateParams().
		WithForum(forum).
		WithContext(ExpectedError(404, "Can't find user with nickname: %s", user.Nickname)))
	CheckIsType(err, operations.NewForumCreateNotFound())
}

func CheckForumCreateUnicode(c *client.Forum, f *Factory) {
	forum := f.RandomForum()
	forum.Title = "Обсужение Unicode 😋"
	f.CreateForum(c, forum, nil)
	CheckForum(c, forum)
}

func CheckForumCreateConflict(c *client.Forum, f *Factory, m *Modify) {
	forum := f.CreateForum(c, nil, nil)
	conflict_forum := f.RandomForum()
	conflict_forum.User = forum.User

	// Slug
	conflict_forum.Slug = m.Case(forum.Slug)

	// Check
	_, err := c.Operations.ForumCreate(operations.NewForumCreateParams().
		WithForum(conflict_forum).
		WithContext(Expected(409, forum, nil)))
	CheckIsType(operations.NewForumCreateConflict(), err)
}
