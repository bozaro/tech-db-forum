package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/models"
	"log"
	"math/rand"
	"net/url"
)

func Fill(url *url.URL) int {

	transport := CreateTransport(url)
	c := client.New(transport, nil)
	_, err := c.Operations.Clear(nil)
	CheckNil(err)

	log.Println("Creating users")
	users := []*models.User{}
	for i := 0; i < 1000; i++ {
		users = append(users, CreateUser(c, nil))
	}

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
