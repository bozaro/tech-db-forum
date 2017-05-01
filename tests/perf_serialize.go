package tests

import (
	"github.com/tinylib/msgp/msgp"
	"io"
)

func (self *PerfData) Load(reader io.Reader) error {
	return nil
}
func (self *PerfData) Save(writer io.Writer) error {
	w := msgp.NewWriter(writer)
	// Write users list
	w.WriteInt(len(self.users))
	for _, user := range self.users {
		user.EncodeMsg(w)
	}
	// Write forum list
	w.WriteInt(len(self.forums))
	for _, forum := range self.forums {
		forum.EncodeMsg(w)
		w.WriteString(forum.User.Nickname)
	}
	// Write thread list
	w.WriteInt(len(self.threads))
	for _, thread := range self.threads {
		thread.EncodeMsg(w)
		w.WriteString(thread.Forum.Slug)
		w.WriteString(thread.Author.Nickname)
	}
	// Write posts list
	w.WriteInt(len(self.posts))
	for _, post := range self.posts {
		post.EncodeMsg(w)
		w.WriteInt32(post.Thread.ID)
		w.WriteString(post.Author.Nickname)
		w.WriteInt64(post.GetParentId())
	}
	return w.Flush()
}
