package tests

import (
	"fmt"
	"github.com/bozaro/tech-db-forum/generated/client"
	"github.com/bozaro/tech-db-forum/generated/client/operations"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/go-openapi/strfmt"
	"time"
)

func init() {
	Register(Checker{
		Name:        "posts_create_simple",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateSimple),
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_same_time",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateSameTime),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_empty",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateEmpty),
		Deps: []string{
			"thread_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_unicode",
		Description: "",
		FnCheck:     CheckPostCreateUnicode,
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_no_thread",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateNoThread),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_no_author",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateNoAuthor),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_with_parent",
		Description: "",
		FnCheck:     CheckPostCreateWithParent,
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_invalid_parent",
		Description: "",
		FnCheck:     Modifications(CheckPostCreateInvalidParent),
		Deps: []string{
			"posts_create_simple",
		},
	})
	Register(Checker{
		Name:        "posts_create_deep_parent",
		Description: "",
		FnCheck:     CheckPostCreateDeepParent,
		Deps: []string{
			"posts_create_simple",
		},
	})
}

func (f *Factory) CreatePosts(c *client.Forum, posts models.Posts, thread *models.Thread) models.Posts {
	if len(posts) == 0 {
		return models.Posts{}
	}
	var postsThread int32 = 0
	for _, post := range posts {
		if post.Thread != 0 {
			if postsThread == 0 {
				postsThread = post.Thread
			}
			if postsThread != post.Thread {
				panic("Invalid test data: can't create multiple posts in differ threads")
			}
		}
	}
	slug := ""
	var postsForum string
	if postsThread == 0 {
		if thread == nil {
			thread = f.CreateThread(c, nil, nil, nil)
		} else {
			slug = thread.Slug
		}
		postsThread = thread.ID
		postsForum = thread.Forum
	}
	if slug == "" {
		slug = fmt.Sprintf("%d", postsThread)
	}
	var author string

	check_forum := postsForum != ""
	if !check_forum {
		postsForum = f.RandomForum().Slug
	}

	var expected models.Posts
	if thread != nil {
		author = thread.Author
	}
	base_id := 42
	example_time := strfmt.DateTime(time.Now())
	for n, post := range posts {
		if post.Author == "" {
			if author == "" {
				author = f.CreateUser(c, nil).Nickname
			}
			post.Author = author
		}
		post.Thread = 0

		expectedPost := *post
		expectedPost.ID = int64(base_id + n)
		expectedPost.Thread = postsThread
		expectedPost.Created = &example_time
		expectedPost.Forum = postsForum
		expected = append(expected, &expectedPost)
	}

	result, err := c.Operations.PostsCreate(operations.NewPostsCreateParams().
		WithSlugOrID(slug).
		WithPosts(posts).
		WithContext(Expected(201, &expected, func(data interface{}) interface{} {
			posts := data.(*models.Posts)

			var first_time *int64
			same_time := true
			for _, post := range *posts {
				if post.Created != nil {
					nano := time.Time(*post.Created).UnixNano()
					if first_time == nil {
						first_time = &nano
					}
					if *first_time != nano {
						same_time = false
					}
				}
			}

			for n, post := range *posts {
				if post.ID != 0 {
					post.ID = int64(base_id + n)
				}
				if !check_forum {
					post.Forum = ""
				}
				if same_time && (post.Created != nil) {
					post.Created = &example_time
				}
			}
			return data
		})))

	CheckNil(err)
	return result.Payload
}

func (f *Factory) CreatePost(c *client.Forum, post *models.Post, thread *models.Thread) *models.Post {
	if post == nil {
		post = f.RandomPost()
	}
	return f.CreatePosts(c, models.Posts{post}, thread)[0]
}

func CheckPost(c *client.Forum, post *models.Post) {
	_, err := c.Operations.PostGetOne(operations.NewPostGetOneParams().
		WithID(post.ID).
		WithContext(Expected(200, &models.PostFull{
			Post: post,
		}, nil)))
	CheckNil(err)
}

