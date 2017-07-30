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

type PostSortTree []OrderedPost

func (a PostSortTree) Len() int           { return len(a) }
func (a PostSortTree) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PostSortTree) Less(i, j int) bool { return a[i].path < a[j].path }

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
	return a.ComparePath(&a[i].Path, &a[j].Path) < 0
}
func (a PPostSortTree) ComparePath(p1 *[]int32, p2 *[]int32) int {
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

func SortPosts(posts []OrderedPost, desc bool, limitType func(OrderedPost) int, limit int) []models.Posts {
	if limit <= 0 {
		limit = len(posts)
	}
	sorted := posts
	// Descending order
	if desc {
		sorted = make([]OrderedPost, len(posts))
		for i, v := range posts {
			sorted[len(posts)-i-1] = v
		}
	}
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
	switch m.Int(4) {
	case 0:
		sortType = nil
	case 1:
		v := SORT_FLAT
		sortType = &v
	case 2:
		v := SORT_TREE
		sortType = &v
		sort.Sort(PostSortTree(posts_tree))
	case 3:
		v := SORT_PARENT
		sortType = &v
		sort.Sort(PostSortTree(posts_tree))
		limitType = func(post OrderedPost) int {
			return post.top
		}
	}

	// Desc
	var desc *bool = m.NullableBool()

	// Slug or ID
	id := m.SlugOrId(thread)

	// Check read all
	all_posts := SortPosts(posts_tree, desc != nil && *desc, limitType, 0)[0]
	full_size := int32(len(all_posts) + 10)
	c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(id).
		WithSort(sortType).
		WithDesc(desc).
		WithLimit(&full_size).
		WithContext(Expected(200, &all_posts, filterPostPage)))

	if limit > 0 {
		// Check read records page by page
		batches := SortPosts(posts_tree, desc != nil && *desc, limitType, int(limit))
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
	thread := p.data.GetThread(-1)
	version := thread.Version

	slugOrId := GetSlugOrId(thread.Slug, int64(thread.ID))
	limit := GetRandomLimit()
	desc := GetRandomDesc()
	order := GetRandomSort()

	part1, err := p.c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(slugOrId).
		WithLimit(&limit).
		WithSort(&order).
		WithSince(nil).
		WithDesc(desc).
		WithContext(Expected(200, nil, nil)))

	CheckNil(err)

	var last_id *int64 = nil
	if len(part1.Payload) > 0 {
		last_id = &part1.Payload[len(part1.Payload)-1].ID
	}

	part2, err := p.c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(slugOrId).
		WithLimit(&limit).
		WithSort(&order).
		WithSince(last_id).
		WithDesc(desc).
		WithContext(Expected(200, nil, nil)))

	CheckNil(err)

	p.Validate(func(v PerfValidator) {
		if v.CheckVersion(version, thread.Version) {
			// Sort
			limitType := func(post *PPost) int32 {
				return post.Index
			}
			var expected []*PPost
			switch order {
			case SORT_FLAT:
				expected = p.data.GetThreadPostsFlat(thread)
			case SORT_TREE:
				expected = p.data.GetThreadPostsTree(thread)
			case SORT_PARENT:
				expected = p.data.GetThreadPostsTree(thread)
				limitType = func(post *PPost) int32 {
					return post.Path[0]
				}
			default:
				panic("Unexpected sort type: " + order)
			}
			reverse := (desc != nil) && (*desc == true)
			// Check
			postSplit := func(full []*PPost) ([]*PPost, []*PPost) {
				count := int32(0)
				last := int32(0)
				if reverse {
					for i := len(full) - 1; i >= 0; i-- {
						item := full[i]
						if last != limitType(item) {
							if count == limit {
								return full[i+1:], full[:i+1]
							}
							count++
							last = limitType(item)
						}
					}
					return full, []*PPost{}
				} else {
					for i, item := range full {
						if last != limitType(item) {
							if count == limit {
								return full[:i], full[i:]
							}
							count++
							last = limitType(item)
						}
					}
					return full, []*PPost{}
				}
			}

			expected1, expected := postSplit(expected)
			expected2, expected := postSplit(expected)

			validate := func(posts []*PPost, actual []*models.Post) {
				v.CheckInt(len(posts), len(actual), "len()")
				for i := range posts {
					item := posts[i]
					if reverse {
						item = posts[len(posts)-i-1]
					}
					item.Validate(v, actual[i], item.Version, fmt.Sprintf("Post[%d]", i))
				}
			}
			validate(expected1, part1.Payload)
			validate(expected2, part2.Payload)

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
