package tests

import (
	"github.com/go-openapi/strfmt"
	"math/rand"
	"sync/atomic"
)

type PerfData struct {
	Status  *PStatus
	users   []*PUser
	forums  []*PForum
	threads []*PThread
	posts   []*PPost

	lastIndex int32

	threadsByForum map[string][]*PThread
	usersByForum   map[string]map[*PUser]bool
	postsByThread  map[int32][]*PPost
	userByNickname map[string]*PUser
	postById       map[int64]*PPost
	threadById     map[int32]*PThread
}

type PStatus struct {
	Version PVersion
	Forum   int
	Post    int
	Thread  int
	User    int
}

type PUser struct {
	Version      PVersion
	AboutHash    PHash
	Email        strfmt.Email
	FullnameHash PHash
	Nickname     string
}

type PThread struct {
	Version     PVersion
	ID          int32
	Slug        string
	Author      *PUser
	Forum       *PForum
	MessageHash PHash
	TitleHash   PHash
	Created     strfmt.DateTime
	Votes       int32
	Posts       int
}

type PForum struct {
	Version   PVersion
	Posts     int
	Slug      string
	Threads   int
	TitleHash PHash
	User      *PUser
}

type PPost struct {
	Version     PVersion
	ID          int64
	Author      *PUser
	Thread      *PThread
	Parent      *PPost
	Created     strfmt.DateTime
	IsEdited    bool
	MessageHash PHash
	Index       int32
	Path        []int32
}

func NewPerfData() *PerfData {
	return &PerfData{
		Status:         &PStatus{},
		forums:         []*PForum{},
		users:          []*PUser{},
		threads:        []*PThread{},
		posts:          []*PPost{},
		threadsByForum: map[string][]*PThread{},
		usersByForum:   map[string]map[*PUser]bool{},
		postsByThread:  map[int32][]*PPost{},
		userByNickname: map[string]*PUser{},
		threadById:     map[int32]*PThread{},
		postById:       map[int64]*PPost{},
	}
}

func (self *PerfData) GetUserByNickname(nickname string) *PUser {
	return self.userByNickname[nickname]
}

func getRandomIndex(count int) int {
	if count == 0 {
		return -1
	}
	return rand.Intn(count)
}

func (self *PerfData) AddForum(forum *PForum) {
	self.forums = append(self.forums, forum)
	self.usersByForum[forum.Slug] = map[*PUser]bool{}
	self.threadsByForum[forum.Slug] = []*PThread{}
	self.Status.Forum++
}

func (self *PerfData) GetForum(index int) *PForum {
	if index < 0 {
		index = getRandomIndex(len(self.forums))
	}
	return self.forums[index]
}

func (self *PerfData) AddUser(user *PUser) {
	self.users = append(self.users, user)
	self.userByNickname[user.Nickname] = user
	self.Status.User++
}

func (self *PerfData) GetUser(index int) *PUser {
	if index < 0 {
		index = getRandomIndex(len(self.users))
	}
	return self.users[index]
}

func (self *PerfData) AddThread(thread *PThread) {
	self.threads = append(self.threads, thread)
	self.threadById[thread.ID] = thread
	self.postsByThread[thread.ID] = []*PPost{}
	self.threadsByForum[thread.Forum.Slug] = append(self.threadsByForum[thread.Forum.Slug], thread)
	self.usersByForum[thread.Forum.Slug][thread.Author] = true
	thread.Forum.Threads++
	self.Status.Thread++
}

func (self *PerfData) GetThread(index int) *PThread {
	if index < 0 {
		index = getRandomIndex(len(self.threads))
	}
	return self.threads[index]
}

func (self *PerfData) GetThreadById(id int32) *PThread {
	return self.threadById[id]
}

func (self *PerfData) GetPostById(id int64) *PPost {
	return self.postById[id]
}

func (self *PerfData) GetForumThreads(forum *PForum) []*PThread {
	result := self.threadsByForum[forum.Slug]
	if result == nil {
		return []*PThread{}
	}
	return result
}

func (self *PerfData) GetForumUsers(forum *PForum) []*PUser {
	users := self.usersByForum[forum.Slug]
	if users == nil {
		return []*PUser{}
	}
	result := make([]*PUser, 0, len(users))
	for k := range users {
		result = append(result, k)
	}
	return result
}

func (self *PerfData) GetThreadPosts(thread *PThread) []*PPost {
	result := self.postsByThread[thread.ID]
	if result == nil {
		return []*PPost{}
	}
	return result
}

func (self *PerfData) AddPost(post *PPost) {
	self.posts = append(self.posts, post)
	self.postById[post.ID] = post
	self.usersByForum[post.Thread.Forum.Slug][post.Author] = true
	self.postsByThread[post.Thread.ID] = append(self.postsByThread[post.Thread.ID], post)

	post.Index = atomic.AddInt32(&self.lastIndex, 1)
	if post.Parent != nil {
		post.Path = append(post.Parent.Path, post.Index)
	} else {
		post.Path = []int32{post.Index}
	}
	post.Thread.Forum.Posts++
	post.Thread.Posts++
	self.Status.Post++
}

func (self *PerfData) GetPost(index int) *PPost {
	if index < 0 {
		index = getRandomIndex(len(self.posts))
	}
	return self.posts[index]
}

func (self *PPost) GetParentId() int64 {
	if self.Parent == nil {
		return 0
	}
	return self.Parent.ID
}

func GetRandomLimit() int32 {
	return 15 + rand.Int31n(5)
}

func GetRandomDesc() *bool {
	switch rand.Intn(3) {
	case 0:
		v := false
		return &v
	case 1:
		v := true
		return &v
	default:
		return nil
	}
}
