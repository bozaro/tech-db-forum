package tests

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/go-openapi/strfmt"
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *PForum) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zbzg uint32
	zbzg, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zbzg > 0 {
		zbzg--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "slug":
			z.Slug, err = dc.ReadString()
			if err != nil {
				return
			}
		case "title":
			err = dc.ReadExactBytes(z.TitleHash[:])
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *PForum) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "slug"
	err = en.Append(0x82, 0xa4, 0x73, 0x6c, 0x75, 0x67)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Slug)
	if err != nil {
		return
	}
	// write "title"
	err = en.Append(0xa5, 0x74, 0x69, 0x74, 0x6c, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.TitleHash[:])
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *PForum) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "slug"
	o = append(o, 0x82, 0xa4, 0x73, 0x6c, 0x75, 0x67)
	o = msgp.AppendString(o, z.Slug)
	// string "title"
	o = append(o, 0xa5, 0x74, 0x69, 0x74, 0x6c, 0x65)
	o = msgp.AppendBytes(o, z.TitleHash[:])
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PForum) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zbai uint32
	zbai, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zbai > 0 {
		zbai--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "slug":
			z.Slug, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "title":
			bts, err = msgp.ReadExactBytes(bts, z.TitleHash[:])
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *PForum) Msgsize() (s int) {
	s = 1 + 5 + msgp.StringPrefixSize + len(z.Slug) + 6 + msgp.ArrayHeaderSize + (16 * (msgp.ByteSize))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PHash) DecodeMsg(dc *msgp.Reader) (err error) {
	err = dc.ReadExactBytes(z[:])
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *PHash) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteBytes(z[:])
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *PHash) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendBytes(o, z[:])
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PHash) UnmarshalMsg(bts []byte) (o []byte, err error) {
	bts, err = msgp.ReadExactBytes(bts, z[:])
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *PHash) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize + (16 * (msgp.ByteSize))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PPost) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zwht uint32
	zwht, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zwht > 0 {
		zwht--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "id":
			z.ID, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "created":
			{
				var zhct string
				zhct, err = dc.ReadString()
				z.Created = parseDateTime(zhct)
			}
			if err != nil {
				return
			}
		case "edited":
			z.IsEdited, err = dc.ReadBool()
			if err != nil {
				return
			}
		case "message":
			err = dc.ReadExactBytes(z.MessageHash[:])
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *PPost) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "id"
	err = en.Append(0x84, 0xa2, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteInt64(z.ID)
	if err != nil {
		return
	}
	// write "created"
	err = en.Append(0xa7, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString((strfmt.DateTime).String(z.Created))
	if err != nil {
		return
	}
	// write "edited"
	err = en.Append(0xa6, 0x65, 0x64, 0x69, 0x74, 0x65, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteBool(z.IsEdited)
	if err != nil {
		return
	}
	// write "message"
	err = en.Append(0xa7, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.MessageHash[:])
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *PPost) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "id"
	o = append(o, 0x84, 0xa2, 0x69, 0x64)
	o = msgp.AppendInt64(o, z.ID)
	// string "created"
	o = append(o, 0xa7, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64)
	o = msgp.AppendString(o, (strfmt.DateTime).String(z.Created))
	// string "edited"
	o = append(o, 0xa6, 0x65, 0x64, 0x69, 0x74, 0x65, 0x64)
	o = msgp.AppendBool(o, z.IsEdited)
	// string "message"
	o = append(o, 0xa7, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65)
	o = msgp.AppendBytes(o, z.MessageHash[:])
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PPost) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zcua uint32
	zcua, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zcua > 0 {
		zcua--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "id":
			z.ID, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		case "created":
			{
				var zxhx string
				zxhx, bts, err = msgp.ReadStringBytes(bts)
				z.Created = parseDateTime(zxhx)
			}
			if err != nil {
				return
			}
		case "edited":
			z.IsEdited, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				return
			}
		case "message":
			bts, err = msgp.ReadExactBytes(bts, z.MessageHash[:])
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *PPost) Msgsize() (s int) {
	s = 1 + 3 + msgp.Int64Size + 8 + msgp.StringPrefixSize + len((strfmt.DateTime).String(z.Created)) + 7 + msgp.BoolSize + 8 + msgp.ArrayHeaderSize + (16 * (msgp.ByteSize))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PThread) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zpks uint32
	zpks, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zpks > 0 {
		zpks--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "id":
			z.ID, err = dc.ReadInt32()
			if err != nil {
				return
			}
		case "slug":
			z.Slug, err = dc.ReadString()
			if err != nil {
				return
			}
		case "message":
			err = dc.ReadExactBytes(z.MessageHash[:])
			if err != nil {
				return
			}
		case "title":
			err = dc.ReadExactBytes(z.TitleHash[:])
			if err != nil {
				return
			}
		case "created":
			{
				var zjfb string
				zjfb, err = dc.ReadString()
				z.Created = parseDateTime(zjfb)
			}
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *PThread) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 5
	// write "id"
	err = en.Append(0x85, 0xa2, 0x69, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteInt32(z.ID)
	if err != nil {
		return
	}
	// write "slug"
	err = en.Append(0xa4, 0x73, 0x6c, 0x75, 0x67)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Slug)
	if err != nil {
		return
	}
	// write "message"
	err = en.Append(0xa7, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.MessageHash[:])
	if err != nil {
		return
	}
	// write "title"
	err = en.Append(0xa5, 0x74, 0x69, 0x74, 0x6c, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.TitleHash[:])
	if err != nil {
		return
	}
	// write "created"
	err = en.Append(0xa7, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString((strfmt.DateTime).String(z.Created))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *PThread) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 5
	// string "id"
	o = append(o, 0x85, 0xa2, 0x69, 0x64)
	o = msgp.AppendInt32(o, z.ID)
	// string "slug"
	o = append(o, 0xa4, 0x73, 0x6c, 0x75, 0x67)
	o = msgp.AppendString(o, z.Slug)
	// string "message"
	o = append(o, 0xa7, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65)
	o = msgp.AppendBytes(o, z.MessageHash[:])
	// string "title"
	o = append(o, 0xa5, 0x74, 0x69, 0x74, 0x6c, 0x65)
	o = msgp.AppendBytes(o, z.TitleHash[:])
	// string "created"
	o = append(o, 0xa7, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64)
	o = msgp.AppendString(o, (strfmt.DateTime).String(z.Created))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PThread) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zcxo uint32
	zcxo, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zcxo > 0 {
		zcxo--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "id":
			z.ID, bts, err = msgp.ReadInt32Bytes(bts)
			if err != nil {
				return
			}
		case "slug":
			z.Slug, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "message":
			bts, err = msgp.ReadExactBytes(bts, z.MessageHash[:])
			if err != nil {
				return
			}
		case "title":
			bts, err = msgp.ReadExactBytes(bts, z.TitleHash[:])
			if err != nil {
				return
			}
		case "created":
			{
				var zeff string
				zeff, bts, err = msgp.ReadStringBytes(bts)
				z.Created = parseDateTime(zeff)
			}
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *PThread) Msgsize() (s int) {
	s = 1 + 3 + msgp.Int32Size + 5 + msgp.StringPrefixSize + len(z.Slug) + 8 + msgp.ArrayHeaderSize + (16 * (msgp.ByteSize)) + 6 + msgp.ArrayHeaderSize + (16 * (msgp.ByteSize)) + 8 + msgp.StringPrefixSize + len((strfmt.DateTime).String(z.Created))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PUser) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zdnj uint32
	zdnj, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zdnj > 0 {
		zdnj--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "about":
			err = dc.ReadExactBytes(z.AboutHash[:])
			if err != nil {
				return
			}
		case "email":
			{
				var zobc string
				zobc, err = dc.ReadString()
				z.Email = strfmt.Email(zobc)
			}
			if err != nil {
				return
			}
		case "name":
			err = dc.ReadExactBytes(z.FullnameHash[:])
			if err != nil {
				return
			}
		case "nick":
			z.Nickname, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *PUser) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "about"
	err = en.Append(0x84, 0xa5, 0x61, 0x62, 0x6f, 0x75, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.AboutHash[:])
	if err != nil {
		return
	}
	// write "email"
	err = en.Append(0xa5, 0x65, 0x6d, 0x61, 0x69, 0x6c)
	if err != nil {
		return err
	}
	err = en.WriteString((strfmt.Email).String(z.Email))
	if err != nil {
		return
	}
	// write "name"
	err = en.Append(0xa4, 0x6e, 0x61, 0x6d, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteBytes(z.FullnameHash[:])
	if err != nil {
		return
	}
	// write "nick"
	err = en.Append(0xa4, 0x6e, 0x69, 0x63, 0x6b)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Nickname)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *PUser) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "about"
	o = append(o, 0x84, 0xa5, 0x61, 0x62, 0x6f, 0x75, 0x74)
	o = msgp.AppendBytes(o, z.AboutHash[:])
	// string "email"
	o = append(o, 0xa5, 0x65, 0x6d, 0x61, 0x69, 0x6c)
	o = msgp.AppendString(o, (strfmt.Email).String(z.Email))
	// string "name"
	o = append(o, 0xa4, 0x6e, 0x61, 0x6d, 0x65)
	o = msgp.AppendBytes(o, z.FullnameHash[:])
	// string "nick"
	o = append(o, 0xa4, 0x6e, 0x69, 0x63, 0x6b)
	o = msgp.AppendString(o, z.Nickname)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PUser) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zsnv uint32
	zsnv, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zsnv > 0 {
		zsnv--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "about":
			bts, err = msgp.ReadExactBytes(bts, z.AboutHash[:])
			if err != nil {
				return
			}
		case "email":
			{
				var zkgt string
				zkgt, bts, err = msgp.ReadStringBytes(bts)
				z.Email = strfmt.Email(zkgt)
			}
			if err != nil {
				return
			}
		case "name":
			bts, err = msgp.ReadExactBytes(bts, z.FullnameHash[:])
			if err != nil {
				return
			}
		case "nick":
			z.Nickname, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *PUser) Msgsize() (s int) {
	s = 1 + 6 + msgp.ArrayHeaderSize + (16 * (msgp.ByteSize)) + 6 + msgp.StringPrefixSize + len((strfmt.Email).String(z.Email)) + 5 + msgp.ArrayHeaderSize + (16 * (msgp.ByteSize)) + 5 + msgp.StringPrefixSize + len(z.Nickname)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *PVersion) DecodeMsg(dc *msgp.Reader) (err error) {
	{
		var zema uint32
		zema, err = dc.ReadUint32()
		(*z) = PVersion(zema)
	}
	if err != nil {
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z PVersion) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteUint32(uint32(z))
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z PVersion) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendUint32(o, uint32(z))
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *PVersion) UnmarshalMsg(bts []byte) (o []byte, err error) {
	{
		var zpez uint32
		zpez, bts, err = msgp.ReadUint32Bytes(bts)
		(*z) = PVersion(zpez)
	}
	if err != nil {
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z PVersion) Msgsize() (s int) {
	s = msgp.Uint32Size
	return
}
