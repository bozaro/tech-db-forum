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

func CheckThreadVoteThreadNotFound(c *client.Forum) {
	user := CreateUser(c, nil)
	thread := RandomThread()
	var err error
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(thread.Slug).
		WithVote(&models.Vote{
			Nickname: user.Nickname,
			Voice:    1,
		}).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewThreadVoteNotFound(), err)

	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(THREAD_FAKE_ID).
		WithVote(&models.Vote{
			Nickname: user.Nickname,
			Voice:    1,
		}).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewThreadVoteNotFound(), err)
}

func CheckThreadVoteUserNotFound(c *client.Forum) {
	user := RandomUser()
	thread := CreateThread(c, nil, nil, nil)
	_, err := c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(thread.Slug).
		WithVote(&models.Vote{
			Nickname: user.Nickname,
			Voice:    1,
		}).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewThreadVoteNotFound(), err)
}

func CheckThreadVoteSimple(c *client.Forum) {
	user1 := CreateUser(c, nil)
	user2 := CreateUser(c, nil)
	user3 := CreateUser(c, nil)
	thread := CreateThread(c, nil, nil, nil)
	var err error
	// Like user1
	thread.Votes = 1
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(thread.Slug).
		WithVote(&models.Vote{
			Nickname: user1.Nickname,
			Voice:    1,
		}).
		WithContext(Expected(200, &thread, nil)))
	CheckNil(err)
	// Like user2
	thread.Votes = 2
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(fmt.Sprintf("%d", thread.ID)).
		WithVote(&models.Vote{
			Nickname: user2.Nickname,
			Voice:    1,
		}).
		WithContext(Expected(200, &thread, nil)))
	// Like user3
	thread.Votes = 3
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(fmt.Sprintf("%d", thread.ID)).
		WithVote(&models.Vote{
			Nickname: user3.Nickname,
			Voice:    1,
		}).
		WithContext(Expected(200, &thread, nil)))
	// Dislike user2
	thread.Votes = 1
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(fmt.Sprintf("%d", thread.ID)).
		WithVote(&models.Vote{
			Nickname: user2.Nickname,
			Voice:    -1,
		}).
		WithContext(Expected(200, &thread, nil)))
	// Dislike user2 again
	thread.Votes = 1
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(thread.Slug).
		WithVote(&models.Vote{
			Nickname: user2.Nickname,
			Voice:    -1,
		}).
		WithContext(Expected(200, &thread, nil)))
	CheckNil(err)
}
