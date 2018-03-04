package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/go-openapi/strfmt"
	"math/rand"
	"sort"
	"time"
)

const (
	SORT_FLAT   = "flat"
	SORT_TREE   = "tree"
	SORT_PARENT = "parent_tree"
)

type OrderedPost struct {
	idx  int
	top  int
	path string
	post *models.Post
}

type PostSortFlat struct {
	posts []OrderedPost
	desc  bool
}

func (a PostSortFlat) Len() int           { return len(a.posts) }
func (a PostSortFlat) Swap(i, j int)      { a.posts[i], a.posts[j] = a.posts[j], a.posts[i] }
func (a PostSortFlat) Less(i, j int) bool { return a.posts[i].idx < a.posts[j].idx != a.desc }

type PostSortTree struct {
	posts    []OrderedPost
	top_desc bool
	all_desc bool
}

func (a PostSortTree) Len() int      { return len(a.posts) }
func (a PostSortTree) Swap(i, j int) { a.posts[i], a.posts[j] = a.posts[j], a.posts[i] }
func (a PostSortTree) Less(i, j int) bool {
	if a.posts[i].top != a.posts[j].top {
		return a.posts[i].top < a.posts[j].top != (a.top_desc || a.all_desc)
	}
	return a.posts[i].path < a.posts[j].path != a.all_desc
}

type PPostSortFlat []*PPost

func (a PPostSortFlat) Len() int      { return len(a) }
func (a PPostSortFlat) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a PPostSortFlat) Less(i, j int) bool {
	return a[i].Index < a[j].Index
}

type PPostSortTree []*PPost

func (a PPostSortTree) Len() int      { return len(a) }
func (a PPostSortTree) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a PPostSortTree) Less(i, j int) bool {
	return comparePath(&a[i].Path, &a[j].Path) < 0
}

type PPostSortParentDesc []*PPost

func (a PPostSortParentDesc) Len() int      { return len(a) }
func (a PPostSortParentDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a PPostSortParentDesc) Less(i, j int) bool {
	iTop := a[i].Path[0]
	jTop := a[j].Path[0]
	if iTop != jTop {
		return iTop > jTop
	}
	return comparePath(&a[i].Path, &a[j].Path) < 0
}
func comparePath(p1 *[]int32, p2 *[]int32) int {
	l1 := len(*p1)
	l2 := len(*p2)
	for i := 0; (i < l1) && (i < l2); i++ {
		v1 := (*p1)[i]
		v2 := (*p2)[i]
		if v1 < v2 {
			return -1
		}
		if v1 > v2 {
			return 1
		}
	}
	if l1 < l2 {
		return -1
	}
	if l1 > l2 {
		return 1
	}
	return 0
}

func init() {
	Register(Checker{
		Name:        "thread_get_posts_simple",
		Description: "",
		FnCheck:     Modifications(CheckThreadGetPostsSimple),
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "thread_get_posts_same_time",
		Description: "",
		FnCheck:     Modifications(CheckThreadGetPostsSameTime),
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "thread_get_posts_notfound",
		Description: "",
		FnCheck:     CheckThreadGetPostsNotFound,
		Deps: []string{
			"thread_get_posts_simple",
		},
	})
	PerfRegister(PerfTest{
		Name:   "thread_get_posts_success",
		Mode:   ModeRead,
		Weight: WeightNormal,
		FnPerf: PerfThreadGetPostsSuccess,
	})
	PerfRegister(PerfTest{
		Name:   "thread_get_posts_not_found",
		Mode:   ModeRead,
		Weight: WeightRare,
		FnPerf: PerfThreadGetPostsNotFound,
	})
}

func (f *Factory) CreateTree(c *client.Forum, thread *models.Thread, tree [][]int) []OrderedPost {
	type node struct {
		parent *node
		path   string
		id     int64
		top    int
	}
	nodes := map[int]*node{}
	keys := []int{}
	for _, t := range tree {
		v := node{top: t[0]}
		k := t[len(t)-1]
		if len(t) > 1 {
			v.parent = nodes[t[len(t)-2]]
		}
		for _, i := range t {
			v.path += fmt.Sprintf("/%04x", i)
		}
		keys = append(keys, k)
		nodes[k] = &v
	}

	sort.Ints(keys)
	result := []OrderedPost{}
	batch := []*node{}

	flushPosts := func() {
		if len(batch) == 0 {
			panic("Internal test error")
		}
		posts := make([]*models.Post, len(batch))
		for i, v := range batch {
			post := f.RandomPost()
			if v.parent != nil {
				post.Parent = v.parent.id
			}
			posts[i] = post
		}
		posts = f.CreatePosts(c, posts, thread)
		for i, v := range batch {
			post := posts[i]
			v.id = post.ID
			result = append(result, OrderedPost{
				idx:  len(result),
				top:  v.top,
				path: v.path,
				post: post,
			})
		}
		batch = []*node{}
	}

	for _, k := range keys {
		v := nodes[k]
		if v.parent != nil && v.parent.id == 0 {
			flushPosts()
		}
		batch = append(batch, v)
	}
	if len(batch) > 0 {
		flushPosts()
	}
	return result
}

