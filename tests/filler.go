package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/models"
	"math/rand"
	"net/url"
	"sync"
	"sync/atomic"
)

type PerfConfig struct {
	UserCount   int
	ForumCount  int
	ThreadCount int
	PostCount   int
	PostBatch   int
}

func NewPerfConfig() *PerfConfig {
	return &PerfConfig{
		UserCount:   1000,
		ForumCount:  20,
		ThreadCount: 1000,
		PostCount:   1000000,
		PostBatch:   100,
	}
}

func FillUsers(c *client.Forum, data *PerfData, parallel int, count int) {
	var need int32 = int32(count)

	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			f := NewFactory()
			for atomic.AddInt32(&need, -1) >= 0 {
				user := f.CreateUser(c, nil)
				data.AddUser(&PUser{
					AboutHash:    Hash(user.About),
					Email:        user.Email,
					FullnameHash: Hash(user.Fullname),
					Nickname:     user.Nickname,
				})
			}
			wg.Done()
		}()
	}

	// wait for the workers to finish
	wg.Wait()
}

func FillThreads(c *client.Forum, data *PerfData, parallel int, count int) {
	var need int32 = int32(count)

	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			f := NewFactory()
			for atomic.AddInt32(&need, -1) >= 0 {
				author := data.GetUser(-1)
				forum := data.GetForum(-1)
				thread := f.RandomThread()
				if rand.Intn(100) >= 25 {
					thread.Slug = ""
				}
				thread.Author = author.Nickname
				thread.Forum = forum.Slug
				thread = f.CreateThread(c, thread, nil, nil)
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
			wg.Done()
		}()
	}

	// wait for the workers to finish
	wg.Wait()
}

func FillPosts(c *client.Forum, data *PerfData, parallel int, count int, batchSize int) {
	var need int32 = int32(count)

	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			f := NewFactory()
			for atomic.AddInt32(&need, -int32(batchSize)) >= 0 {

				batch := make([]*models.Post, 0, batchSize)
				thread := data.GetThread(-1)
				for j := 0; j < batchSize; j++ {
					post := f.RandomPost()
					post.Author = data.GetUser(-1).Nickname
					post.Thread = thread.ID
					batch = append(batch, post)
				}
				for _, post := range f.CreatePosts(c, batch, nil) {
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
			wg.Done()
		}()
	}

	// wait for the workers to finish
	wg.Wait()
}

func Fill(url *url.URL, threads int, config *PerfConfig) *Perf {

	transport := CreateTransport(url)
	c := client.New(transport, nil)
	_, err := c.Operations.Clear(nil)
	CheckNil(err)

	data := NewPerfData(config)
	f := NewFactory()

	log.Infof("Creating users (%d threads)", threads)
	FillUsers(c, data, threads, config.UserCount)

	log.Info("Creating forums")
	for i := 0; i < config.ForumCount; i++ {
		user := data.GetUser(-1)
		forum := f.RandomForum()
		forum.User = user.Nickname
		forum = f.CreateForum(c, forum, nil)
		data.AddForum(&PForum{
			Slug:      forum.Slug,
			TitleHash: Hash(forum.Title),
			User:      user,
		})
	}

	log.Infof("Creating threads (%d threads)", threads)
	FillThreads(c, data, threads, config.ThreadCount)

	log.Infof("Creating posts (%d threads)", threads)
	FillPosts(c, data, threads, config.PostCount, config.PostBatch)

	log.Info("Done")
	return &Perf{c: c,
		data: data,
	}
}
