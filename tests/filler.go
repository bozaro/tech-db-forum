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
	results := make(chan *PerfUser, 64)
	var need int32 = int32(count)

	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			for atomic.AddInt32(&need, -1) >= 0 {
				user := CreateUser(c, nil)
				results <- &PerfUser{
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
	data.Users = make([]*PerfUser, count)
	for i := 0; i < count; i++ {
		data.Users[i] = <-results
	}
	close(results)

	// wait for the workers to finish
	wg.Wait()
}

func Fill(url *url.URL) *PerfData {

	transport := CreateTransport(url)
	c := client.New(transport, nil)
	_, err := c.Operations.Clear(nil)
	CheckNil(err)

	data := &PerfData{}

	log.Info("Creating users (multiple threads)")
	FillUsers(c, data, 8, 10000)

	log.Info("Creating forums")
	forums := []*models.Forum{}
	for i := 0; i < 20; i++ {
		forum := RandomForum()
		forum.User = data.GetUser().Nickname
		forums = append(forums, CreateForum(c, forum, nil))
	}

	log.Info("Creating threads")
	threads := []*models.Thread{}
	for i := 0; i < 1000; i++ {
		thread := RandomThread()
		if rand.Intn(100) >= 5 {
			thread.Slug = ""
		}
		thread.Author = data.GetUser().Nickname
		threads = append(threads, CreateThread(c, thread, forums[rand.Intn(len(forums))], nil))
	}

	log.Info("Creating posts")
	posts := []*models.Post{}
	for i := 0; i < 100; i++ {
		batch := []*models.Post{}
		thread := threads[rand.Intn(len(threads))].ID
		for j := 0; j < 100; j++ {
			post := RandomPost()
			post.Author = data.GetUser().Nickname
			post.Thread = thread
			batch = append(batch, post)
		}
		posts = append(posts, CreatePosts(c, batch, nil)...)
	}

	log.Info("Done")
	return data
}
