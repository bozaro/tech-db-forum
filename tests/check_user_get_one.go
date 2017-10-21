package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
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

func CheckUserGetOneSimple(c *client.Forum, f *Factory) {
	user := f.CreateUser(c, nil)
	CheckUser(c, user)
}

func CheckUserGetOneNotFound(c *client.Forum, f *Factory) {
	user := f.RandomUser()
	_, err := c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(user.Nickname).
		WithContext(ExpectedError(404, "Can't find user by nickname: %s", user.Nickname)))
	CheckIsType(operations.NewUserGetOneNotFound(), err)
}

func CheckUserGetOneNocase(c *client.Forum, f *Factory, m *Modify) {
	user := f.CreateUser(c, nil)
	nickname := m.Case(user.Nickname)
	_, err := c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(nickname).
		WithContext(Expected(200, user, nil)))
	CheckNil(err)
}

func (self *PUser) Validate(v PerfValidator, user *models.User, version PVersion) {
	v.CheckHash(self.AboutHash, user.About, "About")
	v.CheckStr(self.Email.String(), user.Email.String(), "Email")
	v.CheckHash(self.FullnameHash, user.Fullname, "Fullname")
	v.CheckStr(self.Nickname, user.Nickname, "Nickname")
	v.Finish(version, self.Version)
}

func PerfUserGetOneSuccess(p *Perf, f *Factory) {
	user := p.data.GetUser(-1)
	version := user.Version
	s := p.Session()
	result, err := p.c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(GetRandomCase(user.Nickname)).
		WithContext(s.Expected(200)))
	CheckNil(err)

	s.Validate(func(v PerfValidator) {
		user.Validate(v, result.Payload, version)
	})
}

func PerfUserGetOneNotFound(p *Perf, f *Factory) {
	nickname := f.RandomNickname()
	_, err := p.c.Operations.UserGetOne(operations.NewUserGetOneParams().
		WithNickname(nickname).
		WithContext(ExpectedError(404, "Can't find user by nickname: %s", nickname)))
	CheckIsType(operations.NewUserGetOneNotFound(), err)
}
