package tests

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/drhodes/golorem"
	"github.com/go-openapi/strfmt"
	"math/rand"
	"time"
)

const ABC_NICK = "0123456789abcdefghijklmnopqrstuvwxyz_."
const ABC_SLUG = "0123456789abcdefghijklmnopqrstuvwxyz_-"
const POST_FAKE_ID int64 = 2139800938
const THREAD_FAKE_ID = "2139800939"

var nick_id *Shortid
var slug_id *Shortid

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

func RandomEmail() strfmt.Email {
	return strfmt.Email(RandomNickname() + "@" + lorem.Host())
}

func RandomNickname() string {
	return lorem.Word(1, 10) + "." + nick_id.Generate()
}

func RandomUser() *models.User {
	return &models.User{
		About:    lorem.Paragraph(1, 10),
		Email:    RandomEmail(),
		Fullname: randomdata.FullName(-1),
		Nickname: RandomNickname(),
	}
}

func RandomForum() *models.Forum {
	return &models.Forum{
		Posts: 0,
		Slug:  slug_id.Generate(),
		Title: lorem.Sentence(1, 10),
	}
}

func RandomThread() *models.Thread {
	created := strfmt.DateTime(RandomTime())
	return &models.Thread{
		Message: lorem.Paragraph(1, 20),
		Slug:    slug_id.Generate(),
		Title:   lorem.Sentence(1, 10),
		Created: &created,
	}
}

func RandomPost() *models.Post {
	return &models.Post{
		Message:  lorem.Paragraph(1, 20),
		IsEdited: false,
	}
}

func RandomPosts(count int) []*models.Post {
	posts := make([]*models.Post, count)
	for i := range posts {
		posts[i] = RandomPost()
	}
	return posts
}
