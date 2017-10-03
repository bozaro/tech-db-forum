package tests

import (
	"github.com/bozaro/golorem"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/go-openapi/strfmt"
	"math/rand"
	"time"
)

const ABC_NICK = "0123456789abcdefghijklmnopqrstuvwxyz"
const ABC_SLUG = "0123456789abcdefghijklmnopqrstuvwxyz_-"
const POST_FAKE_ID int64 = 2139800938
const THREAD_FAKE_ID = "2139800939"

var nick_id *Shortid
var slug_id *Shortid

type Factory struct {
	lorem *lorem.Lorem
}

func NewFactory() *Factory {
	return &Factory{
		lorem: lorem.New(),
	}
}

func init() {
	nick_id = NewShortid(ABC_NICK)
	slug_id = NewShortid(ABC_SLUG)
}

func RandomTime() time.Time {
	year := int64(time.Hour) * 24 * 365
	return time.
		Now().
		Add(time.Duration(rand.Int63n(year*2) - year)).
		Round(time.Millisecond)
}

func (self *Factory) RandomEmail() strfmt.Email {
	return strfmt.Email(self.RandomNickname() + "@" + self.lorem.Host())
}

func (self *Factory) RandomNickname() string {
	return self.lorem.Word(1, 10) + "." + nick_id.Generate()
}

func (self *Factory) RandomUser() *models.User {
	return &models.User{
		About:    self.lorem.Paragraph(1, 10),
		Email:    self.RandomEmail(),
		Fullname: self.lorem.FullName(-1),
		Nickname: self.RandomNickname(),
	}
}

func (self *Factory) RandomForum() *models.Forum {
	return &models.Forum{
		Posts: 0,
		Slug:  slug_id.Generate(),
		Title: self.lorem.Sentence(1, 10),
	}
}
func (self *Factory) RandomSlug() string {
	return slug_id.Generate()
}
func (self *Factory) RandomThread() *models.Thread {
	created := strfmt.DateTime(RandomTime())
	return &models.Thread{
		Message: self.lorem.Paragraph(1, 20),
		Slug:    self.RandomSlug(),
		Title:   self.lorem.Sentence(1, 10),
		Created: &created,
	}
}

func (self *Factory) RandomPost() *models.Post {
	return &models.Post{
		Message:  self.lorem.Paragraph(1, 20),
		IsEdited: false,
	}
}

func (self *Factory) RandomPosts(count int) []*models.Post {
	posts := make([]*models.Post, count)
	for i := range posts {
		posts[i] = self.RandomPost()
	}
	return posts
}
