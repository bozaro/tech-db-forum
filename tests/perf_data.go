package tests

import (
	"github.com/go-openapi/strfmt"
)

import (
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

const POST_POWER = 16

type PVersion uint32
type PHash uint32

type PerfData struct {
	mutex sync.RWMutex

	Status  *PStatus
	users   []*PUser
	forums  []*PForum
	threads []*PThread
	posts   []*PPost

	lastIndex int32

	threadsByForum    map[string][]*PThread
	usersByForum      map[string][]*PUser
	postsByThreadFlat map[int32][]*PPost
	postsByThreadTree map[int32][]*PPost
	postsByThreadParentDesc map[int32][]*PPost
	userByNickname    map[string]*PUser
	forumBySlug       map[string]*PForum
	postById          map[int64]*PPost
	threadById        map[int32]*PThread
	threadBySlug      map[string]*PThread
}

type PStatus struct {
	Version PVersion
	Forum   int32
	Post    int64
	Thread  int32
	User    int32
}

type PUser struct {
	Version      PVersion
	AboutHash    PHash
	Email        strfmt.Email
	FullnameHash PHash
	Nickname     string
}

type PThread struct {
	mutex sync.RWMutex

	Version     PVersion
	ID          int32
	Slug        string
	Author      *PUser
	Forum       *PForum
	MessageHash PHash
	TitleHash   PHash
	Created     strfmt.DateTime
	Voices      map[*PUser]int32
	Votes       int32
	Posts       int32
}

type PForum struct {
	Version   PVersion
	Posts     int64
	Slug      string
	Threads   int32
	TitleHash PHash
	User      *PUser
}

type PPost struct {
	Version     PVersion
	ID          int64
	Author      *PUser
	Thread      *PThread
	Parent      *PPost
	Created     strfmt.DateTime
	IsEdited    bool
	MessageHash PHash
	Index       int32
	Path        []int32
}

func NewPerfData(config *PerfConfig) *PerfData {
	return &PerfData{
		Status:            &PStatus{},
		forums:            make([]*PForum, 0, config.ForumCount),
		users:             make([]*PUser, 0, config.UserCount),
		threads:           make([]*PThread, 0, config.ThreadCount),
		posts:             make([]*PPost, 0, config.PostCount),
		threadsByForum:    map[string][]*PThread{},
		usersByForum:      map[string][]*PUser{},
		postsByThreadTree: map[int32][]*PPost{},
		postsByThreadParentDesc: map[int32][]*PPost{},
		postsByThreadFlat: map[int32][]*PPost{},
		userByNickname:    map[string]*PUser{},
		forumBySlug:       map[string]*PForum{},
		threadBySlug:      map[string]*PThread{},
		threadById:        map[int32]*PThread{},
		postById:          map[int64]*PPost{},
	}
}

func (self *PerfData) GetForumBySlug(slug string) *PForum {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return self.forumBySlug[slug]
}

func (self *PerfData) GetUserByNickname(nickname string) *PUser {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return self.userByNickname[nickname]
}

func getRandomIndex(count int, p float64) int {
	if count == 0 {
		return -1
	}
	q := rand.Float64()
	if p > 0 && p != 1 {
		if q >= 0.5 {
			q = (q - 0.5) * 2
			q = math.Pow(q, p)
			q = q/2 + 0.5
		} else {
			q = (0.5 - q) * 2
			q = math.Pow(q, p)
			q = 0.5 - q/2
		}
	}
	return int(q * float64(count))
}

func (self *PerfData) AddForum(forum *PForum) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if _, ok := self.forumBySlug[forum.Slug]; ok {
		panic("Internal error: forum.Slug = " + forum.Slug)
	}
	self.forums = append(self.forums, forum)
	self.forumBySlug[forum.Slug] = forum
	self.usersByForum[forum.Slug] = []*PUser{}
	self.threadsByForum[forum.Slug] = []*PThread{}
	self.Status.Forum++
}

