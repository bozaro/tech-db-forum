package checkers

import (
	"fmt"
	"log"
	"net/http"
)

const (
	REPORT_FAILED  = 0
	REPORT_SKIPPED = 1
	REPORT_SUCCESS = 2
)

func (self *Report) AddError(err interface{}) {
	self.messages = append(self.messages, fmt.Sprint(err))
	self.result = REPORT_FAILED
}

func (self *Report) RoundTrip(req *http.Request, res *http.Response, example *http.Response, message *string) {
	if message != nil {
		self.AddError(message)
	}
}

type Report struct {
	messages []string
	result   int
}

func (self *Report) Show() {
	for _, message := range self.messages {
		log.Println(message)
	}
}
