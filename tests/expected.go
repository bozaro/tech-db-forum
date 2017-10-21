package tests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aryann/difflib"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/go-openapi/swag"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

const (
	KEY_STATUS = "expected-status"
	KEY_BODY   = "expected-body"
	KEY_FILTER = "expected-filter"
	KEY_SKIP   = "dont-check"

	MaxIdleConnections int = 20
)

type Filter func(interface{}) interface{}

type Validator struct {
	report *Report
	code   int
	body   interface{}
	filter Filter
}

var HttpTransport *http.Transport

func init() {
	HttpTransport = &http.Transport{
		MaxIdleConnsPerHost: MaxIdleConnections,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
}

func ExpectedError(statusCode int, format string, a ...interface{}) context.Context {
	message := fmt.Sprintf(format, a...)
	return Expected(statusCode, &models.Error{Message: message}, func(data interface{}) interface{} {
		e := data.(*models.Error)
		e.Message = message
		return data
	})
}

func Expected(statusCode int, body interface{}, prepare Filter) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KEY_STATUS, statusCode)
	if body != nil {
		ctx = context.WithValue(ctx, KEY_BODY, body)
	}
	if prepare != nil {
		ctx = context.WithValue(ctx, KEY_FILTER, prepare)
	}
	return ctx
}

func NewValidator(ctx context.Context, report *Report) *Validator {
	v := Validator{report: report}
	if ctx != nil {
		if ctx.Value(KEY_STATUS) != nil {
			v.code = ctx.Value(KEY_STATUS).(int)
			if ctx.Value(KEY_BODY) != nil {
				v.body = ctx.Value(KEY_BODY)
			}
			if ctx.Value(KEY_FILTER) != nil {
				v.filter = ctx.Value(KEY_FILTER).(Filter)
			}
		}
	}
	if v.filter == nil {
		v.filter = func(data interface{}) interface{} {
			return data
		}
	}
	return &v
}

func (self *Validator) validate(req *http.Request, res *http.Response) bool {
	defer func() {
		if r := recover(); r != nil {
			self.report.AddError(r)
		}
	}()
	if self.code != 0 {
		body, err := GetBody(&res.Body)
		if err != nil {
			panic(err)
		}

		if res.StatusCode != self.code {
			self.report.RoundTrip(req, res, self.Example(req), &[]difflib.DiffRecord{
				{
					Delta:   difflib.LeftOnly,
					Payload: fmt.Sprintf("Status: %d %s", self.code, http.StatusText(self.code)),
				},
				{
					Delta:   difflib.RightOnly,
					Payload: fmt.Sprintf("Status: %d %s", res.StatusCode, http.StatusText(res.StatusCode)),
				},
			}, nil)
			return false
		}
		delta, err := GetDelta(body, self.body, self.filter)
		if (res.StatusCode != self.code) || (delta != nil) || (err != nil) {
			self.report.RoundTrip(req, res, self.Example(req), delta, err)
			return false
		}
		self.report.RoundTrip(req, res, nil, nil, nil)
	}
	return true
}

func (self *Validator) Example(req *http.Request) *http.Response {
	if self.body == "" {
		return nil
	}
	json_body := []byte(ToJsonPretty(self.body))
	return &http.Response{
		Proto:      req.Proto,
		StatusCode: self.code,
		Status:     fmt.Sprintf("%d %s", self.code, http.StatusText(self.code)),
		Header: http.Header{
			"Content-Type":   []string{"application/json"},
			"Content-Length": []string{fmt.Sprintf("%d", len(json_body))},
		},
		Body: ioutil.NopCloser(bytes.NewReader(json_body)),
	}
}

func (self *Validator) RoundTrip(req *http.Request) (*http.Response, error) {
	req_body, err := GetBody(&req.Body)
	if err != nil {
		self.report.AddError(err)
		return nil, err
	}
	res, err := HttpTransport.RoundTrip(req)
	if err != nil {
		self.report.AddError(err)
		return nil, err
	}
	if req.Body != nil {
		req.Body = ioutil.NopCloser(bytes.NewReader(req_body))
	}
	if self.validate(req, res) {
		return res, nil
	}
	return nil, errors.New("Unexpected error")
}

func ToJsonPretty(obj interface{}) string {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func GetDiff(actual string, expected string) *[]difflib.DiffRecord {
	if actual == expected {
		return nil
	}
	delta := difflib.Diff(
		strings.Split(expected, "\n"),
		strings.Split(actual, "\n"),
	)
	return &delta
}

func Colorize(color int, message string) string {
	if runtime.GOOS == "windows" {
		return message
	}
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, message)
}

func GetDelta(data []byte, expected interface{}, prepare Filter) (*[]difflib.DiffRecord, error) {
	if expected == nil {
		return nil, nil
	}
	expected_obj := prepare(expected)
	var actual interface{} = reflect.New(reflect.TypeOf(expected).Elem()).Interface()
	if err := swag.ReadJSON(data, actual); err != nil {
		return nil, err
	}

	actual_obj := prepare(actual)
	actual_json, _ := swag.WriteJSON(actual_obj)
	expected_json, _ := swag.WriteJSON(expected_obj)
	if bytes.Equal(expected_json, actual_json) {
		return nil, nil
	}
	return GetDiff(ToJsonPretty(actual_obj), ToJsonPretty(expected_obj)), nil
}
