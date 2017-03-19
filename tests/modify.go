package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/models"
)

type Modify int

func (self *Modify) Int(n int) int {
	i := int(*self)
	*self = Modify(i / n)
	return i % n
}

func InvertCase(str string) string {
	data := []byte(str)
	for i, c := range data {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			c ^= 'a' - 'A'
		}
		data[i] = c
	}
	return string(data)
}

func (self *Modify) Bool() bool {
	return self.Int(2) > 0
}

func (self *Modify) NullableBool() *bool {
	switch self.Int(3) {
	case 0:
		return nil
	case 1:
		v := false
		return &v
	case 2:
		v := true
		return &v
	default:
		panic("Unexpected value")
	}
}

func (self *Modify) Case(source string) string {
	if self.Bool() {
		return InvertCase(source)
	} else {
		return source
	}
}

func (self *Modify) SlugOrId(thread *models.Thread) string {
	switch self.Int(3) {
	case 0:
		return thread.Slug
	case 1:
		return InvertCase(thread.Slug)
	case 2:
		return fmt.Sprintf("%d", thread.ID)
	default:
		panic("Unexpected value")
	}
}
