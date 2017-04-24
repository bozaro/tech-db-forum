package tests

import (
	"github.com/go-openapi/strfmt"
	"math/rand"
)

type PerfData struct {
	Status  *PStatus
	Users   []*PUser
	Forums  []*PForum
	Threads []*PThread
	Posts   []*PPost

	userByNickname map[string]*PUser
	postById       map[int64]*PPost
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
}

func (self *PerfData) GetUserByNickname(nickname string) *PUser {
	return self.userByNickname[nickname]
}

func (self *PerfData) GetForumIndex() int {
	return 0
}

func (self *PerfData) GetUserIndex() int {
	return getRandomIndex(len(self.Users))
}

func (self *PerfData) GetThreadIndex() int {
	return getRandomIndex(len(self.Threads))
}

func (self *PerfData) GetPostIndex() int {
	return getRandomIndex(len(self.Posts))
}

func getRandomIndex(count int) int {
	if count == 0 {
		return -1
	}
	return rand.Intn(count)
}

func (self *PerfData) AddForum(forum *PForum) {
	self.Forums = append(self.Forums, forum)
	self.Status.Forum++
}

func (self *PerfData) GetForum(index int) *PForum {
	if index < 0 {
		index = self.GetForumIndex()
	}
	return self.Forums[index]
}

func (self *PerfData) AddUser(user *PUser) {
	self.Users = append(self.Users, user)
	self.userByNickname[user.Nickname] = user
	self.Status.User++
}

func (self *PerfData) GetUser(index int) *PUser {
	if index < 0 {
		index = self.GetUserIndex()
	}
	return self.Users[index]
}

func (self *PerfData) AddThread(thread *PThread) {
	self.Threads = append(self.Threads, thread)
	thread.Forum.Threads++
	self.Status.Thread++
}

func (self *PerfData) GetThread(index int) *PThread {
	if index < 0 {
		index = self.GetThreadIndex()
	}
	return self.Threads[index]
}

func (self *PerfData) GetPostById(id int64) *PPost {
	return self.postById[id]
}

func (self *PerfData) AddPost(post *PPost) {
	self.Posts = append(self.Posts, post)
	self.postById[post.ID] = post
	post.Thread.Forum.Posts++
	post.Thread.Posts++
	self.Status.Post++
}

func (self *PerfData) GetPost(index int) *PPost {
	if index < 0 {
		index = self.GetPostIndex()
	}
	return self.Posts[index]
}
