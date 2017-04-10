package tests

import "github.com/bozaro/tech-db-forum/generated/client"

type Perf struct {
	c    *client.Forum
	data *PerfData
}

type PerfTest struct {
	Name   string
	Mode   PerfMode
	Weight PerfWeight
}

type PerfMode int

const (
	ModeRead PerfMode = iota
	ModeWrite
)

type PerfWeight int

const (
	WeightRare PerfWeight = 1
)

var registeredPerfs []PerfTest

func PerfRegister(test PerfTest) {
	registeredPerfs = append(registeredPerfs, test)
}

func (self *Perf) Validate(callback func(validator PerfValidator)) {

}
