package tests

import (
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"sort"
)

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

func CreateTree(c *client.Forum, thread *models.Thread) []*models.Post {
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
		id     int64
	}
	nodes := map[int]*node{}
	keys := []int{}
	for _, t := range tree {
		v := node{}
		k := t[len(t)-1]
		if len(t) > 1 {
			v.parent = nodes[t[len(t)-2]]
		}
		keys = append(keys, k)
		nodes[k] = &v
	}
	sort.Ints(keys)
	result := []*models.Post{}
	for _, k := range keys {
		v := nodes[k]
		post := RandomPost()
		if v.parent != nil {
			post.Parent = v.parent.id
		}
		post = CreatePost(c, post, thread)
		v.id = post.ID
		result = append(result, post)
	}
	return result
}

func SortPosts(posts []*models.Post, desc bool, limit int) [][]*models.Post {
	if limit <= 0 {
		limit = len(posts)
	}
	sorted := posts
	if desc {
		sorted = make([]*models.Post, len(posts))
		for i, v := range posts {
			sorted[len(posts)-i-1] = v
		}
	}
	return [][]*models.Post{sorted}
}

func CheckThreadGetPostsSimple(c *client.Forum, m *Modify) {
	thread := CreateThread(c, nil, nil, nil)
	tree := CreateTree(c, thread)

	// Desc
	var desc *bool = m.NullableBool()

	// Slug or ID
	id := m.SlugOrId(thread)

	// Check read all
	all_posts := SortPosts(tree, desc != nil && *desc, 0)[0]
	c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(id).
		WithDesc(desc).
		WithContext(Expected(200, &all_posts, nil)))

	// Check read by 3 records
	/* todo: #14
	var limit int32 = 3
	batches := SortPosts(tree, desc != nil && *desc, int(limit))
	var since *strfmt.DateTime = nil
	for _, batch := range batches {
		c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
			WithSlugOrID(id).
			WithLimit(&limit).
			WithDesc(desc).
			WithSince(since).
			WithContext(Expected(200, &batch, nil)))
		//since = &threads[m-1].Created
	}

	// Check read after all
	c.Operations.ThreadGetPosts(operations.NewThreadGetPostsParams().
		WithSlugOrID(id).
		WithLimit(&limit).
		WithDesc(desc).
		WithSince(since).
		WithContext(Expected(200, &[]models.Thread{}, nil)))
	*/
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
