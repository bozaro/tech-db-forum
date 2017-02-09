package tests

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/drhodes/golorem"
	"github.com/go-openapi/strfmt"
	"github.com/ventu-io/go-shortid"
	"math/rand"
	"time"
)

const ABC_NICK = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_."
const ABC_SLUG = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-"
const POST_FAKE_ID int64 = 2139800938
const THREAD_FAKE_ID = "2139800939"

var nick_id *shortid.Shortid
var slug_id *shortid.Shortid

func init() {
	nick_id = shortid.MustNew(0, ABC_NICK, 1)
	slug_id = shortid.MustNew(0, ABC_SLUG, 1)
}

func RandomTime() time.Time {
	year := int64(time.Hour) * 24 * 365
	return time.
		Now().
		Add(time.Duration(rand.Int63n(year*2) - year)).
		Round(time.Millisecond)
}

func RandomMarker() string {
	return slug_id.MustGenerate()
}

func RandomEmail() strfmt.Email {
	return strfmt.Email(RandomNickname() + "@" + lorem.Host())
}

func RandomNickname() string {
	return lorem.Word(1, 10) + "." + nick_id.MustGenerate()
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
		Slug:  slug_id.MustGenerate(),
		Title: lorem.Sentence(1, 10),
	}
}

func RandomThread() *models.Thread {
	return &models.Thread{
		Message: lorem.Paragraph(1, 20),
		Slug:    slug_id.MustGenerate(),
		Title:   lorem.Sentence(1, 10),
		Created: strfmt.DateTime(RandomTime()),
	}
}
func RandomPost() *models.Post {
	edited := false
	return &models.Post{
		Message:  lorem.Paragraph(1, 20),
		IsEdited: &edited,
	}
}