func (self *PerfData) GetForum(index int) *PForum {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	if index < 0 {
		index = getRandomIndex(len(self.forums), 1)
	}
	return self.forums[index]
}

func (self *PerfData) AddUser(user *PUser) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if _, ok := self.userByNickname[user.Nickname]; ok {
		panic("Internal error: user.Nickname = " + user.Nickname)
	}
	self.users = append(self.users, user)
	self.userByNickname[user.Nickname] = user
	self.Status.User++
}

func (self *PerfData) GetUser(index int) *PUser {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	if index < 0 {
		index = getRandomIndex(len(self.users), 1)
	}
	return self.users[index]
}

func (self *PerfData) AddThread(thread *PThread) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if thread.Slug != "" {
		if _, ok := self.threadBySlug[thread.Slug]; ok {
			panic("Internal error: thread.Slug = " + thread.Slug)
		}
		self.threadBySlug[thread.Slug] = thread
	}
	self.threads = append(self.threads, thread)
	self.threadById[thread.ID] = thread
	self.postsByThreadTree[thread.ID] = []*PPost{}
	self.postsByThreadParentDesc[thread.ID] = []*PPost{}
	self.postsByThreadFlat[thread.ID] = []*PPost{}
	self.threadsByForum[thread.Forum.Slug] = append(self.threadsByForum[thread.Forum.Slug], thread)
	self.usersByForum[thread.Forum.Slug] = append(self.usersByForum[thread.Forum.Slug], thread.Author)
	thread.Forum.Threads++
	self.Status.Thread++
}

func (self *PerfData) GetThread(index int, p float64, offset float64) *PThread {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	if index < 0 {
		length := len(self.threads)
		delta := int((offset+0.5)*float64(length)) % length
		if delta < 0 {
			delta = length - delta
		}

		index = getRandomIndex(length, p)
		index = (index + delta) % length
	}
	return self.threads[index]
}

func (self *PerfData) GetThreadById(id int32) *PThread {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return self.threadById[id]
}

func (self *PerfData) GetPostById(id int64) *PPost {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return self.postById[id]
}

func (self *PerfData) GetForumThreads(forum *PForum) []*PThread {
	/*self.mutex.RLock()
	defer self.mutex.RUnlock()

	return append([]*PThread{}, self.threadsByForum[forum.Slug]...)*/
	return self.threadsByForum[forum.Slug]
}

func (self *PerfData) GetForumThreadsByCreated(forum *PForum, since *strfmt.DateTime, desc bool, limit int) []*PThread {
	threads := self.GetForumThreads(forum)
	idx := 0
	if since != nil {
		idx = sort.Search(len(threads), func(i int) bool { return !time.Time(threads[i].Created).Before(time.Time(*since)) })
		if idx >= len(threads) {
			idx = len(threads) - 1
		}
	}
	result := make([]*PThread, 0, limit)
	if desc {
		if since == nil {
			idx = len(threads) - 1
		}
		for i := idx; i >= 0; i-- {
			if len(result) == limit {
				break
			}
			if since == nil || !time.Time(*since).Before(time.Time(threads[i].Created)) {
				result = append(result, threads[i])
			}
		}
	} else {
		for i := idx; i < len(threads); i++ {
			if len(result) == limit {
				break
			}
			if since == nil || !time.Time(threads[i].Created).Before(time.Time(*since)) {
				result = append(result, threads[i])
			}
		}
	}
	return result
}
func (self *PerfData) GetForumUsers(forum *PForum) []*PUser {
	/*self.mutex.RLock()
	defer self.mutex.RUnlock()

	return append([]*PUser{}, self.usersByForum[forum.Slug]...)*/
	return self.usersByForum[forum.Slug]
}

