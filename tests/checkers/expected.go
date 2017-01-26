package checkers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

const (
	KEY_STATUS = "expected-status"
	KEY_BODY   = "expected-body"
	KEY_FILTER = "expected-filter"
)

type Filter func(interface{}) interface{}

type Validator struct {
	code   int
	body   interface{}
	filter Filter
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

func NewValidator(ctx context.Context) *Validator {
	v := Validator{}
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
			log.Println(r)
		}
	}()
	if self.code != 0 {
		body := []byte{}
		if res.Body != nil {
			ibody := res.Body
			defer ibody.Close()
			var err error
			body, err = ioutil.ReadAll(ibody)
			if err != nil {
				panic(err)
			}
		}

		if (res.StatusCode != self.code) || !AreEqual(body, self.body, self.filter) {
			log.Println("----------------")
			log.Println(string(body))
			expected_json, _ := json.MarshalIndent(self.body, "", "  ")
			log.Println(string(expected_json))
			log.Println("++++++++++++++++")

			log.Println("Unexpected status code:", res.StatusCode, "!=", self.code, string(body))
			panic("Ops...")
		}

		if res.Body != nil {
			res.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
	}
	return true
}

func (self *Validator) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Println(*req)
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if self.validate(req, res) {
		return res, nil
	}
	fmt.Println("!!!!!!!!!!!!!!!!!!")
	return nil, errors.New("Unexpected error")
}

func ToJson(obj interface{}) string {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func AreEqual(data []byte, expected interface{}, prepare Filter) bool {
	if expected == nil {
		return true
	}
	var actual interface{} = reflect.New(reflect.TypeOf(expected).Elem()).Interface()
	if err := json.Unmarshal(data, actual); err != nil {
		log.Println(err)
		return false
	}

	expected_json := ToJson(prepare(expected))
	actual_json := ToJson(prepare(actual))
	if expected_json == actual_json {
		return true
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(expected_json, actual_json, false)
	fmt.Println("====>")
	fmt.Println(dmp.DiffPrettyText(diffs))
	fmt.Println("====<")
	return false
}
