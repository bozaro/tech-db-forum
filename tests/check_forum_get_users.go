package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"math/rand"
	"sort"
	"strings"
	"time"
)

func init() {
	Register(Checker{
		Name:        "forum_get_users_simple",
		Description: "",
		FnCheck:     Modifications(CheckForumGetUsersSimple),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_users_collation",
		Description: "Данный тест проверяет корректность сортировки пользователей.",
		FnCheck:     Modifications(CheckForumGetUsersCollation),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_users_notfound",
		Description: "",
		FnCheck:     Modifications(CheckForumGetUsersNotFound),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_users_empty",
		Description: "",
		FnCheck:     Modifications(CheckForumGetUsersEmpty),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "forum_get_users_vote",
		Description: "",
		FnCheck:     Modifications(CheckForumGetUsersVote),
		Deps: []string{
			"posts_create_simple",
			"thread_create_vote_simple",
		},
	})
	PerfRegister(PerfTest{
		Name:   "forum_get_users_success",
		Mode:   ModeRead,
		Weight: WeightNormal,
		FnPerf: PerfForumGetUsersSuccess,
	})
	PerfRegister(PerfTest{
		Name:   "forum_get_users_not_found",
		Mode:   ModeRead,
		Weight: WeightRare,
		FnPerf: PerfForumGetUsersNotFound,
	})
}

type PUserByNickname []*PUser

func (a PUserByNickname) Len() int      { return len(a) }
func (a PUserByNickname) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a PUserByNickname) Less(i, j int) bool {
	return strings.ToLower(a[i].Nickname) < strings.ToLower(a[j].Nickname)
}

func CheckForumGetUsersSimple(c *client.Forum, f *Factory, m *Modify) {
	forum := f.CreateForum(c, nil, nil)
	threads := []*models.Thread{}
	users := []*models.User{}
	created := time.Now()
	created.Round(time.Millisecond)

	forum_users := map[string]*models.User{}
	// Пост, который не участвует в данном форуме
	f.CreatePost(c, nil, nil)
	// Пользователи
	for i := 0; i < 8; i++ {
		user := f.CreateUser(c, nil)
		users = append(users, user)
	}
	// Ветви
	for i := 0; i < 4; i++ {
		user := users[i/2]
		thread := f.CreateThread(c, nil, forum, user)
		threads = append(threads, thread)
		forum_users[user.Nickname] = user
	}
	// Посты
	for i := 0; i < 10; i++ {
		user := users[i%(len(users)-1)+1]
		thread := threads[rand.Intn(len(threads))]
		post := f.RandomPost()
		post.Author = user.Nickname
		f.CreatePost(c, post, thread)
		forum_users[user.Nickname] = user
	}

	// Список пользователей
	all_expected := models.Users{}
	for _, user := range forum_users {
		all_expected = append(all_expected, user)
	}
	sort.Sort(UserByNickname(all_expected))

	// Desc
	desc := m.NullableBool()
	small := time.Millisecond
	if desc != nil && *desc {
		small = -small
		sort.Sort(sort.Reverse(UserByNickname(all_expected)))
	}

	slug := m.Case(forum.Slug)

	// Check read all
	c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(slug).
		WithDesc(desc).
		WithContext(Expected(200, &all_expected, nil)))

	// Check read by 4 records
	limit := int32(4)
	for n := 0; n < len(all_expected); n += int(limit) - 1 {
		m := n + int(limit)
		if m > len(all_expected) {
			m = len(all_expected)
		}
		expected := all_expected[n:m]
		var since *string = nil
		if n > 0 {
			since = &all_expected[n-1].Nickname
		}
		c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
			WithSlug(slug).
			WithLimit(&limit).
			WithDesc(desc).
			WithSince(since).
			WithContext(Expected(200, &expected, nil)))
	}

	// Check read after all
	c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(slug).
		WithLimit(&limit).
		WithDesc(desc).
		WithSince(&all_expected[len(all_expected)-1].Nickname).
		WithContext(Expected(200, &models.Threads{}, nil)))
}

func CheckForumGetUsersCollation(c *client.Forum, f *Factory, m *Modify) {
	forum := f.CreateForum(c, nil, nil)
	threads := []*models.Thread{}
	users := []*models.User{}
	created := time.Now()
	created.Round(time.Millisecond)

	forum_users := map[string]*models.User{}
	// Пост, который не участвует в данном форуме
	f.CreatePost(c, nil, nil)
	// Суффиксы пользователей
	prefix := nick_id.Generate() + "."
	suffixes := []string{
		"joe",
		"_joe",
		".joe",
		"Jill",
		"bill",
		"Bob",
		"Zod",
	}
	// Пользователи
	for _, suffix := range suffixes {
		user := f.RandomUser()
		user.Nickname = prefix + suffix
		user = f.CreateUser(c, user)
		users = append(users, user)
	}
	// Ветви
	for i := 0; i < 4; i++ {
		user := users[i/2]
		thread := f.CreateThread(c, nil, forum, user)
		threads = append(threads, thread)
		forum_users[user.Nickname] = user
	}
	// Посты
	for i := 0; i < 10; i++ {
		user := users[i%(len(users)-1)+1]
		thread := threads[rand.Intn(len(threads))]
		post := f.RandomPost()
		post.Author = user.Nickname
		f.CreatePost(c, post, thread)
		forum_users[user.Nickname] = user
	}

	// Список пользователей
	all_expected := models.Users{}
	for _, user := range forum_users {
		all_expected = append(all_expected, user)
	}
	sort.Sort(UserByNickname(all_expected))

	// Desc
	desc := m.NullableBool()
	small := time.Millisecond
	if desc != nil && *desc {
		small = -small
		sort.Sort(sort.Reverse(UserByNickname(all_expected)))
	}

	// Check read all
	c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(forum.Slug).
		WithDesc(desc).
		WithContext(Expected(200, &all_expected, nil)))

	// Check read by 4 records
	limit := int32(4)
	for n := 0; n < len(all_expected); n += int(limit) - 1 {
		m := n + int(limit)
		if m > len(all_expected) {
			m = len(all_expected)
		}
		expected := all_expected[n:m]
		var since *string = nil
		if n > 0 {
			since = &all_expected[n-1].Nickname
		}
		c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
			WithSlug(forum.Slug).
			WithLimit(&limit).
			WithDesc(desc).
			WithSince(since).
			WithContext(Expected(200, &expected, nil)))
	}

	// Check read after all
	c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(forum.Slug).
		WithLimit(&limit).
		WithDesc(desc).
		WithSince(&all_expected[len(all_expected)-1].Nickname).
		WithContext(Expected(200, &models.Threads{}, nil)))
}

