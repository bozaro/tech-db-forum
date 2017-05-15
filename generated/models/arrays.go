package models

import (
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

func easyjsonUnmarshalPosts(in *jlexer.Lexer, out *Posts) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Posts, 0, 8)
			} else {
				*out = Posts{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 *Post
			if in.IsNull() {
				in.Skip()
				v1 = nil
			} else {
				if v1 == nil {
					v1 = new(Post)
				}
				(*v1).UnmarshalEasyJSON(in)
			}
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonMarshalPosts(out *jwriter.Writer, in Posts) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			if v3 == nil {
				out.RawString("null")
			} else {
				(*v3).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Posts) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonMarshalPosts(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Posts) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonMarshalPosts(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Posts) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonUnmarshalPosts(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Posts) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonUnmarshalPosts(l, v)
}
func easyjsonUnmarshalThreads(in *jlexer.Lexer, out *Threads) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Threads, 0, 8)
			} else {
				*out = Threads{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v4 *Thread
			if in.IsNull() {
				in.Skip()
				v4 = nil
			} else {
				if v4 == nil {
					v4 = new(Thread)
				}
				(*v4).UnmarshalEasyJSON(in)
			}
			*out = append(*out, v4)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonMarshalThreads(out *jwriter.Writer, in Threads) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v5, v6 := range in {
			if v5 > 0 {
				out.RawByte(',')
			}
			if v6 == nil {
				out.RawString("null")
			} else {
				(*v6).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Threads) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonMarshalThreads(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Threads) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonMarshalThreads(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Threads) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonUnmarshalThreads(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Threads) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonUnmarshalThreads(l, v)
}
func easyjsonUnmarshalUsers(in *jlexer.Lexer, out *Users) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(Users, 0, 8)
			} else {
				*out = Users{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v7 *User
			if in.IsNull() {
				in.Skip()
				v7 = nil
			} else {
				if v7 == nil {
					v7 = new(User)
				}
				(*v7).UnmarshalEasyJSON(in)
			}
			*out = append(*out, v7)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonMarshalUsers(out *jwriter.Writer, in Users) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v8, v9 := range in {
			if v8 > 0 {
				out.RawByte(',')
			}
			if v9 == nil {
				out.RawString("null")
			} else {
				(*v9).MarshalEasyJSON(out)
			}
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v Users) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonMarshalUsers(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Users) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonMarshalUsers(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Users) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonUnmarshalUsers(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Users) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonUnmarshalUsers(l, v)
}
