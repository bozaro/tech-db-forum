package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"math/rand"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type PerfConfig struct {
	UserCount   int
	ForumCount  int
	ThreadCount int
	PostCount   int
	PostBatch   int
	VoteCount   int

	Validate float32
}

func NewPerfConfig() *PerfConfig {
	return &PerfConfig{
		UserCount:   1000,
		ForumCount:  20,
		ThreadCount: 10000,
		PostCount:   1500000,
		PostBatch:   100,
		VoteCount:   100000,
		Validate:    1.0,
	}
}

func FillUsers(perf *Perf, parallel int, timeout time.Time, count int) {
	var need int32 = int32(count)
	c := perf.c
	data := perf.data

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
	waitWaitGroup(&wg, timeout, count, &need)
}

func FillThreads(perf *Perf, parallel int, timeout time.Time, count int) {
	var need int32 = int32(count)
	c := perf.c
	data := perf.data

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
					Voices:      map[*PUser]int32{},
					MessageHash: Hash(thread.Message),
					TitleHash:   Hash(thread.Title),
					Created:     *thread.Created,
				})
			}
			wg.Done()
		}()
	}

	// wait for the workers to finish
	waitWaitGroup(&wg, timeout, count, &need)
}

func FillPosts(perf *Perf, parallel int, timeout time.Time, count int, batchSize int) {
	var need int32 = int32(count)
	c := perf.c
	data := perf.data

	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			f := NewFactory()
			for atomic.AddInt32(&need, -int32(batchSize)) >= 0 {
				offset := float64(int32(count)-atomic.LoadInt32(&need)) / float64(count)
				batch := make([]*models.Post, 0, batchSize)
				thread := data.GetThread(-1, POST_POWER, offset)
				thread.mutex.Lock() // todo: Потом исправить

				parents := data.GetThreadPostsFlat(thread)

				for j := 0; j < batchSize; j++ {
					var parent *PPost
					if (len(parents) > 0) && (rand.Intn(4) == 0) {
						parent = parents[rand.Intn(len(parents))]
					}
					post := f.RandomPost()
					post.Author = data.GetUser(-1).Nickname
					post.Thread = thread.ID
					if parent != nil {
						post.Parent = parent.ID
					}
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
				thread.mutex.Unlock()
			}
			wg.Done()
		}()
	}

	// wait for the workers to finish
	waitWaitGroup(&wg, timeout, count, &need)
}

func VoteThreads(perf *Perf, parallel int, timeout time.Time, count int) {
	var need int32 = int32(count)
	c := perf.c
	data := perf.data

	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			for atomic.AddInt32(&need, -1) >= 0 {
				user := data.GetUser(-1)

				thread := data.GetThread(-1, 1, 0.5)
				thread.mutex.Lock() // todo: Потом исправить

				old_voice := thread.Voices[user]
				var new_voice int32
				if old_voice != 0 {
					new_voice = -old_voice
				} else if rand.Intn(8) < 5 {
					new_voice = 1
				} else {
					new_voice = -1
				}
				thread.Voices[user] = new_voice
				thread.Votes += new_voice - old_voice

				result, err := c.Operations.ThreadVote(operations.NewThreadVoteParams().
					WithSlugOrID(fmt.Sprintf("%d", thread.ID)).
					WithVote(&models.Vote{
						Nickname: user.Nickname,
						Voice:    new_voice,
					}).
					WithContext(Expected(200, nil, nil)))
				CheckNil(err)
				if result.Payload.Votes != thread.Votes {
					panic(fmt.Sprintf("Unexpected votes count: %d != %d", result.Payload.Votes, thread.Votes))
				}
				thread.mutex.Unlock()
			}
			wg.Done()
		}()
	}

	// wait for the workers to finish
	waitWaitGroup(&wg, timeout, count, &need)
}

func waitWaitGroup(wg *sync.WaitGroup, timeout time.Time, total int, need *int32) bool {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	logProgress := func() {
		completed := total - int(atomic.LoadInt32(need))
		if completed > total {
			completed = total
		}
		log.Infof("  %d of %d (%.2f%%)", completed, total, 100*float32(completed)/float32(total))
	}

	for {
		select {
		case <-done:
			logProgress()
			return true
		case <-time.After(time.Second * 10):
			logProgress()
		}
		if time.Now().After(timeout) {
			log.Panic("Timeout")
			return false
		}
	}
}

func NewPerf(url *url.URL, config *PerfConfig) *Perf {
	transport := CreateTransport(url)
	report := Report{
		OnlyError: true,
		Result:    Success,
	}
	c := client.New(&CheckerTransport{transport, &report}, nil)

	data := NewPerfData(config)
	return &Perf{c: c,
		data:     data,
		validate: config.Validate,
	}
}

func (self *Perf) Fill(threads int, timeout_sec int, config *PerfConfig) {
	f := NewFactory()

	timeout := time.Now().Add(time.Second * time.Duration(timeout_sec))

	log.Infof("Clear data")
	_, err := self.c.Operations.Clear(nil)
	CheckNil(err)

	log.Infof("Creating users (%d threads)", threads)
	FillUsers(self, threads, timeout, config.UserCount)

	log.Info("Creating forums")
	for i := 0; i < config.ForumCount; i++ {
		user := self.data.GetUser(-1)
		forum := f.RandomForum()
		forum.User = user.Nickname
		forum = f.CreateForum(self.c, forum, nil)
		self.data.AddForum(&PForum{
			Slug:      forum.Slug,
			TitleHash: Hash(forum.Title),
			User:      user,
		})
	}

	log.Infof("Creating threads (%d threads)", threads)
	FillThreads(self, threads, timeout, config.ThreadCount)

	log.Infof("Vote threads (%d threads)", threads)
	VoteThreads(self, threads, timeout, config.VoteCount)

	log.Infof("Creating posts (%d threads)", threads)
	FillPosts(self, threads, timeout, config.PostCount, config.PostBatch)

	log.Info("Prepare for saving data")
	self.data.Normalize()

	log.Info("Done")
}
