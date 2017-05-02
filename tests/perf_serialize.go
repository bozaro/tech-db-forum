package tests

import (
	"github.com/tinylib/msgp/msgp"
	"io"
)

func LoadPerfData(reader io.Reader) (*PerfData, error) {
	var err error
	r := msgp.NewReader(reader)
	c := NewPerfConfig()
	if c.UserCount, err = r.ReadInt(); err != nil {
		return nil, err
	}
	if c.ForumCount, err = r.ReadInt(); err != nil {
		return nil, err
	}
	if c.ThreadCount, err = r.ReadInt(); err != nil {
		return nil, err
	}
	if c.PostCount, err = r.ReadInt(); err != nil {
		return nil, err
	}

	d := NewPerfData(c)
	// Read users list
	for i := 0; i < c.UserCount; i++ {
		user := PUser{}
		if err = user.DecodeMsg(r); err != nil {
			return nil, err
		}
		d.AddUser(&user)
	}
	// Read forum list
	for i := 0; i < c.ForumCount; i++ {
		forum := PForum{}
		if err = forum.DecodeMsg(r); err != nil {
			return nil, err
		}
		if nickname, err := r.ReadString(); err == nil {
			forum.User = d.GetUserByNickname(nickname)
		} else {
			return nil, err
		}
		d.AddForum(&forum)
	}
	// Read thread list
	for i := 0; i < c.ThreadCount; i++ {
		thread := PThread{}
		if err = thread.DecodeMsg(r); err != nil {
			return nil, err
		}
		if slug, err := r.ReadString(); err == nil {
			thread.Forum = d.GetForumBySlug(slug)
		} else {
			return nil, err
		}
		if nickname, err := r.ReadString(); err == nil {
			thread.Author = d.GetUserByNickname(nickname)
		} else {
			return nil, err
		}
		d.AddThread(&thread)
	}
	// Read posts list
	for i := 0; i < c.PostCount; i++ {
		post := PPost{}
		if err = post.DecodeMsg(r); err != nil {
			return nil, err
		}
		if thread, err := r.ReadInt32(); err == nil {
			post.Thread = d.GetThreadById(thread)
		} else {
			return nil, err
		}
		if nickname, err := r.ReadString(); err == nil {
			post.Author = d.GetUserByNickname(nickname)
		} else {
			return nil, err
		}
		if parent, err := r.ReadInt64(); err == nil {
			if parent != 0 {
				post.Parent = d.GetPostById(parent)
			}
		} else {
			return nil, err
		}
		d.AddPost(&post)
	}
	return d, nil
}

func (self *PerfData) Save(writer io.Writer) error {
	w := msgp.NewWriter(writer)

	w.WriteInt(len(self.users))
	w.WriteInt(len(self.forums))
	w.WriteInt(len(self.threads))
	w.WriteInt(len(self.posts))
	// Write users list
	for _, user := range self.users {
		user.EncodeMsg(w)
	}
	// Write forum list
	for _, forum := range self.forums {
		forum.EncodeMsg(w)
		w.WriteString(forum.User.Nickname)
	}
	// Write thread list
	for _, thread := range self.threads {
		thread.EncodeMsg(w)
		w.WriteString(thread.Forum.Slug)
		w.WriteString(thread.Author.Nickname)
	}
	// Write posts list
	for _, post := range self.posts {
		post.EncodeMsg(w)
		w.WriteInt32(post.Thread.ID)
		w.WriteString(post.Author.Nickname)
		w.WriteInt64(post.GetParentId())
	}
	return w.Flush()
}
