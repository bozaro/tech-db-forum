package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"math/rand"
	"strings"
)

const (
	RELATED_FORUM  = "forum"
	RELATED_USER   = "user"
	RELATED_THREAD = "thread"
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

func CheckPostGetOneSimple(c *client.Forum, f *Factory) {
	post := f.CreatePost(c, nil, nil)
	CheckPost(c, post)
}

func CheckPostGetOneRelated(c *client.Forum, f *Factory, m *Modify) {
	user := f.CreateUser(c, nil)
	forum := f.CreateForum(c, nil, nil)
	forum.Threads = 1
	forum.Posts = 1
	thread := f.CreateThread(c, nil, forum, nil)
	temp := f.RandomPost()
	temp.Author = user.Nickname
	post := f.CreatePost(c, temp, thread)
	expected := models.PostFull{
		Post: post,
	}

	related := []string{}
	// User
	if m.Bool() {
		related = append(related, RELATED_USER)
		expected.Author = user
	}
	// Thread
	if m.Bool() {
		related = append(related, RELATED_THREAD)
		expected.Thread = thread
	}
	// Forum
	if m.Bool() {
		related = append(related, RELATED_FORUM)
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

func CheckPostGetOneNotFound(c *client.Forum, f *Factory, m *Modify) {
	related := []string{}
	// User
	if m.Bool() {
		related = append(related, RELATED_USER)
	}
	// Thread
	if m.Bool() {
		related = append(related, RELATED_THREAD)
	}
	// Forum
	if m.Bool() {
		related = append(related, RELATED_FORUM)
	}

	// Check
	_, err := c.Operations.PostGetOne(operations.NewPostGetOneParams().
		WithID(POST_FAKE_ID).
		WithRelated(related).
		WithContext(ExpectedError(404, "Can't find post with id: %d", POST_FAKE_ID)))
	CheckIsType(err, operations.NewPostGetOneNotFound())
}

func (self *PPost) Validate(v PerfValidator, post *models.Post, version PVersion, prefix string) {
	v.CheckInt64(self.ID, post.ID, prefix+".ID")
	v.CheckStr(self.Thread.Forum.Slug, post.Forum, prefix+".Forum")
	v.CheckInt(int(self.Thread.ID), int(post.Thread), prefix+".Thread")
	v.CheckStr(self.Author.Nickname, post.Author, prefix+".Author")
	v.CheckHash(self.MessageHash, post.Message, prefix+".Message")
	v.CheckInt64(self.GetParentId(), post.Parent, prefix+".Parent")
	v.CheckBool(self.IsEdited, post.IsEdited, prefix+".IsEditer")
	v.CheckDate(&self.Created, post.Created, prefix+".Created")
	v.Finish(version, self.Version)
}

func PerfPostGetOneSuccess(p *Perf, f *Factory) {
	post := p.data.GetPost(-1, POST_POWER)
	postVersion := post.Version

	userVersion := post.Author.Version
	threadVersion := post.Thread.Version
	forumVersion := post.Thread.Forum.Version

	related := GetRandomRelated()
	s := p.Session()
	result, err := p.c.Operations.PostGetOne(operations.NewPostGetOneParams().
		WithID(post.ID).
		WithRelated(related).
		WithContext(s.Expected(200)))
	CheckNil(err)

	relatedStr := strings.Join(related, ",")
	s.Validate(func(v PerfValidator) {
		payload := result.Payload
		post.Validate(v, payload.Post, postVersion, "Post")

		if strings.Contains(relatedStr, RELATED_USER) {
			CheckIsType(payload.Author, &models.User{})
			post.Author.Validate(v, payload.Author, userVersion)
		}
		if strings.Contains(relatedStr, RELATED_FORUM) {
			CheckIsType(payload.Forum, &models.Forum{})
			post.Thread.Forum.Validate(v, payload.Forum, forumVersion)
		}
		if strings.Contains(relatedStr, RELATED_THREAD) {
			CheckIsType(payload.Thread, &models.Thread{})
			post.Thread.Validate(v, payload.Thread, threadVersion)
		}
	})
}

func PerfPostGetOneNotFound(p *Perf, f *Factory) {
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
		WithContext(ExpectedError(404, "Can't find post with id: %d", id)))
	CheckIsType(operations.NewPostGetOneNotFound(), err)
}

func GetRandomRelated() []string {
	related := []string{}
	r := rand.Int()
	// User
	if r&1 != 0 {
		related = append(related, RELATED_USER)
	}
	// Thread
	if r&2 != 0 {
		related = append(related, RELATED_THREAD)
	}
	// Forum
	if r&4 != 0 {
		related = append(related, RELATED_FORUM)
	}
	return related
}