func PagePosts(posts []OrderedPost, limitType func(OrderedPost) int, limit int) []models.Posts {
	if limit <= 0 {
		limit = len(posts)
	}
	sorted := posts
	// Pagination
	result := []models.Posts{}
	page := models.Posts{}
	last := -1
	size := 0
	for _, post := range sorted {
		if last != limitType(post) {
			last = limitType(post)
			if size == limit {
				result = append(result, page)
				page = []*models.Post{}
				size = 0
			}
			size++
		}
		page = append(page, post.post)

	}
	if len(page) > 0 {
		result = append(result, page)
	}
	return result
}

func CheckThreadGetPostsSimple(c *client.Forum, f *Factory, m *Modify) {
	tree := [][]int{
		{1},
		{1, 2},
		{1, 2, 14},
		{1, 2, 15},
		{1, 3},
		{1, 3, 5},
		{1, 3, 5, 6},
		{1, 3, 5, 6, 7},
		{1, 3, 5, 8},
		{13},
		{1, 3, 5, 8, 10},
		{1, 3, 5, 9},
		{13, 18},
		{12},
		{12, 16},
		{11},
		{11, 17},
		{11, 19},
		{11, 17, 20},
		{1, 4},
	}
	CheckThreadGetPosts(c, f, m, tree, 3)
}

func CheckThreadGetPostsSameTime(c *client.Forum, f *Factory, m *Modify) {
	tree := [][]int{}
	id := 0
	top := []int{}
	for i := 0; i < 5; i++ {
		id++
		tree = append(tree, []int{id})
		top = append(top, id)
	}
	for i := 0; i < len(top)*10; i++ {
		tid := top[rand.Intn(len(top))]
		id++
		tree = append(tree, []int{tid, id})
	}

	CheckThreadGetPosts(c, f, m, tree, -1)
}

func CheckThreadGetPosts(c *client.Forum, f *Factory, m *Modify, tree [][]int, limit int32) {
	thread := f.CreateThread(c, nil, nil, nil)
	posts_tree := f.CreateTree(c, thread, tree)

	// Sort order
	var sortType *string
	limitType := func(post OrderedPost) int {
		return post.idx
	}
	sortFunc := func(posts []OrderedPost, desc bool) {
		sort.Sort(PostSortFlat{posts, desc})
	}
	switch m.Int(4) {
	case 0:
		sortType = nil
	case 1:
		v := SORT_FLAT
		sortType = &v
	case 2:
		v := SORT_TREE
		sortType = &v
		sortFunc = func(posts []OrderedPost, desc bool) {
			sort.Sort(PostSortTree{posts, false, desc})
		}
	case 3:
		v := SORT_PARENT
		sortType = &v
		sortFunc = func(posts []OrderedPost, desc bool) {
			sort.Sort(PostSortTree{posts, desc, false})
		}
		limitType = func(post OrderedPost) int {
			return post.top
		}
	}

	// Desc
	var desc *bool = m.NullableBool()

	// Slug or ID
	id := m.SlugOrId(thread)

	sortFunc(posts_tree, desc != nil && *desc)

	// Check read all
	all_posts := PagePosts(posts_tree, limitType, 0)[0]
	full_size := int32(len(all_posts) + 10)
	c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(id).
		WithSort(sortType).
		WithDesc(desc).
		WithLimit(&full_size).
		WithContext(Expected(200, &all_posts, filterPostPage)))

	if limit > 0 {
		// Check read records page by page
		batches := PagePosts(posts_tree, limitType, int(limit))
		var lastId *int64 = nil
		for _, batch := range batches {
			_, err := c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
				WithSlugOrID(id).
				WithSort(sortType).
				WithLimit(&limit).
				WithDesc(desc).
				WithSince(lastId).
				WithContext(Expected(200, &batch, filterPostPage)))
			CheckNil(err)
			lastId = &batch[len(batch)-1].ID
		}

		// Check read after all
		c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
			WithSlugOrID(id).
			WithSort(sortType).
			WithLimit(&limit).
			WithDesc(desc).
			WithSince(lastId).
			WithContext(Expected(200, &models.Posts{}, filterPostPage)))
	}
}

