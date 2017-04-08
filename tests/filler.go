package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/models"
	"math/rand"
	"net/url"
	"sync"
	"sync/atomic"
)

func FillUsers(c *client.Forum, parallel int, count int) []*models.User {
	results := make(chan *models.User, 64)
	var need int32 = int32(count)

	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			for atomic.AddInt32(&need, -1) >= 0 {
				results <- CreateUser(c, nil)
			}
			wg.Done()
		}()
	}

	// get result
	result := make([]*models.User, count)
	for i := 0; i < count; i++ {
		result[i] = <-results
	}
	close(results)

	// wait for the workers to finish
	wg.Wait()
	return result
}

func Fill(url *url.URL) int {

	transport := CreateTransport(url)
	c := client.New(transport, nil)
	_, err := c.Operations.Clear(nil)
	CheckNil(err)

	log.Info("Creating users (multiple threads)")
	users := FillUsers(c, 8, 10000)

	log.Info("Creating forums")
	forums := []*models.Forum{}
	for i := 0; i < 20; i++ {
		forums = append(forums, CreateForum(c, nil, users[rand.Intn(len(users))]))
	}

	log.Info("Creating threads")
	threads := []*models.Thread{}
	for i := 0; i < 1000; i++ {
		thread := RandomThread()
		if rand.Intn(100) >= 5 {
			thread.Slug = ""
		}
		threads = append(threads, CreateThread(c, thread, forums[rand.Intn(len(forums))], users[rand.Intn(len(users))]))
	}

	log.Info("Creating posts")
	posts := []*models.Post{}
	for i := 0; i < 10000; i++ {
		batch := []*models.Post{}
		thread := threads[rand.Intn(len(threads))].ID
		for j := 0; j < 100; j++ {
			post := RandomPost()
			post.Author = users[rand.Intn(len(users))].Nickname
			post.Thread = thread
			batch = append(batch, post)
		}
		posts = append(posts, CreatePosts(c, batch, nil)...)
	}

	log.Info("Done")
	return 0
}
