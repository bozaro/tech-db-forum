package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"math/rand"
	"strings"
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
	PerfRegister(PerfTest{
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

func (self *PPost) Validate(v PerfValidator, post *models.Post, version PVersion) {
	v.CheckInt64(self.ID, post.ID, "Post.ID")
	v.CheckStr(self.Thread.Forum.Slug, post.Forum, "Post.Forum")
	v.CheckInt(int(self.Thread.ID), int(post.Thread), "Post.Thread")
	v.CheckStr(self.Author.Nickname, post.Author, "Post.Author")
	v.CheckHash(self.MessageHash, post.Message, "Post.Message")
	v.CheckInt64(self.GetParentId(), post.Parent, "Post.Parent")
	v.CheckBool(self.IsEdited, post.IsEdited, "Post.IsEditer")
	v.CheckDate(&self.Created, post.Created, "Created")
	v.Finish(version, self.Version)
}

func PerfPostGetOneSuccess(p *Perf) {
	post := p.data.GetPost(-1)
	postVersion := post.Version

	userVersion := post.Author.Version
	threadVersion := post.Thread.Version
	forumVersion := post.Thread.Forum.Version

	related := GetRandomRelated()
	result, err := p.c.Operations.PostGetOne(operations.NewPostGetOneParams().
		WithID(post.ID).
		WithRelated(related).
		WithContext(Expected(200, nil, nil)))
	CheckNil(err)

	relatedStr := strings.Join(related, ",")
	p.Validate(func(v PerfValidator) {
		payload := result.Payload
		post.Validate(v, payload.Post, postVersion)

		if strings.Contains(relatedStr, "user") {
			CheckIsType(payload.Author, &models.User{})
			post.Author.Validate(v, payload.Author, userVersion)
		}
		if strings.Contains(relatedStr, "forum") {
			CheckIsType(payload.Forum, &models.Forum{})
			post.Thread.Forum.Validate(v, payload.Forum, forumVersion)
		}
		if strings.Contains(relatedStr, "thread") {
			CheckIsType(payload.Thread, &models.Thread{})
			post.Thread.Validate(v, payload.Thread, threadVersion)
		}
	})
}

func PerfPostGetOneNotFound(p *Perf) {
	related := GetRandomRelated()
	var id int64
	for {
		id = rand.Int63n(100000000)
		if p.data.GetPostById(id) == nil {
			break
		}
	}
	_, err := p.c.Operations.PostGetOne(operations.NewPostGetOneParams().
		WithID(id).
		WithRelated(related).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewPostGetOneNotFound(), err)
}

func GetRandomRelated() []string {
	related := []string{}
	r := rand.Int()
	// User
	if r&1 != 0 {
		related = append(related, "user")
	}
	// Thread
	if r&2 != 0 {
		related = append(related, "thread")
	}
	// Forum
	if r&4 != 0 {
		related = append(related, "forum")
	}
	return related
}