func (self *PerfData) GetForumUsersByNickname(forum *PForum, since *string, desc bool, limit int) []*PUser {
	users := self.GetForumUsers(forum)
	idx := 0
	lower := ""
	if since != nil {
		lower = strings.ToLower(*since)
		idx = sort.Search(len(users), func(i int) bool { return strings.ToLower(users[i].Nickname) >= lower })
		if idx >= len(users) {
			idx = len(users) - 1
		}
	}
	result := make([]*PUser, 0, limit)
	if desc {
		if since == nil {
			idx = len(users) - 1
		}
		for i := idx; i >= 0; i-- {
			if len(result) == limit {
				break
			}
			if since == nil || strings.ToLower(users[i].Nickname) < lower {
				result = append(result, users[i])
			}
		}
	} else {
		for i := idx; i < len(users); i++ {
			if len(result) == limit {
				break
			}
			if since == nil || strings.ToLower(users[i].Nickname) > lower {
				result = append(result, users[i])
			}
		}
	}
	return result
}

func (self *PerfData) GetThreadPostsFlat(thread *PThread) []*PPost {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return self.postsByThreadFlat[thread.ID]
}

func (self *PerfData) GetThreadPostsTree(thread *PThread) []*PPost {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return self.postsByThreadTree[thread.ID]
}

func (self *PerfData) GetThreadPostsParentDesc(thread *PThread) []*PPost {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	return self.postsByThreadParentDesc[thread.ID]
}

func (self *PerfData) AddPost(post *PPost) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.posts = append(self.posts, post)
	self.postById[post.ID] = post
	self.usersByForum[post.Thread.Forum.Slug] = append(self.usersByForum[post.Thread.Forum.Slug], post.Author)

	self.lastIndex++
	post.Index = self.lastIndex
	if post.Parent != nil {
		// Явное копирование массива, т.к. append не всегда ведёт себя адекватно в многопоточном окружении
		path := make([]int32, 0, len(post.Parent.Path)+1)
		path = append(path, post.Parent.Path...)
		path = append(path, post.Index)
		post.Path = path
	} else {
		post.Path = []int32{post.Index}
	}

	self.postsByThreadFlat[post.Thread.ID] = append(self.postsByThreadFlat[post.Thread.ID], post)
	self.postsByThreadTree[post.Thread.ID] = append(self.postsByThreadTree[post.Thread.ID], post)
	self.postsByThreadParentDesc[post.Thread.ID] = append(self.postsByThreadParentDesc[post.Thread.ID], post)

	post.Thread.Forum.Posts++
	post.Thread.Posts++
	self.Status.Post++
}

func (self *PerfData) Normalize() {
	for _, posts := range self.postsByThreadTree {
		sort.Sort(PPostSortTree(posts))
	}
	for _, threads := range self.threadsByForum {
		sort.Sort(PThreadByCreated(threads))
	}
	for _, posts := range self.postsByThreadParentDesc {
		sort.Sort(PPostSortParentDesc(posts))
	}
	for key, users := range self.usersByForum {
		sort.Sort(PUserByNickname(users))
		size := 0
		var last *PUser
		for _, user := range users {
			if last != user {
				size++
				last = user
			}
		}
		uniq := make([]*PUser, 0, size)
		for _, user := range users {
			if last != user {
				uniq = append(uniq, user)
				last = user
			}
		}
		self.usersByForum[key] = uniq
	}
}

func (self *PerfData) GetPost(index int, p float64) *PPost {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	if index < 0 {
		index = getRandomIndex(len(self.posts), p)
	}
	return self.posts[index]
}

func (self *PPost) GetParentId() int64 {
	if self.Parent == nil {
		return 0
	}
	return self.Parent.ID
}

func GetRandomLimit() int32 {
	return 15 + rand.Int31n(5)
}

func GetRandomDesc() *bool {
	switch rand.Intn(3) {
	case 0:
		v := false
		return &v
	case 1:
		v := true
		return &v
	default:
		return nil
	}
}

func dateTimeToInt64(value strfmt.DateTime) int64 {
	return time.Time(value).UnixNano()
}

func int64ToDateTime(value int64) strfmt.DateTime {
	return strfmt.DateTime(time.Unix(0, value))
}
