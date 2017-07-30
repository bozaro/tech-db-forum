package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
)

func init() {
	Register(Checker{
		Name:        "thread_create_vote_simple",
		Description: "",
		FnCheck:     CheckThreadVoteSimple,
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "thread_create_vote_notfound_thread",
		Description: "",
		FnCheck:     CheckThreadVoteThreadNotFound,
		Deps: []string{
			"thread_create_vote_simple",
		},
	})
	Register(Checker{
		Name:        "thread_create_vote_notfound_user",
		Description: "",
		FnCheck:     CheckThreadVoteUserNotFound,
		Deps: []string{
			"thread_create_vote_simple",
		},
	})
}

func CheckThreadVoteThreadNotFound(c *client.Forum, f *Factory) {
	user := f.CreateUser(c, nil)
	thread := f.RandomThread()
	var err error
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(thread.Slug).
		WithVote(&models.Vote{
			Nickname: user.Nickname,
			Voice:    1,
		}).
		WithContext(ExpectedError(404, "Can't find thread by slug: %s", thread.Slug)))
	CheckIsType(operations.NewThreadVoteNotFound(), err)

	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(THREAD_FAKE_ID).
		WithVote(&models.Vote{
			Nickname: user.Nickname,
			Voice:    1,
		}).
		WithContext(ExpectedError(404, "Can't find thread by id: %d", THREAD_FAKE_ID)))
	CheckIsType(operations.NewThreadVoteNotFound(), err)
}

func CheckThreadVoteUserNotFound(c *client.Forum, f *Factory) {
	user := f.RandomUser()
	thread := f.CreateThread(c, nil, nil, nil)
	_, err := c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(thread.Slug).
		WithVote(&models.Vote{
			Nickname: user.Nickname,
			Voice:    1,
		}).
		WithContext(ExpectedError(404, "Can't find user by nickname: %s", user.Nickname)))
	CheckIsType(operations.NewThreadVoteNotFound(), err)
}

func CheckThreadVoteSimple(c *client.Forum, f *Factory) {
	user1 := f.CreateUser(c, nil)
	user2 := f.CreateUser(c, nil)
	user3 := f.CreateUser(c, nil)
	thread := f.CreateThread(c, nil, nil, nil)
	var err error
	// Like user1
	thread.Votes = 1
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(thread.Slug).
		WithVote(&models.Vote{
			Nickname: user1.Nickname,
			Voice:    1,
		}).
		WithContext(Expected(200, thread, filterThread)))
	CheckNil(err)
	// Like user2
	thread.Votes = 2
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(fmt.Sprintf("%d", thread.ID)).
		WithVote(&models.Vote{
			Nickname: user2.Nickname,
			Voice:    1,
		}).
		WithContext(Expected(200, thread, filterThread)))
	// Like user3
	thread.Votes = 3
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(fmt.Sprintf("%d", thread.ID)).
		WithVote(&models.Vote{
			Nickname: user3.Nickname,
			Voice:    1,
		}).
		WithContext(Expected(200, thread, filterThread)))
	// Dislike user2
	thread.Votes = 1
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(fmt.Sprintf("%d", thread.ID)).
		WithVote(&models.Vote{
			Nickname: user2.Nickname,
			Voice:    -1,
		}).
		WithContext(Expected(200, thread, filterThread)))
	// Dislike user2 again
	thread.Votes = 1
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(thread.Slug).
		WithVote(&models.Vote{
			Nickname: user2.Nickname,
			Voice:    -1,
		}).
		WithContext(Expected(200, thread, filterThread)))
	CheckNil(err)
}
