package tests

import (
	"bytes"
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/models"
	"io"
	"io/ioutil"
	"reflect"
)

type UserByNickname []models.User

func (a UserByNickname) Len() int           { return len(a) }
func (a UserByNickname) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UserByNickname) Less(i, j int) bool { return a[i].Nickname < a[j].Nickname }

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
