package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"sort"
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
		Name:        "thread_get_posts_notfound",
		Description: "",
		FnCheck:     CheckThreadGetPostsNotFound,
		Deps: []string{
			"thread_get_posts_simple",
		},
	})
}

func CreateTree(c *client.Forum, thread *models.Thread) []OrderedPost {
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
		{1, 3, 5, 8, 10},
		{1, 3, 5, 9},
		{1, 4},
		{11},
		{11, 17},
		{11, 17, 20},
		{11, 19},
		{12},
		{12, 16},
		{13},
		{13, 18},
	}

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
			v.path += fmt.Sprintf("%02x", i)
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
			post := RandomPost()
			if v.parent != nil {
				post.Parent = v.parent.id
			}
			posts[i] = post
		}
		posts = CreatePosts(c, posts, thread)
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

func SortPosts(posts []OrderedPost, desc bool, limitType func(OrderedPost) int, limit int) [][]*models.Post {
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
	result := [][]*models.Post{}
	page := []*models.Post{}
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

func CheckThreadGetPostsSimple(c *client.Forum, m *Modify) {
	thread := CreateThread(c, nil, nil, nil)
	tree := CreateTree(c, thread)

	// Sort order
	var sortType *string
	limitType := func(post OrderedPost) int {
		return post.idx
	}
	switch m.Int(4) {
	case 0:
		sortType = nil
	case 1:
		v := "flat"
		sortType = &v
	case 2:
		v := "tree"
		sortType = &v
		sort.Sort(PostSortTree(tree))
	case 3:
		v := "parent_tree"
		sortType = &v
		sort.Sort(PostSortTree(tree))
		limitType = func(post OrderedPost) int {
			return post.top
		}
	}

	// Desc
	var desc *bool = m.NullableBool()

	// Slug or ID
	id := m.SlugOrId(thread)

	// Check read all
	marker_filter := func(data interface{}) interface{} {
		page := data.(*models.PostPage)
		page.Marker = ""
		return page
	}
	all_posts := SortPosts(tree, desc != nil && *desc, limitType, 0)[0]
	c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(id).
		WithSort(sortType).
		WithDesc(desc).
		WithContext(Expected(200, &models.PostPage{
			RandomMarker(),
			all_posts,
		}, marker_filter)))

	// Check read by 3 records
	var limit int32 = 3
	batches := SortPosts(tree, desc != nil && *desc, limitType, int(limit))
	var marker *string = nil
	for _, batch := range batches {
		page, err := c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
			WithSlugOrID(id).
			WithSort(sortType).
			WithLimit(&limit).
			WithDesc(desc).
			WithMarker(marker).
			WithContext(Expected(200, &models.PostPage{
				RandomMarker(),
				batch,
			}, marker_filter)))
		CheckNil(err)
		marker = &page.Payload.Marker
	}

	// Check read after all
	c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(id).
		WithSort(sortType).
		WithLimit(&limit).
		WithDesc(desc).
		WithMarker(marker).
		WithContext(Expected(200, &models.PostPage{
			Marker: *marker,
			Posts:  []*models.Post{},
		}, nil)))
}

func CheckThreadGetPostsNotFound(c *client.Forum) {
	thread := RandomThread()
	_, err := c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(thread.Slug).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewThreadGetPostsNotFound(), err)

	_, err = c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(THREAD_FAKE_ID).
		WithContext(Expected(404, nil, nil)))
	CheckIsType(operations.NewThreadGetPostsNotFound(), err)
}