func CheckForumGetUsersEmpty(c *client.Forum, f *Factory, m *Modify) {
	var limit *int32
	var since *string
	var desc *bool

	forum := f.CreateForum(c, nil, nil)
	// Limit
	if m.Bool() {
		v := int32(10)
		limit = &v
	}
	// Since
	if m.Bool() {
		v := f.RandomNickname()
		since = &v
	}
	// Desc
	switch m.Int(3) {
	case 1:
		v := true
		desc = &v
	case 2:
		v := false
		desc = &v
	}

	// Check
	c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(forum.Slug).
		WithLimit(limit).
		WithSince(since).
		WithDesc(desc).
		WithContext(Expected(200, &models.Users{}, nil)))
}

func CheckForumGetUsersVote(c *client.Forum, f *Factory, m *Modify) {
	var limit *int32
	var since *string
	var desc *bool

	forum := f.CreateForum(c, nil, nil)

	user := f.CreateUser(c, nil)
	author := f.CreateUser(c, nil)
	thread := f.CreateThread(c, nil, forum, author)
	var err error
	// Like user1
	thread.Votes = 1
	_, err = c.Operations.ThreadVote(operations.NewThreadVoteParams().
		WithSlugOrID(m.SlugOrId(thread)).
		WithVote(&models.Vote{
			Nickname: user.Nickname,
			Voice:    1,
		}).
		WithContext(Expected(200, thread, filterThread)))
	CheckNil(err)

	thread.Votes = 1
	CheckThread(c, thread)

	// Check
	c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(forum.Slug).
		WithLimit(limit).
		WithSince(since).
		WithDesc(desc).
		WithContext(Expected(200, &models.Users{author}, nil)))
}

func CheckForumGetUsersNotFound(c *client.Forum, f *Factory, m *Modify) {
	var limit *int32
	var since *string
	var desc *bool

	forum := f.RandomForum()
	// Limit
	if m.Bool() {
		v := int32(10)
		limit = &v
	}
	// Since
	if m.Bool() {
		v := f.RandomNickname()
		since = &v
	}
	// Desc
	switch m.Int(3) {
	case 1:
		v := true
		desc = &v
	case 2:
		v := false
		desc = &v
	}

	// Check
	_, err := c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(forum.Slug).
		WithLimit(limit).
		WithSince(since).
		WithDesc(desc).
		WithContext(ExpectedError(404, "Can't find forum by slug: %s", forum.Slug)))
	CheckIsType(operations.NewForumGetUsersNotFound(), err)
}

func PerfForumGetUsersSuccess(p *Perf, f *Factory) {
	forum := p.data.GetForum(-1)
	version := forum.Version

	slug := forum.Slug
	limit := GetRandomLimit()
	var since *string
	if rand.Int()&1 == 0 {
		nick := GetRandomCase(p.data.GetUser(-1).Nickname)
		since = &nick
	}
	desc := GetRandomDesc()
	s := p.Session()
	result, err := p.c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(GetRandomCase(slug)).
		WithLimit(&limit).
		WithSince(since).
		WithDesc(desc).
		WithContext(s.Expected(200)))

	CheckNil(err)

	s.Validate(func(v PerfValidator) {
		if v.CheckVersion(version, forum.Version) {
			expected := p.data.GetForumUsersByNickname(forum, since, (desc != nil) && *desc, int(limit))
			// Check
			if len(expected) > int(limit) {
				expected = expected[0:limit]
			}

			payload := result.Payload
			v.CheckInt(len(expected), len(payload), "len()")
			for i, item := range expected {
				item.Validate(v, payload[i], item.Version)
			}
			v.Finish(version, forum.Version)
		}
	})
}

func PerfForumGetUsersNotFound(p *Perf, f *Factory) {
	slug := f.RandomSlug()
	limit := GetRandomLimit()
	var since *string
	if rand.Int()&1 == 0 {
		nick := GetRandomCase(p.data.GetUser(-1).Nickname)
		since = &nick
	}
	desc := GetRandomDesc()
	_, err := p.c.Operations.ForumGetUsers(operations.NewForumGetUsersParams().
		WithSlug(slug).
		WithLimit(&limit).
		WithSince(since).
		WithDesc(desc).
		WithContext(ExpectedError(404, "Can't find forum by slug: %s", slug)))
	CheckIsType(operations.NewForumGetUsersNotFound(), err)
}
