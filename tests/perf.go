package tests

import "github.com/bozaro/tech-db-forum/generated/client"
import (
	"crypto/md5"
)

type Perf struct {
	c    *client.Forum
	data *PerfData
}

type PerfHash [16]byte

type PerfTest struct {
	Name   string
	Mode   PerfMode
	Weight PerfWeight
	FnPerf func(p *Perf)
}

type PerfMode int

const (
	ModeRead PerfMode = iota
	ModeWrite
)

type PerfWeight int

const (
	WeightRare   PerfWeight = 1
	WeightNormal            = 10
)

var registeredPerfs []PerfTest

func PerfRegister(test PerfTest) {
	registeredPerfs = append(registeredPerfs, test)
}

func (self *Perf) Validate(callback func(validator PerfValidator)) {

}

func Hash(data string) PerfHash {
	return PerfHash(md5.Sum([]byte(data)))
}
