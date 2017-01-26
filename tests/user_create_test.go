package main

import (
	"github.com/bozaro/tech-db-forum/tests/client/operations"
	"github.com/bozaro/tech-db-forum/tests/models"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func CreateUser(t *testing.T, user *models.User) *models.User {
	if user == nil {
		user = RandomUser()
	}

	request := *user
	request.Nickname = ""

	_, err := c.Operations.UserCreate(operations.NewUserCreateParams().
		WithNickname(user.Nickname).
		WithProfile(&request).
		WithContext(Expected(t, 201, user, nil)))
	assert.Nil(t, err)

	return user
}

func CheckUser(t *testing.T, user *models.User) {
	_, err := c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(user.Nickname).
		WithContext(Expected(t, 200, user, nil)))
	assert.Nil(t, err)
}

func TestUserCreateSimple(t *testing.T) {
	CreateUser(t, nil)
}

func TestUserCreateUnicode(t *testing.T) {
	user := RandomUser()
	user.Fullname = "Маркиз О-де-Колóн"
	user.About = "Бездельник третьего разряда 😋"
	CreateUser(t, user)
	CheckUser(t, user)
}

func TestUserCreateConflict(t *testing.T) {
	user1 := CreateUser(t, nil)
	user2 := CreateUser(t, nil)

	pass := 0
	for true {
		pass++

		expected := []models.User{}
		conflict_user := RandomUser()

		modify := pass
		// Email
		switch modify % 4 {
		case 1:
			conflict_user.Email = user1.Email
			expected = append(expected, *user1)
		case 2:
			conflict_user.Email = strfmt.Email(strings.ToLower(user1.Email.String()))
			expected = append(expected, *user1)
		case 3:
			conflict_user.Email = strfmt.Email(strings.ToUpper(user1.Email.String()))
			expected = append(expected, *user1)
		}
		modify /= 4
		// Nickname
		switch modify % 5 {
		case 1:
			conflict_user.Nickname = user2.Nickname
			expected = append(expected, *user2)
		case 2:
			conflict_user.Nickname = strings.ToLower(user2.Nickname)
			expected = append(expected, *user2)
		case 3:
			conflict_user.Nickname = strings.ToUpper(user2.Nickname)
			expected = append(expected, *user2)
		case 4:
			conflict_user.Nickname = user1.Nickname
			if len(expected) == 0 {
				expected = append(expected, *user1)
			}
		default:
		}
		modify /= 5
		// Done?
		if modify != 0 {
			break
		}
		// Check
		nickname := conflict_user.Nickname
		conflict_user.Nickname = ""
		c.Operations.UserCreate(operations.NewUserCreateParams().
			WithNickname(nickname).
			WithProfile(conflict_user).
			WithContext(Expected(t, 409, &expected, func(users interface{}) interface{} {
				result := UserByNickname(reflect.ValueOf(users).Elem().Interface().([]models.User))
				sort.Sort(result)
				return result
			})))
	}
}
