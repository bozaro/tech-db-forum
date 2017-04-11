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
	data.Users = make([]*PUser, count)
	for i := 0; i < count; i++ {
		data.Users[i] = <-results
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

	data := &PerfData{}

	log.Info("Creating users (multiple threads)")
	FillUsers(c, data, 8, 1000)

	log.Info("Creating forums")
	for i := 0; i < 20; i++ {
		user := data.GetUser(-1)
		forum := RandomForum()
		forum.User = user.Nickname
		forum = CreateForum(c, forum, nil)
		data.Forums = append(data.Forums, &PForum{
			Slug:      forum.Slug,
			TitleHash: Hash(forum.Title),
			User:      user,
		})
	}

	log.Info("Creating threads")
	threads := []*models.Thread{}
	for i := 0; i < 1000; i++ {
		author := data.GetUser(-1)
		forum := data.GetForum(-1)
		thread := RandomThread()
		if rand.Intn(100) >= 25 {
			thread.Slug = ""
		}
		thread.Author = author.Nickname
		thread.Forum = forum.Slug
		threads = append(threads, CreateThread(c, thread, nil, nil))
		forum.Threads++
	}

	log.Info("Creating posts")
	posts := []*models.Post{}
	for i := 0; i < 100; i++ {
		batch := []*models.Post{}
		thread := threads[rand.Intn(len(threads))].ID
		for j := 0; j < 100; j++ {
			post := RandomPost()
			post.Author = data.GetUser(-1).Nickname
			post.Thread = thread
			batch = append(batch, post)
		}
		posts = append(posts, CreatePosts(c, batch, nil)...)
	}

	data.Status = &PStatus{
		User:   len(data.Users),
		Forum:  len(data.Forums),
		Thread: len(threads),
		Post:   len(posts),
	}

	log.Info("Done")
	return &Perf{c: c,
		data: data,
	}
}
