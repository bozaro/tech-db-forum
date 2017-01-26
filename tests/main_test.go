package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bozaro/tech-db-forum/tests/client"
	"github.com/bozaro/tech-db-forum/tests/client/operations"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
)

const (
	KEY_TESTING = "testing"
	KEY_STATUS  = "expected-status"
	KEY_BODY    = "expected-body"
	KEY_FILTER  = "expected-filter"
)

type Filter func(interface{}) interface{}

type CheckerTransport struct {
	t runtime.ClientTransport
}

type CheckerRoundTripper struct {
	t      *testing.T
	code   int
	body   interface{}
	filter Filter
}

func (self *CheckerTransport) Submit(operation *runtime.ClientOperation) (interface{}, error) {
	tracker := CheckerRoundTripper{}
	if operation.Context != nil {
		tracker.t = operation.Context.Value(KEY_TESTING).(*testing.T)
		tracker.code = operation.Context.Value(KEY_STATUS).(int)
		if operation.Context.Value(KEY_BODY) != nil {
			tracker.body = operation.Context.Value(KEY_BODY)
		}
		if operation.Context.Value(KEY_FILTER) != nil {
			tracker.filter = operation.Context.Value(KEY_FILTER).(Filter)
		}
	}
	if tracker.filter == nil {
		tracker.filter = func(data interface{}) interface{} {
			return data
		}
	}
	operation.Client = &http.Client{Transport: &tracker}
	return self.t.Submit(operation)
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

func (self *CheckerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Println(*req)
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		panic(err)
	}
	if self.t != nil {
		body := []byte{}
		if res.Body != nil {
			body, err = ioutil.ReadAll(res.Body)
		}
		res.Body.Close()
		if err != nil {
			panic(err)
		}
		res.Body.Close()
		if (res.StatusCode != self.code) || !AreEqual(body, self.body, self.filter) {
			log.Println("----------------")
			log.Println(string(body))
			expected_json, _ := json.MarshalIndent(self.body, "", "  ")
			log.Println(string(expected_json))
			log.Println("++++++++++++++++")

			log.Println("Unexpected status code:", res.StatusCode, "!=", self.code, string(body))
			assert.Fail(self.t, "Ops...")
		}

		if res.Body != nil {
			res.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
	}
	return res, err
}

var c *client.Forum

func Expected(t *testing.T, statusCode int, body interface{}, prepare Filter) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KEY_TESTING, t)
	ctx = context.WithValue(ctx, KEY_STATUS, statusCode)
	if body != nil {
		ctx = context.WithValue(ctx, KEY_BODY, body)
	}
	if prepare != nil {
		ctx = context.WithValue(ctx, KEY_FILTER, prepare)
	}
	return ctx

}

func TestStatusSmoke(t *testing.T) {
	c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(t, 200, nil, nil)))
}

//go:generate swagger generate client --target . --spec ../swagger.yml
func TestMain(m *testing.M) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:5000").WithSchemes([]string{"http"})
	transport := httptransport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
	c = client.New(&CheckerTransport{transport}, nil)
	os.Exit(m.Run())
}
