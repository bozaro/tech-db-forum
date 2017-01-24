package main

import (
	"github.com/ventu-io/go-shortid"
	"github.com/go-openapi/strfmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/bozaro/tech-db-forum/tests/models"
	"strings"
)

const ABC = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_."

var sid *shortid.Shortid

func init() {
	sid = shortid.MustNew(0, ABC, 1)
}

func RandomEmail() strfmt.Email {
	email := strings.Split(randomdata.Email(), "@")[1]
	return strfmt.Email(sid.MustGenerate() + "@" + email)
}

func RandomNickname() string {
	return sid.MustGenerate()
}

func RandomUser() models.User {
	return models.User{
		About:    randomdata.Paragraph(),
		Email:    RandomEmail(),
		Fullname: randomdata.FullName(-1),
		Nickname: RandomNickname(),
	}
}