func filterPostPage(data interface{}) interface{} {
	page := data.(*models.Posts)
	for _, post := range *page {
		if post.Created != nil {
			created := strfmt.DateTime(time.Time(*post.Created).UTC())
			post.Created = &created
		}
	}
	return page
}

func CheckThreadGetPostsNotFound(c *client.Forum, f *Factory) {
	thread := f.RandomThread()
	_, err := c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(thread.Slug).
		WithContext(ExpectedError(404, "Can't find thread by slug: %s", thread.Slug)))
	CheckIsType(operations.NewThreadGetPostsNotFound(), err)

	_, err = c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(THREAD_FAKE_ID).
		WithContext(ExpectedError(404, "Can't find forum by id: %d", THREAD_FAKE_ID)))
	CheckIsType(operations.NewThreadGetPostsNotFound(), err)
}

func PerfThreadGetPostsSuccess(p *Perf, f *Factory) {
	thread := p.data.GetThread(-1, POST_POWER, 0.5)
	version := thread.Version

	slugOrId := GetSlugOrId(thread.Slug, int64(thread.ID))
	limit := GetRandomLimit()
	desc :=  GetRandomDesc()
	order := GetRandomSort()

	// Sort
	limitType := func(post *PPost) int32 {
		return post.Index
	}
	var expected []*PPost
	reverse := (desc != nil) && (*desc == true)
	switch order {
	case SORT_FLAT:
		expected = p.data.GetThreadPostsFlat(thread)
	case SORT_TREE:
		expected = p.data.GetThreadPostsTree(thread)
	case SORT_PARENT:
		expected = p.data.GetThreadPostsTree(thread)
		if reverse {
			expected = p.data.GetThreadPostsParentDesc(thread)
			reverse = false
		}
		limitType = func(post *PPost) int32 {
			return post.Path[0]
		}
	default:
		panic("Unexpected sort type: " + order)
	}
	var last_id *int64 = nil
	index := -1
	if reverse {
		index = len(expected)
	}
	if rand.Int()&1 == 0 {
		if len(expected) > 0 {
			rnd := rand.Intn(len(expected))
			a := limitType(expected[rnd])
			if reverse {
				for index = rnd + 1; index < len(expected); index++ {
					item := expected[index]
					if a != limitType(item) {
						last_id = &expected[index].ID
						break
					}
				}
			} else {
				for index = rnd - 1; index >= 0; index-- {
					item := expected[index]
					if a != limitType(item) {
						last_id = &expected[index].ID
						break
					}
				}
			}
		}
	}

	s := p.Session()
	posts, err := p.c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(slugOrId).
		WithLimit(&limit).
		WithSort(&order).
		WithSince(last_id).
		WithDesc(desc).
		WithContext(s.Expected(200)))

	CheckNil(err)

	s.Validate(func(v PerfValidator) {
		if v.CheckVersion(version, thread.Version) {

			// Check
			result := make([]*PPost, 0, int(limit))
			count := int32(0)
			last := int32(0)
			if reverse {
				for i := index - 1; i >= 0; i-- {
					item := expected[i]
					if last != limitType(item) {
						if count == limit {
							break
						}
						count++
						last = limitType(item)
					}
					result = append(result, item)
				}
			} else {
				for i := index + 1; i < len(expected); i++ {
					item := expected[i]
					if last != limitType(item) {
						if count == limit {
							break
						}
						count++
						last = limitType(item)
					}
					result = append(result, item)
				}
			}

			v.CheckInt(len(result), len(posts.Payload), "len()")
			for i, item := range result {
				item.Validate(v, posts.Payload[i], item.Version, fmt.Sprintf("Post[%d]", i))
			}

			v.Finish(version, thread.Version)
		}
	})
}

func PerfThreadGetPostsNotFound(p *Perf, f *Factory) {
	slug := f.RandomSlug()
	var id int32
	for {
		id = rand.Int31n(100000000)
		if p.data.GetThreadById(id) == nil {
			break
		}
	}
	slugOrId := GetSlugOrId(slug, int64(id))

	limit := GetRandomLimit()
	order := GetRandomSort()
	desc := GetRandomDesc()
	_, err := p.c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(slugOrId).
		WithLimit(&limit).
		WithSort(&order).
		WithDesc(desc).
		WithContext(ExpectedError(404, "Can't find thread by slug or id: %s", slugOrId)))
	CheckIsType(operations.NewThreadGetPostsNotFound(), err)
}

func GetRandomSort() string {
	switch rand.Intn(3) {
	case 0:
		return SORT_FLAT
	case 1:
		return SORT_TREE
	case 2:
		return SORT_PARENT
	}
	panic("Invalid internal state")
}
