package tests

import "github.com/bozaro/tech-db-forum/generated/client"
import (
	"hash/crc32"
	"io"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Perf struct {
	c        *client.Forum
	data     *PerfData
	validate float32
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
	WeightNever  PerfWeight = 0
	WeightRare              = 1
	WeightNormal            = 10
)

var (
	registeredPerfsWeight int32 = 0
	registeredPerfs       []PerfTest
)
var crc32q = crc32.MakeTable(0xD5828281)

func PerfRegister(test PerfTest) {
	registeredPerfs = append(registeredPerfs, test)
	registeredPerfsWeight += int32(test.Weight)
}

func (self *Perf) Session() *PerfSession {
	return &PerfSession{validate: rand.Float32() < self.validate}
}

func Hash(data string) PHash {
	return PHash(crc32.Checksum([]byte(data), crc32q))
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

func (self *Perf) Run(threads int, duration int, step int) float64 {
	var done int32 = 0
	var counter int64 = 0
	// spawn four worker goroutines
	var wg sync.WaitGroup
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		go func() {
			f := NewFactory()
			runtime.LockOSThread()
			for {
				if atomic.LoadInt32(&done) != 0 {
					break
				}
				p := GetRandomPerfTest()
				p.FnPerf(self, f)
				atomic.AddInt64(&counter, 1)
			}
			runtime.UnlockOSThread()
			wg.Done()
		}()
	}

	lst := atomic.LoadInt64(&counter)
	best := -1.0
	cnt := duration / step
	for {
		time.Sleep(time.Second * time.Duration(step))
		cur := atomic.LoadInt64(&counter)
		rps := float64(cur-lst) / float64(step)
		if best < rps {
			best = rps
		}
		log.Infof("Requests per second: %5.02f", rps)
		lst = cur
		cnt--
		if duration >= 0 && cnt <= 0 {
			break
		}
	}
	log.Infof("Requests per second: %5.02f (best)", best)
	done = 1

	// wait for the workers to finish
	wg.Wait()
	return best
}

func (self *Perf) Load(reader io.Reader) error {
	var err error
	self.data, err = LoadPerfData(reader)
	return err
}
func (self *Perf) Save(writer io.Writer) error {
	return self.data.Save(writer)
}
