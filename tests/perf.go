package tests

import "github.com/bozaro/tech-db-forum/generated/client"
import (
	"crypto/md5"
	"io"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Perf struct {
	c    *client.Forum
	data *PerfData
}

type PerfTest struct {
	Name   string
	Mode   PerfMode
	Weight PerfWeight
	FnPerf func(p *Perf, f *Factory)
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

var (
	registeredPerfsWeight int32 = 0
	registeredPerfs       []PerfTest
)

func PerfRegister(test PerfTest) {
	registeredPerfs = append(registeredPerfs, test)
	registeredPerfsWeight += int32(test.Weight)
}

func (self *Perf) Validate(callback func(validator PerfValidator)) {
	callback(&PerfSession{})
}

func Hash(data string) PHash {
	return PHash(md5.Sum([]byte(data)))
}

func GetRandomPerfTest() *PerfTest {
	index := rand.Int31n(registeredPerfsWeight)
	for _, item := range registeredPerfs {
		index -= int32(item.Weight)
		if index < 0 {
			return &item
		}
	}
	panic("Invalid state")
}

func (self *Perf) Run(threads int) {
	var done int32 = 0
	var counter int64 = 0
	// spawn four worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			f := NewFactory()
			for {
				if atomic.LoadInt32(&done) != 0 {
					break
				}
				p := GetRandomPerfTest()
				p.FnPerf(self, f)
				atomic.AddInt64(&counter, 1)
			}
			wg.Done()
		}()
	}

	lst := atomic.LoadInt64(&counter)
	step := time.Duration(10)
	for i := 0; i < 18; i++ {
		time.Sleep(time.Second * step)
		cur := atomic.LoadInt64(&counter)
		log.Infof("Requests per second: %5.02f", float32(cur-lst)/float32(step))
		lst = cur
	}
	done = 1

	// wait for the workers to finish
	wg.Wait()
}

func (self *Perf) Load(reader io.Reader) error {
	var err error
	self.data, err = LoadPerfData(reader)
	return err
}
func (self *Perf) Save(writer io.Writer) error {
	return self.data.Save(writer)
}