func CheckPostCreateSimple(c *client.Forum, f *Factory, m *Modify) {
	thread := f.CreateThread(c, nil, nil, nil)
	if m.Bool() {
		thread.Slug = ""
	}
	f.CreatePost(c, nil, thread)
}

func CheckPostCreateSameTime(c *client.Forum, f *Factory, m *Modify) {
	thread := f.CreateThread(c, nil, nil, nil)
	if m.Bool() {
		thread.Slug = ""
	}
	f.CreatePosts(c, []*models.Post{
		f.RandomPost(),
		f.RandomPost(),
	}, thread)
}

func CheckPostCreateEmpty(c *client.Forum, f *Factory, m *Modify) {
	thread := f.CreateThread(c, nil, nil, nil)
	f.CreatePosts(c, []*models.Post{}, thread)

	_, err := c.Operations.PostsCreate(operations.NewPostsCreateParams().
		WithSlugOrID(m.SlugOrId(thread)).
		WithPosts(models.Posts{}).
		WithContext(Expected(201, &models.Posts{}, nil)))

	CheckNil(err)
}

func CheckPostCreateNoThread(c *client.Forum, f *Factory, m *Modify) {
	posts := []*models.Post{}
	if m.Bool() {
		post := f.RandomPost()
		post.Author = f.CreateUser(c, nil).Nickname
		posts = append(posts, post)
	}
	var err error
	_, err = c.Operations.PostsCreate(operations.NewPostsCreateParams().
		WithSlugOrID(THREAD_FAKE_ID).
		WithPosts(posts).
		WithContext(ExpectedError(404, "Can't find post thread by id: %s", THREAD_FAKE_ID)))
	CheckIsType(operations.NewPostsCreateNotFound(), err)

	slug := f.RandomThread().Slug
	_, err = c.Operations.PostsCreate(operations.NewPostsCreateParams().
		WithSlugOrID(slug).
		WithPosts(posts).
		WithContext(ExpectedError(404, "Can't find post thread by slug: %s", slug)))
	CheckIsType(operations.NewPostsCreateNotFound(), err)
}

func CheckPostCreateNoAuthor(c *client.Forum, f *Factory, m *Modify) {
	post := f.RandomPost()
	post.Author = f.RandomNickname()
	thread := f.CreateThread(c, nil, nil, nil)

	_, err := c.Operations.PostsCreate(operations.NewPostsCreateParams().
		WithSlugOrID(m.SlugOrId(thread)).
		WithPosts([]*models.Post{post}).
		WithContext(ExpectedError(404, "Can't find post author by nickname: %s", post.Author)))
	CheckIsType(operations.NewPostsCreateNotFound(), err)
}

func CheckPostCreateInvalidParent(c *client.Forum, f *Factory, m *Modify) {
	post := f.RandomPost()
	thread := f.CreateThread(c, nil, nil, nil)
	post.Author = thread.Author
	parentId := POST_FAKE_ID
	if m.Bool() {
		parentId = f.CreatePost(c, nil, nil).ID
	}
	post.Parent = parentId
	_, err := c.Operations.PostsCreate(operations.NewPostsCreateParams().
		WithSlugOrID(m.SlugOrId(thread)).
		WithPosts([]*models.Post{post}).
		WithContext(ExpectedError(409, "Parent post was created in another thread")))
	CheckIsType(operations.NewPostsCreateConflict(), err)
}

func CheckPostCreateWithParent(c *client.Forum, f *Factory) {
	post := f.RandomPost()
	thread := f.CreateThread(c, nil, nil, nil)
	post.Author = thread.Author
	post.Parent = f.CreatePost(c, nil, thread).ID
	f.CreatePost(c, post, thread)
}

func CheckPostCreateDeepParent(c *client.Forum, f *Factory) {
	thread := f.CreateThread(c, nil, nil, nil)
	var parent int64
	for level := 0; level < 0x10; level++ {
		post := f.RandomPost()
		post.Author = thread.Author
		post.Parent = parent
		f.CreatePost(c, post, thread)
		parent = post.ID
	}
}

func CheckPostCreateUnicode(c *client.Forum, f *Factory) {
	post := f.RandomPost()
	post.Message = "大象销售。不贵。"
	f.CreatePost(c, post, nil)
}
