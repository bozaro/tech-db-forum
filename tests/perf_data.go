package tests

import (
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/go-openapi/strfmt"
	"math/rand"
)

type PerfData struct {
	Users []*PerfUser
}

type PerfUser struct {
	AboutHash    PerfHash
	Email        strfmt.Email
	FullnameHash PerfHash
	Nickname     string
}

func (self *PerfData) Status(with_pending bool) models.Status {
	return models.Status{
		Forum:  0,
		Post:   0,
		User:   0,
		Thread: 0,
	}
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

func (self *PerfData) GetForumData(index int, with_pending bool) models.Forum {
	return models.Forum{
		Posts:   0,
		Slug:    "slug",
		Threads: 0,
		Title:   "title",
		User:    "jack",
	}
}

func (self *PerfData) GetUser() *PerfUser {
	return self.GetUserData(self.GetUserIndex(), false)
}

func (self *PerfData) GetUserData(index int, with_pending bool) *PerfUser {
	return self.Users[index]
}
