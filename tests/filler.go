package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/models"
	"log"
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

	log.Println("Creating users")
	/*users := []*models.User{}
	for i := 0; i < 1000; i++ {
		users = append(users, CreateUser(c, nil))
	}*/
	users := FillUsers(c, 8, 1000)

	log.Println("Creating forums")
	forums := []*models.Forum{}
	for i := 0; i < 20; i++ {
		forums = append(forums, CreateForum(c, nil, users[rand.Intn(len(users))]))
	}

	log.Println("Creating threads")
	threads := []*models.Thread{}
	for i := 0; i < 1000; i++ {
		thread := RandomThread()
		if rand.Intn(100) >= 5 {
			thread.Slug = ""
		}
		threads = append(threads, CreateThread(c, thread, forums[rand.Intn(len(forums))], users[rand.Intn(len(users))]))
	}

	log.Println("Creating posts")
	posts := []*models.Post{}
	for i := 0; i < 10000; i++ {
		post := RandomPost()
		post.Author = users[rand.Intn(len(users))].Nickname
		post.Thread = threads[rand.Intn(len(users))].ID
		posts = append(posts, CreatePost(c, post, nil))
	}

	log.Println("Done")
	return 0
}
