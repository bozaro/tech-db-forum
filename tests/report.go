package tests

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	REPORT_FAILED  = 0
	REPORT_SKIPPED = 1
	REPORT_SUCCESS = 2
)

func (self *Report) AddError(err interface{}) {
	if self.result != REPORT_FAILED {
		self.messages = append(self.messages, fmt.Sprint(err))
		self.result = REPORT_FAILED
	}
}
func (self *Report) Skip(message string) {
	self.messages = append(self.messages, message)
	self.result = REPORT_SKIPPED
}
func (self *Report) Checkpoint(message string) bool {
	if self.result == REPORT_FAILED {
		return false
	}
	self.messages = []string{}
	log.Println("  " + message)
	return true
}

func (self *Report) RoundTrip(req *http.Request, res *http.Response, example *http.Response, message *string) {
	if self.result == REPORT_FAILED {
		return
	}
	msg := ""
	if message != nil {
		msg += "!!! ERROR:\n"
		msg += *message
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
	}
	msg += ">>> REQUEST:\n"
	msg += RequestToText(req)
	if res != nil {
		msg += "<<< ACTUAL RESPONSE:\n"
		msg += ResponseToText(res)
	}
	if message != nil {
		if example != nil {
			msg += "<<< EXPECTED RESPONSE EXAMPLE:\n"
			msg += ResponseToText(example)
		}
		self.result = REPORT_FAILED
	}
	self.messages = append(self.messages, msg)
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

func RequestToText(req *http.Request) string {
	msg := req.Method + " " + req.URL.String() + " " + req.Proto + "\n"
	for key, vals := range req.Header {
		for _, val := range vals {
			msg += key + ": " + val + "\n"
		}
	}
	msg += "\n"

	body, err := GetBody(&req.Body)
	if err == nil {
		msg += string(body)
	}
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	return msg
}

func ResponseToText(res *http.Response) string {
	msg := res.Proto + " " + res.Status + "\n"
	for key, vals := range res.Header {
		for _, val := range vals {
			msg += key + ": " + val + "\n"
		}
	}
	msg += "\n"

	body, err := GetBody(&res.Body)
	if err == nil {
		msg += string(body)
	}
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	return msg
}
