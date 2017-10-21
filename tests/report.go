package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aryann/difflib"
	"html/template"
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
		// Добавляем сообщение в отчет
		if len(self.Pass) == 0 {
			self.Pass = []ReportPass{{Name: ""}}
		}
		//debug.PrintStack()
		pass := &self.Pass[len(self.Pass)-1]
		pass.Failure = fmt.Sprintf("%s", err)
		self.Result = Failed
		log.Errorf("%v", err)
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
	log.Debug("  " + message)
	return true
}

func (self *Report) RoundTrip(req *http.Request, res *http.Response, example *http.Response, delta *[]difflib.DiffRecord, err error) {
	if self.Result == Failed {
		return
	}
	if delta == nil && err == nil && self.OnlyError {
		return
	}
	reportMessage := ReportMessage{
		Failed:   delta != nil || err != nil,
		Url:      req.URL.String(),
		Request:  RequestInfo(req),
		Response: ResponseInfo(res),
		Example:  ResponseInfo(example),
	}
	if reportMessage.Failed {
		log.Warningf("Request:\n%s", reportMessage.Request.String())
		log.Warningf("Actual response:\n%s", reportMessage.Response.String())
		if example != nil {
			log.Warningf("Expected response like:\n%s", reportMessage.Example.String())
		}
		self.Result = Failed
		if delta != nil {
			reportMessage.Delta = template.HTML(DeltaToHtml(*delta))
			log.Errorf("Delta:\n%s", DeltaToText(*delta))
		}
		if err != nil {
			log.Errorf("Error:\n%s", err)
			reportMessage.Error = err.Error()
		}
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
	OnlyError bool
	Pass      []ReportPass
	SkippedBy []string
	Result    ResultType
}

type ReportPass struct {
	Name     string
	Failure  string
	Messages []ReportMessage
}

type ReportHttp struct {
	Title    string
	Header   http.Header
	BodyRaw  string
	BodyJson string
}

type ReportMessage struct {
	Failed   bool
	Url      string
	Delta    template.HTML
	Request  *ReportHttp
	Response *ReportHttp
	Example  *ReportHttp
	Error    string
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
		buf.WriteString(`<tr><td class="line-num line-num-l">`)
		if d.Delta == difflib.Common || d.Delta == difflib.LeftOnly {
			i++
			fmt.Fprintf(buf, "%d</td><td", i)
			if d.Delta == difflib.LeftOnly {
				fmt.Fprint(buf, ` class="deleted"`)
			}
			fmt.Fprintf(buf, `><pre><code class="javascript">%s</code></pre>`, d.Payload)
		} else {
			buf.WriteString("</td><td>")
		}
		buf.WriteString("</td><td")
		if d.Delta == difflib.Common || d.Delta == difflib.RightOnly {
			j++
			if d.Delta == difflib.RightOnly {
				fmt.Fprint(buf, ` class="added"`)
			}
			fmt.Fprintf(buf, `><pre><code class="javascript">%s</code></pre></td><td class="line-num line-num-r">%d`, d.Payload, j)
		} else {
			buf.WriteString(`></td><td class="line-num line-num-r">`)
		}
		buf.WriteString("</td></tr>\n")
	}
	return buf.String()
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
		Title:    req.Method + " " + req.URL.String() + " " + req.Proto,
		Header:   req.Header,
		BodyRaw:  body,
		BodyJson: PrettyJson(body),
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
	body := string(context)
	return &ReportHttp{
		Title:    res.Proto + " " + res.Status,
		Header:   res.Header,
		BodyRaw:  body,
		BodyJson: PrettyJson(body),
	}
}

func PrettyJson(body string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(body), "", "  ")
	if err != nil {
		return body
	}
	return out.String()
}

func (self *ReportHttp) String() string {
	msg := self.Title + "\n"
	for key, vals := range self.Header {
		for _, val := range vals {
			msg += key + ": " + val + "\n"
		}
	}
	msg += "\n"
	msg += self.BodyJson
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	return msg
}
