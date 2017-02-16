package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"strings"
)

func init() {
	Register(Checker{
		Name:        "post_update_simple",
		Description: "",
		FnCheck:     CheckPostUpdateSimple,
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "post_update_empty",
		Description: "",
		FnCheck:     CheckPostUpdateEmpty,
		Deps: []string{
			"post_update_simple",
		},
	})
	Register(Checker{
		Name:        "post_update_notfound",
		Description: "",
		FnCheck:     CheckPostUpdateNotFound,
		Deps: []string{
			"post_update_simple",
		},
	})
	Register(Checker{
		Name:        "post_update_same",
		Description: "",
		FnCheck:     CheckPostUpdateSame,
		Deps: []string{
			"post_update_simple",
		},
	})
	Register(Checker{
		Name:        "post_update_case",
		Description: "",
		FnCheck:     CheckPostUpdateCase,
		Deps: []string{
			"post_update_simple",
		},
	})
}

func CheckPostUpdateSimple(c *client.Forum) {
	post := CreatePost(c, nil, nil)
	temp := RandomPost()

	update := &models.PostUpdate{}
	update.Message = temp.Message

	edited := true
	expected := *post
	expected.IsEdited = &edited
	expected.Message = update.Message

	c.Operations.PostUpdate(operations.NewPostUpdateParams().
		WithID(post.ID).
		WithPost(update).
		WithContext(Expected(200, &expected, nil)))

	CheckPost(c, &expected)
}

func CheckPostUpdateEmpty(c *client.Forum) {
	post := CreatePost(c, nil, nil)
	c.Operations.PostUpdate(operations.NewPostUpdateParams().
		WithID(post.ID).
		WithPost(&models.PostUpdate{}).
		WithContext(Expected(200, post, nil)))

	CheckPost(c, post)
}

func CheckPostUpdateSame(c *client.Forum) {
	post := CreatePost(c, nil, nil)
	c.Operations.PostUpdate(operations.NewPostUpdateParams().
		WithID(post.ID).
		WithPost(&models.PostUpdate{
			Message: post.Message,
		}).
		WithContext(Expected(200, post, nil)))

	CheckPost(c, post)
}

func CheckPostUpdateCase(c *client.Forum) {
	edited := true
	post := CreatePost(c, nil, nil)
	post.Message = strings.ToUpper(post.Message)
	post.IsEdited = &edited
	c.Operations.PostUpdate(operations.NewPostUpdateParams().
		WithID(post.ID).
		WithPost(&models.PostUpdate{
			Message: post.Message,
		}).
		WithContext(Expected(200, post, nil)))

	CheckPost(c, post)
}

func CheckPostUpdateNotFound(c *client.Forum) {
	post := RandomPost()
	_, err := c.Operations.PostUpdate(operations.NewPostUpdateParams().
		WithID(post.ID).
		WithPost(&models.PostUpdate{
			Message: post.Message,
		}).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewPostUpdateNotFound(), err)
}
