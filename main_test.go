package main

import (
	"bytes"
	"encoding/json"
	"github.com/bozaro/tech-db-forum/client"
	"github.com/bozaro/tech-db-forum/client/operations"
	"github.com/bozaro/tech-db-forum/models"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
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
)

type Checker struct {
	// Имя текущей проверки.
	Name string
	// Функция для текущей проверки.
	FnCheck func(c *client.Forum)
	// Тесты, без которых проверка не имеет смысл.
	Deps []string
}

type CheckerTransport struct {
	t runtime.ClientTransport
}

type CheckerRoundTripper struct {
	t    *testing.T
	code int
	body interface{}
}

func (self *CheckerTransport) Submit(operation *runtime.ClientOperation) (interface{}, error) {
	tracker := CheckerRoundTripper{}
	if operation.Context != nil {
		tracker.t = operation.Context.Value(KEY_TESTING).(*testing.T)
		tracker.code = operation.Context.Value(KEY_STATUS).(int)
		tracker.body = operation.Context.Value(KEY_BODY)
	}
	operation.Client = &http.Client{Transport: &tracker}
	return self.t.Submit(operation)
}

func AreEqual(data []byte, expected interface{}) bool {
	if expected == nil {
		return true
	}
	var actual interface{} = reflect.New(reflect.TypeOf(expected).Elem()).Interface()
	log.Println("================")
	log.Println(reflect.TypeOf(expected))
	log.Println(reflect.TypeOf(actual))
	log.Println("----------------")
	log.Println(string(data))
	d, _ := json.MarshalIndent(expected, "", "  ")
	log.Println(string(d))
	log.Println("++++++++++++++++")
	if err := json.Unmarshal(data, actual); err != nil {
		log.Println(err)
		return false
	}
	log.Println(actual)
	log.Println(expected)
	return reflect.DeepEqual(actual, expected)
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
		if (res.StatusCode != self.code) || !AreEqual(body, self.body) {
			log.Println("Unexpected status code:", res.StatusCode, "!=", self.code, string(body))
		}

		if res.Body != nil {
			res.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
	}
	return res, err
}

var c *client.Forum

func Expected(t *testing.T, statusCode int, body interface{}) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, KEY_TESTING, t)
	ctx = context.WithValue(ctx, KEY_STATUS, statusCode)
	if body != nil {
		ctx = context.WithValue(ctx, KEY_BODY, body)
	}
	return ctx

}

func TestClearSmoke(t *testing.T) {
	c.Operations.Clear(operations.NewClearParams().
		WithContext(Expected(t, 200, nil)))
}

func TestStatusSmoke(t *testing.T) {
	c.Operations.Status(operations.NewStatusParams().
		WithContext(Expected(t, 200, nil)))
}

func TestUserCreateSimple(t *testing.T) {
	expected_user := models.User{
		About:    "",
		Email:    "j.sparrow@see",
		Fullname: "Jack Sparrow",
		Nickname: "j.sparrow3",
	}
	user := expected_user
	user.Nickname = ""

	_, err := c.Operations.UserCreate(operations.NewUserCreateParams().
		WithNickname(expected_user.Nickname).
		WithProfile(&user).
		WithContext(Expected(t, 201, &expected_user)))
	assert.Nil(t, err)

	conflict_user := models.User{
		About:    "",
		Email:    "j.sparrow@see.te",
		Fullname: "Jack Sparrow",
		Nickname: "j.sparrow3",
	}
	_, err = c.Operations.UserCreate(operations.NewUserCreateParams().
		WithNickname(expected_user.Nickname).
		WithProfile(&conflict_user).
		WithContext(Expected(t, 409, &[]models.User{expected_user})))
}

// go:generate swagger generate client --target . --spec swagger.yml
func TestMain(m *testing.M) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:5000").WithSchemes([]string{"http"})
	transport := httptransport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
	c = client.New(&CheckerTransport{transport}, nil)
	os.Exit(m.Run())
}
