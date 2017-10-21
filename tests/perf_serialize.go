package tests

import (
	"errors"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/tinylib/msgp/msgp"
	"io"
)

const (
	SERIALIZE_VERSION int = 1
)

func LoadPerfData(reader io.Reader) (*PerfData, error) {
	var err error
	r := msgp.NewReader(reader)
	if project, err := r.ReadString(); err == nil {
		if project != Project {
			return nil, errors.New("Unsupported state file format")
		}
	} else {
		return nil, err
	}
	if ver, err := r.ReadInt(); err == nil {
		if ver != SERIALIZE_VERSION {
			return nil, errors.New(fmt.Sprintf("Unsupported state file version (expected: %d, actual: %d)", SERIALIZE_VERSION, ver))
		}
	} else {
		return nil, err
	}

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
		if user.Nickname, err = r.ReadString(); err != nil {
			return nil, err
		}
		if email, err := r.ReadString(); err == nil {
			user.Email = strfmt.Email(email)
		} else {
			return nil, err
		}
		if err = user.AboutHash.DecodeMsg(r); err != nil {
			return nil, err
		}
		if err = user.FullnameHash.DecodeMsg(r); err != nil {
			return nil, err
		}
		d.AddUser(&user)
	}
	// Read forum list
	for i := 0; i < c.ForumCount; i++ {
		forum := PForum{}
		if forum.Slug, err = r.ReadString(); err != nil {
			return nil, err
		}
		if err = forum.TitleHash.DecodeMsg(r); err != nil {
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
		thread := PThread{
			Voices: map[*PUser]int32{},
		}

		if thread.ID, err = r.ReadInt32(); err != nil {
			return nil, err
		}
		if thread.Slug, err = r.ReadString(); err != nil {
			return nil, err
		}
		if err = thread.MessageHash.DecodeMsg(r); err != nil {
			return nil, err
		}
		if err = thread.TitleHash.DecodeMsg(r); err != nil {
			return nil, err
		}
		if created, err := r.ReadInt64(); err == nil {
			thread.Created = int64ToDateTime(created)
		} else {
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
		// Read votes
		if voice_count, err := r.ReadInt(); err == nil {
			for j := 0; j < voice_count; j++ {
				var user *PUser
				if nickname, err := r.ReadString(); err == nil {
					user = d.GetUserByNickname(nickname)
				} else {
					return nil, err
				}
				if voice, err := r.ReadInt32(); err == nil {
					thread.Voices[user] = voice
					thread.Votes += voice
				} else {
					return nil, err
				}
			}
		} else {
			return nil, err
		}

		d.AddThread(&thread)
	}
	// Read posts list
	for i := 0; i < c.PostCount; i++ {
		post := PPost{}

		if post.ID, err = r.ReadInt64(); err != nil {
			return nil, err
		}
		if created, err := r.ReadInt64(); err == nil {
			post.Created = int64ToDateTime(created)
		} else {
			return nil, err
		}
		if post.IsEdited, err = r.ReadBool(); err != nil {
			return nil, err
		}
		if err = post.MessageHash.DecodeMsg(r); err != nil {
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
	d.Normalize()
	return d, nil
}

func (self *PerfData) Save(writer io.Writer) error {
	w := msgp.NewWriter(writer)
	w.WriteString(Project)
	w.WriteInt(SERIALIZE_VERSION)

	w.WriteInt(len(self.users))
	w.WriteInt(len(self.forums))
	w.WriteInt(len(self.threads))
	w.WriteInt(len(self.posts))
	// Write users list
	for _, user := range self.users {
		w.WriteString(user.Nickname)
		w.WriteString(user.Email.String())
		user.AboutHash.EncodeMsg(w)
		user.FullnameHash.EncodeMsg(w)
	}
	// Write forum list
	for _, forum := range self.forums {
		w.WriteString(forum.Slug)
		forum.TitleHash.EncodeMsg(w)
		w.WriteString(forum.User.Nickname)
	}
	// Write thread list
	for _, thread := range self.threads {
		w.WriteInt32(thread.ID)
		w.WriteString(thread.Slug)
		thread.MessageHash.EncodeMsg(w)
		thread.TitleHash.EncodeMsg(w)
		w.WriteInt64(dateTimeToInt64(thread.Created))

		w.WriteString(thread.Forum.Slug)
		w.WriteString(thread.Author.Nickname)
		// Write votes
		w.WriteInt(len(thread.Voices))
		for user, voice := range thread.Voices {
			w.WriteString(user.Nickname)
			w.WriteInt32(voice)
		}
	}
	// Write posts list
	for _, post := range self.posts {
		w.WriteInt64(post.ID)
		w.WriteInt64(dateTimeToInt64(post.Created))
		w.WriteBool(post.IsEdited)
		post.MessageHash.EncodeMsg(w)

		w.WriteInt32(post.Thread.ID)
		w.WriteString(post.Author.Nickname)
		w.WriteInt64(post.GetParentId())
	}
	return w.Flush()
}

func (z *PHash) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zajw uint32
		zajw, err = dc.ReadUint32()
		(*z) = PHash(zajw)
	}
	if err != nil {
		return
	}
	return
}

func (z PHash) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteUint32(uint32(z))
	if err != nil {
		return
	}
	return
}
