package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/models"
	"strings"
)

type Modify int

func (self *Modify) Int(n int) int {
	i := int(*self)
	*self = Modify(i / n)
	return i % n
}

func (self *Modify) Bool() bool {
	return self.Int(2) > 0
}

func (self *Modify) Case(source string) string {
	switch self.Int(3) {
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
func (self *Modify) SlugOrId(thread *models.Thread) string {
	switch self.Int(4) {
	case 0:
		return thread.Slug
	case 1:
		return strings.ToLower(thread.Slug)
	case 2:
		return strings.ToUpper(thread.Slug)
	case 3:
		return fmt.Sprintf("%d", thread.ID)
	default:
		panic("Unexpected value")
	}
}
