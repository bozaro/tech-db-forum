package tests

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/bozaro/tech-db-forum/models"
	"github.com/drhodes/golorem"
	"github.com/go-openapi/strfmt"
	"github.com/ventu-io/go-shortid"
)

const ABC = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_."

var sid *shortid.Shortid

func init() {
	sid = shortid.MustNew(0, ABC, 1)
}

func RandomEmail() strfmt.Email {
	return strfmt.Email(RandomNickname() + "@" + lorem.Host())
}

func RandomNickname() string {
	return lorem.Word(1, 10) + "." + sid.MustGenerate()
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
		Slug:  sid.MustGenerate(),
		Title: lorem.Sentence(1, 10),
	}
}
