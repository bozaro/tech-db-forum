package tests

import (
	"github.com/go-openapi/strfmt"
	"math/rand"
)

type PerfData struct {
	Status *PStatus
	Users  []*PUser
	Forums []*PForum
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

type PForum struct {
	Version   PVersion
	Posts     int
	Slug      string
	Threads   int
	TitleHash PHash
	User      *PUser
}

func (self *PerfData) GetForumIndex() int {
	return 0
}

func (self *PerfData) GetUserIndex() int {
	return getRandomIndex(len(self.Users))
}

func getRandomIndex(count int) int {
	if count == 0 {
		return -1
	}
	return rand.Intn(count)
}

func (self *PerfData) GetForum(index int) *PForum {
	if index < 0 {
		index = self.GetForumIndex()
	}
	return self.Forums[index]
}

func (self *PerfData) GetUser(index int) *PUser {
	if index < 0 {
		index = self.GetUserIndex()
	}
	return self.Users[index]
}
