package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/models"
	"math/rand"
	"net/url"
	"sync"
	"sync/atomic"
)

func FillUsers(c *client.Forum, data *PerfData, parallel int, count int) {
	results := make(chan *PUser, 64)
	var need int32 = int32(count)

	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			for atomic.AddInt32(&need, -1) >= 0 {
				user := CreateUser(c, nil)
				results <- &PUser{
					AboutHash:    Hash(user.About),
					Email:        user.Email,
					FullnameHash: Hash(user.Fullname),
					Nickname:     user.Nickname,
				}
			}
			wg.Done()
		}()
	}

	// get result
	for i := 0; i < count; i++ {
		data.AddUser(<-results)
	}
	close(results)

	// wait for the workers to finish
	wg.Wait()
}

func Fill(url *url.URL) *Perf {

	transport := CreateTransport(url)
	c := client.New(transport, nil)
	_, err := c.Operations.Clear(nil)
	CheckNil(err)

	data := NewPerfData()

	log.Info("Creating users (multiple threads)")
	FillUsers(c, data, 8, 1000)

	log.Info("Creating forums")
	for i := 0; i < 20; i++ {
		user := data.GetUser(-1)
		forum := RandomForum()
		forum.User = user.Nickname
		forum = CreateForum(c, forum, nil)
		data.AddForum(&PForum{
			Slug:      forum.Slug,
			TitleHash: Hash(forum.Title),
			User:      user,
		})
	}

	log.Info("Creating threads")
	for i := 0; i < 1000; i++ {
		author := data.GetUser(-1)
		forum := data.GetForum(-1)
		thread := RandomThread()
		if rand.Intn(100) >= 25 {
			thread.Slug = ""
		}
		thread.Author = author.Nickname
		thread.Forum = forum.Slug
		thread = CreateThread(c, thread, nil, nil)
		data.AddThread(&PThread{
			ID:          thread.ID,
			Slug:        thread.Slug,
			Author:      author,
			Forum:       forum,
			MessageHash: Hash(thread.Message),
			TitleHash:   Hash(thread.Title),
			Created:     *thread.Created,
		})
	}

	log.Info("Creating posts")
	for i := 0; i < 10000; i++ {
		batch := []*models.Post{}
		thread := data.GetThread(-1)
		for j := 0; j < 100; j++ {
			post := RandomPost()
			post.Author = data.GetUser(-1).Nickname
			post.Thread = thread.ID
			batch = append(batch, post)
		}
		for _, post := range CreatePosts(c, batch, nil) {
			data.AddPost(&PPost{
				ID:          post.ID,
				Author:      data.GetUserByNickname(post.Author),
				Thread:      thread,
				Parent:      data.GetPostById(post.Parent),
				Created:     *post.Created,
				IsEdited:    false,
				MessageHash: Hash(post.Message),
			})
		}
	}

	log.Info("Done")
	return &Perf{c: c,
		data: data,
	}
}
