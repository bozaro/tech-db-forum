package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
)

func init() {
	Register(Checker{
		Name:        "post_get_one_simple",
		Description: "",
		FnCheck:     CheckPostGetOneSimple,
		Deps: []string{
			"post_create_simple",
		},
	})
	Register(Checker{
		Name:        "post_get_one_notfound",
		Description: "",
		FnCheck:     CheckPostGetOneNotFound,
		Deps: []string{
			"post_get_one_simple",
		},
	})
	Register(Checker{
		Name:        "post_get_one_related",
		Description: "",
		FnCheck:     CheckPostGetOneRelated,
		Deps: []string{
			"post_get_one_simple",
		},
	})
}

func CheckPostGetOneSimple(c *client.Forum) {
	post := CreatePost(c, nil, nil)
	CheckPost(c, post)
}

func CheckPostGetOneRelated(c *client.Forum) {
	pass := 0
	for true {
		pass++
		Checkpoint(c, fmt.Sprintf("Pass %d", pass))

		user := CreateUser(c, nil)
		forum := CreateForum(c, nil, nil)
		forum.Threads = 1
		thread := CreateThread(c, nil, forum, nil)
		temp := RandomPost()
		temp.Author = user.Nickname
		post := CreatePost(c, temp, thread)
		expected := models.PostFull{
			Post: post,
		}

		modify := pass
		related := []string{}
		// User
		if (modify & 1) == 1 {
			related = append(related, "user")
			expected.Author = user
		}
		modify >>= 1
		// Thread
		if (modify & 1) == 1 {
			related = append(related, "thread")
			expected.Thread = thread
		}
		modify >>= 1
		// Forum
		if (modify & 1) == 1 {
			related = append(related, "forum")
			expected.Forum = forum
		}
		modify >>= 1
		// Done?
		if modify != 0 {
			break
		}
		// Check
		c.Operations.PostGetOne(operations.NewPostGetOneParams().
			WithID(post.ID).
			WithRelated(related).
			WithContext(Expected(200, &expected, nil)))
	}
}

func CheckPostGetOneNotFound(c *client.Forum) {
	pass := 0
	for true {
		Checkpoint(c, fmt.Sprintf("Pass %d", pass))

		modify := pass
		related := []string{}
		// User
		if (modify & 1) == 1 {
			related = append(related, "user")
		}
		modify >>= 1
		// Thread
		if (modify & 1) == 1 {
			related = append(related, "thread")
		}
		modify >>= 1
		// Forum
		if (modify & 1) == 1 {
			related = append(related, "forum")
		}
		modify >>= 1
		// Done?
		if modify != 0 {
			break
		}
		pass++
		// Check
		_, err := c.Operations.PostGetOne(operations.NewPostGetOneParams().
			WithID(POST_FAKE_ID).
			WithRelated(related).
			WithContext(Expected(404, nil, nil)))
		CheckIsType(err, operations.NewPostGetOneNotFound())

	}

}
