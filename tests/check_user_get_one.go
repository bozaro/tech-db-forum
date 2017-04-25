package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
)

func init() {
	Register(Checker{
		Name:        "user_get_one_simple",
		Description: "",
		FnCheck:     CheckUserGetOneSimple,
		Deps: []string{
			"user_create_simple",
		},
	})
	Register(Checker{
		Name:        "user_get_one_notfound",
		Description: "",
		FnCheck:     CheckUserGetOneNotFound,
		Deps: []string{
			"user_get_one_simple",
		},
	})
	Register(Checker{
		Name:        "user_get_one_nocase",
		Description: "",
		FnCheck:     Modifications(CheckUserGetOneNocase),
		Deps: []string{
			"user_get_one_simple",
		},
	})
	PerfRegister(PerfTest{
		Name:   "user_get_one_success",
		Mode:   ModeRead,
		Weight: WeightNormal,
		FnPerf: PerfUserGetOneSuccess,
	})
	PerfRegister(PerfTest{
		Name:   "user_get_one_not_found",
		Mode:   ModeRead,
		Weight: WeightRare,
		FnPerf: PerfUserGetOneNotFound,
	})
}

func CheckUserGetOneSimple(c *client.Forum) {
	user := CreateUser(c, nil)
	CheckUser(c, user)
}

func CheckUserGetOneNotFound(c *client.Forum) {
	user := RandomUser()
	_, err := c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(user.Nickname).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewUserGetOneNotFound(), err)
}

func CheckUserGetOneNocase(c *client.Forum, m *Modify) {
	user := CreateUser(c, nil)
	nickname := m.Case(user.Nickname)
	_, err := c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(nickname).
		WithContext(Expected(200, user, nil)))
	CheckNil(err)
}

func PerfUserGetOneSuccess(p *Perf) {
	user := p.data.GetUser(-1)
	version := user.Version
	result, err := p.c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(GetRandomCase(user.Nickname)).
		WithContext(Expected(200, nil, nil)))
	CheckNil(err)

	p.Validate(func(v PerfValidator) {
		payload := result.Payload
		v.CheckHash(user.AboutHash, payload.About, "Incorrect About")
		v.CheckStr(user.Email.String(), payload.Email.String(), "Incorrect Email")
		v.CheckHash(user.FullnameHash, payload.Fullname, "Incorrect Fullname")
		v.CheckStr(user.Nickname, payload.Nickname, "Incorrect Nickname")
		v.Finish(version, user.Version)
	})
}

func PerfUserGetOneNotFound(p *Perf) {
	nickname := RandomUser().Nickname
	_, err := p.c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(nickname).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewUserGetOneNotFound(), err)
}
