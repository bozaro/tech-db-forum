package tests

import (
	"bytes"
	"fmt"
	"github.com/aryann/difflib"
	"html/template"
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
		//self.messages = append(self.messages, fmt.Sprint(err))
		self.Result = Failed
	}
}

func (self *Report) Skip(message string) {
	self.SkippedBy = append(self.SkippedBy, message)
	self.Result = Skipped
}

func (self *Report) Checkpoint(message string) bool {
	if self.Result == Failed {
		return false
	}
	self.Pass = append(self.Pass, ReportPass{Name: message})
	log.Println("  " + message)
	return true
}

func (self *Report) RoundTrip(req *http.Request, res *http.Response, example *http.Response, delta *[]difflib.DiffRecord) {
	if self.Result == Failed {
		return
	}
	msg := ""
	if delta != nil {
		msg += "!!! ERROR:\n"
		msg += DeltaToText(*delta)
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
	if delta != nil {
		if example != nil {
			msg += "<<< EXPECTED RESPONSE EXAMPLE:\n"
			msg += ResponseToText(example)
		}
		self.Result = Failed
	}
	if len(self.Pass) == 0 {
		self.Pass = []ReportPass{{Name: ""}}
	}
	pass := &self.Pass[len(self.Pass)-1]

	reportMessage := ReportMessage{}
	if delta != nil {
		reportMessage.Delta = template.HTML(DeltaToHtml(*delta))
	}
	pass.Messages = append(pass.Messages, reportMessage)
}

type Report struct {
	Checker   Checker
	Pass      []ReportPass
	SkippedBy []string
	Result    ResultType
}

type ReportPass struct {
	Name     string
	Messages []ReportMessage
}

type ReportMessage struct {
	Delta template.HTML
}

func DeltaToText(delta []difflib.DiffRecord) string {
	result := make([]string, len(delta))
	for i, item := range delta {
		switch item.Delta {
		case difflib.LeftOnly:
			result[i] = Colorize(31, item.String())
		case difflib.RightOnly:
			result[i] = Colorize(32, item.String())
		default:
			result[i] = item.String()
		}
	}
	return strings.Join(result, "\n")
}

func DeltaToHtml(delta []difflib.DiffRecord) string {
	buf := bytes.NewBufferString("")
	i, j := 0, 0
	for _, d := range delta {
		buf.WriteString(`<tr><td class="line-num">`)
		if d.Delta == difflib.Common || d.Delta == difflib.LeftOnly {
			i++
			fmt.Fprintf(buf, "%d</td><td", i)
			if d.Delta == difflib.LeftOnly {
				fmt.Fprint(buf, ` class="deleted"`)
			}
			fmt.Fprintf(buf, "><pre>%s</pre>", d.Payload)
		} else {
			buf.WriteString("</td><td>")
		}
		buf.WriteString("</td><td")
		if d.Delta == difflib.Common || d.Delta == difflib.RightOnly {
			j++
			if d.Delta == difflib.RightOnly {
				fmt.Fprint(buf, ` class="added"`)
			}
			fmt.Fprintf(buf, `><pre>%s</pre></td><td class="line-num">%d`, d.Payload, j)
		} else {
			buf.WriteString("></td><td>")
		}
		buf.WriteString("</td></tr>\n")
	}
	return buf.String()
}

func (self *Report) Show() {
	/*for _, message := range self.messages {
		log.Println(message)
	}*/
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
