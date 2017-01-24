package main

import (
	"github.com/bozaro/tech-db-forum/client"
	"github.com/bozaro/tech-db-forum/client/operations"
	"github.com/bozaro/tech-db-forum/models"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"log"
	"os"
	"testing"
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

var c *client.Forum

func (self *CheckerTransport) Submit(operation *runtime.ClientOperation) (interface{}, error) {
	log.Println(operation.Method)
	return self.t.Submit(operation)
}

func TestStatusSmoke(t *testing.T) {
	c.Operations.Status(operations.NewStatusParams().
		WithContext(context.WithValue(context.Background(), "expected-status", 200)))
}

func TestUserCreateSimple(t *testing.T) {
	expected_user := models.User{
		About:    "",
		Email:    "j.sparrow@see",
		Fullname: "Jack Sparrow",
		Nickname: "j.sparrow",
	}
	user := expected_user
	user.Nickname = ""

	log.Println(42)

	created_user, err := c.Operations.UserCreate(operations.NewUserCreateParams().
		WithNickname(expected_user.Nickname).
		WithProfile(&user).
		WithContext(context.WithValue(context.Background(), "expected-status", 201)))
	assert.Nil(t, err)
	assert.NotNil(t, created_user.Payload)
	assert.Equal(t, *created_user.Payload, expected_user)
}

// go:generate swagger generate client --target . --spec swagger.yml
func TestMain(m *testing.M) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:5000").WithSchemes([]string{"http"})
	transport := httptransport.New(cfg.Host, cfg.BasePath, cfg.Schemes)
	c = client.New(&CheckerTransport{transport}, nil)
	os.Exit(m.Run())
}
