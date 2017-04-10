package tests

import "github.com/bozaro/tech-db-forum/generated/models"

type PerfData struct {
}

func (self *PerfData) Status() models.Status {
	return models.Status{
		Forum:  0,
		Post:   0,
		User:   0,
		Thread: 0,
	}
}
