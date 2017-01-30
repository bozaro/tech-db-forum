package tests

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/drhodes/golorem"
	"github.com/go-openapi/strfmt"
	"github.com/ventu-io/go-shortid"
)

const ABC_NICK = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_."
const ABC_SLUG = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_-"

var nick_id *shortid.Shortid
var slug_id *shortid.Shortid

func init() {
	nick_id = shortid.MustNew(0, ABC_NICK, 1)
	slug_id = shortid.MustNew(0, ABC_SLUG, 1)
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
	}
}
func RandomPost() *models.Post {
	return &models.Post{
		Message: lorem.Paragraph(1, 20),
	}
}
