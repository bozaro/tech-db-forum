package tests

import "github.com/bozaro/tech-db-forum/generated/client"
import (
	"crypto/md5"
)

type Perf struct {
	c    *client.Forum
	data *PerfData
}

type PHash [16]byte
type PVersion uint32

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
	callback(&PerfSession{})
}

func Hash(data string) PHash {
	return PHash(md5.Sum([]byte(data)))
}

func (self *Perf) Run() {
	for pass := 0; pass < 100; pass++ {
		for _, p := range registeredPerfs {
			if p.FnPerf == nil {
				log.Warning(p.Name)
				continue
			}
			log.Info(p.Name)
			p.FnPerf(self)
		}
	}
}
