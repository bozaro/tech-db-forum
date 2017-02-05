package tests

import (
	"bytes"
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/models"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"time"
)

type UserByNickname []models.User

func (a UserByNickname) Len() int           { return len(a) }
func (a UserByNickname) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UserByNickname) Less(i, j int) bool { return a[i].Nickname < a[j].Nickname }

type ThreadByCreated []models.Thread

func (a ThreadByCreated) Len() int      { return len(a) }
func (a ThreadByCreated) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ThreadByCreated) Less(i, j int) bool {
	return time.Time(a[i].Created).Before(time.Time(a[j].Created))
}

func CheckNil(err interface{}) {
	if err != nil {
		panic(err)
	}
}

func CheckIsType(expectedType interface{}, object interface{}) {
	if !ObjectsAreEqual(reflect.TypeOf(object), reflect.TypeOf(expectedType)) {
		panic(fmt.Sprintf("Object expected to be of type %v, but was %v", reflect.TypeOf(expectedType), reflect.TypeOf(object)))
	}
}

func ObjectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}
	return reflect.DeepEqual(expected, actual)
}

func GetBody(stream *io.ReadCloser) ([]byte, error) {
	if *stream == nil {
		return []byte{}, nil
	}
	ibody := *stream
	defer ibody.Close()
	body, err := ioutil.ReadAll(ibody)
	if err != nil {
		return body, err
	}
	*stream = ioutil.NopCloser(bytes.NewReader(body))
	return body, nil

}

func ModifyCase(modify *int, source string) string {
	v := *modify % 3
	*modify /= 3
	switch v {
	case 0:
		return source
	case 1:
		return strings.ToLower(source)
	case 2:
		return strings.ToUpper(source)
	default:
		panic("Unexpected value")
	}
}
