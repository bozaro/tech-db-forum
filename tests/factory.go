package main

import (
	"github.com/bozaro/tech-db-forum/tests/models"
	"github.com/go-openapi/strfmt"
	"github.com/ventu-io/go-shortid"
)

const ABC = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_."

var sid *shortid.Shortid

func init() {
	sid = shortid.MustNew(0, ABC, 1)
}

func RandomEmail() strfmt.Email {
	return strfmt.Email(sid.MustGenerate() + "@see")
}

func RandomNickname() string {
	return sid.MustGenerate()
}

func RandomUser() models.User {
	return models.User{
		About:    "",
		Email:    RandomEmail(),
		Fullname: "Jack Sparrow",
		Nickname: RandomNickname(),
	}
}
