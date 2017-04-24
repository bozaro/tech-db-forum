package tests

import (
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
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "post_get_one_notfound",
		Description: "",
		FnCheck:     Modifications(CheckPostGetOneNotFound),
		Deps: []string{
			"post_get_one_simple",
		},
	})
	Register(Checker{
		Name:        "post_get_one_related",
		Description: "",
		FnCheck:     Modifications(CheckPostGetOneRelated),
		Deps: []string{
			"post_get_one_simple",
		},
	})
	/*PerfRegister(PerfTest{
		Name:   "post_get_one_success",
		Mode:   ModeRead,
		Weight: WeightNormal,
		FnPerf: PerfPostGetOneSuccess,
	})
	PerfRegister(PerfTest{
		Name:   "post_get_one_not_found",
		Mode:   ModeRead,
		Weight: WeightRare,
		FnPerf: PerfPostGetOneNotFound,
	})
	PerfRegister(PerfTest{
		Name:   "post_get_one_related",
		Mode:   ModeRead,
		Weight: WeightNormal,
		FnPerf: PerfPostGetOneRelated,
	})*/
}

func CheckPostGetOneSimple(c *client.Forum) {
	post := CreatePost(c, nil, nil)
	CheckPost(c, post)
}

func CheckPostGetOneRelated(c *client.Forum, m *Modify) {
	user := CreateUser(c, nil)
	forum := CreateForum(c, nil, nil)
	forum.Threads = 1
	forum.Posts = 1
	thread := CreateThread(c, nil, forum, nil)
	temp := RandomPost()
	temp.Author = user.Nickname
	post := CreatePost(c, temp, thread)
	expected := models.PostFull{
		Post: post,
	}

	related := []string{}
	// User
	if m.Bool() {
		related = append(related, "user")
		expected.Author = user
	}
	// Thread
	if m.Bool() {
		related = append(related, "thread")
		expected.Thread = thread
	}
	// Forum
	if m.Bool() {
		related = append(related, "forum")
		expected.Forum = forum
	}

	// Check
	c.Operations.PostGetOne(operations.NewPostGetOneParams().
		WithID(post.ID).
		WithRelated(related).
		WithContext(Expected(200, &expected, filterPostFull)))
}

func filterPostFull(data interface{}) interface{} {
	full := data.(*models.PostFull)
	if full.Thread != nil {
		filterThread(full.Thread)
	}
	return full
}

func CheckPostGetOneNotFound(c *client.Forum, m *Modify) {
	related := []string{}
	// User
	if m.Bool() {
		related = append(related, "user")
	}
	// Thread
	if m.Bool() {
		related = append(related, "thread")
	}
	// Forum
	if m.Bool() {
		related = append(related, "forum")
	}

	// Check
	_, err := c.Operations.PostGetOne(operations.NewPostGetOneParams().
		WithID(POST_FAKE_ID).
		WithRelated(related).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(err, operations.NewPostGetOneNotFound())
}
