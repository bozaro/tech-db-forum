package tests

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ResultType int

const (
	Failed ResultType = iota
	Skipped
	Success
)

func (self *Report) AddError(err interface{}) {
	if self.Result != Failed {
		self.messages = append(self.messages, fmt.Sprint(err))
		self.Result = Failed
	}
}
func (self *Report) Skip(message string) {
	self.messages = append(self.messages, message)
	self.Result = Skipped
}
func (self *Report) Checkpoint(message string) bool {
	if self.Result == Failed {
		return false
	}
	self.messages = []string{}
	log.Println("  " + message)
	return true
}

func (self *Report) RoundTrip(req *http.Request, res *http.Response, example *http.Response, message *string) {
	if self.Result == Failed {
		return
	}
	if message == nil {
		// TODO: Сделать адекватный вывод ошибок
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
		self.Result = Failed
	}
	self.messages = append(self.messages, msg)
}

type Report struct {
	Checker  Checker
	messages []string
	Result   ResultType
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
