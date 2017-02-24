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

	reportMessage := ReportMessage{
		Url:      req.URL.String(),
		Request:  RequestInfo(req),
		Response: ResponseInfo(res),
		Example:  ResponseInfo(example),
	}
	if delta != nil {
		reportMessage.Delta = template.HTML(DeltaToHtml(*delta))
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
	msg += reportMessage.Request.String()
	if res != nil {
		msg += "<<< ACTUAL RESPONSE:\n"
		msg += reportMessage.Response.String()
	}
	if delta != nil {
		if example != nil {
			msg += "<<< EXPECTED RESPONSE EXAMPLE:\n"
			msg += reportMessage.Example.String()
		}
		self.Result = Failed
	}
	// Добавляем сообщение в отчет
	if len(self.Pass) == 0 {
		self.Pass = []ReportPass{{Name: ""}}
	}
	pass := &self.Pass[len(self.Pass)-1]
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

type ReportHttp struct {
	Title  string
	Header http.Header
	Body   string
}

type ReportMessage struct {
	Url      string
	Delta    template.HTML
	Request  *ReportHttp
	Response *ReportHttp
	Example  *ReportHttp
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

func RequestInfo(req *http.Request) *ReportHttp {
	context, err := GetBody(&req.Body)
	body := ""
	if err == nil {
		body += string(context)
	}
	if len(body) > 0 && !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	return &ReportHttp{
		Title:  req.Method + " " + req.URL.String() + " " + req.Proto,
		Header: req.Header,
		Body:   body,
	}
}

func ResponseInfo(res *http.Response) *ReportHttp {
	if res == nil {
		return nil
	}
	context, err := GetBody(&res.Body)
	if err != nil {
		panic(err)
	}
	return &ReportHttp{
		Title:  res.Proto + " " + res.Status,
		Header: res.Header,
		Body:   string(context),
	}
}

func (self *ReportHttp) String() string {
	msg := self.Title + "\n"
	for key, vals := range self.Header {
		for _, val := range vals {
			msg += key + ": " + val + "\n"
		}
	}
	msg += "\n"
	msg += self.Body
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	return msg
}
