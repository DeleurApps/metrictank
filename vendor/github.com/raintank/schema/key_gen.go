package schema

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *AMKey) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "MKey":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			if err != nil {
				err = msgp.WrapError(err, "MKey")
				return
			}
			for zb0002 > 0 {
				zb0002--
				field, err = dc.ReadMapKeyPtr()
				if err != nil {
					err = msgp.WrapError(err, "MKey")
					return
				}
				switch msgp.UnsafeString(field) {
				case "Key":
					err = dc.ReadExactBytes((z.MKey.Key)[:])
					if err != nil {
						err = msgp.WrapError(err, "MKey", "Key")
						return
					}
				case "Org":
					z.MKey.Org, err = dc.ReadUint32()
					if err != nil {
						err = msgp.WrapError(err, "MKey", "Org")
						return
					}
				default:
					err = dc.Skip()
					if err != nil {
						err = msgp.WrapError(err, "MKey")
						return
					}
				}
			}
		case "Archive":
			err = z.Archive.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Archive")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *AMKey) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "MKey"
	// map header, size 2
	// write "Key"
	err = en.Append(0x82, 0xa4, 0x4d, 0x4b, 0x65, 0x79, 0x82, 0xa3, 0x4b, 0x65, 0x79)
	if err != nil {
		return
	}
	err = en.WriteBytes((z.MKey.Key)[:])
	if err != nil {
		err = msgp.WrapError(err, "MKey", "Key")
		return
	}
	// write "Org"
	err = en.Append(0xa3, 0x4f, 0x72, 0x67)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.MKey.Org)
	if err != nil {
		err = msgp.WrapError(err, "MKey", "Org")
		return
	}
	// write "Archive"
	err = en.Append(0xa7, 0x41, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65)
	if err != nil {
		return
	}
	err = z.Archive.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Archive")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *AMKey) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "MKey"
	// map header, size 2
	// string "Key"
	o = append(o, 0x82, 0xa4, 0x4d, 0x4b, 0x65, 0x79, 0x82, 0xa3, 0x4b, 0x65, 0x79)
	o = msgp.AppendBytes(o, (z.MKey.Key)[:])
	// string "Org"
	o = append(o, 0xa3, 0x4f, 0x72, 0x67)
	o = msgp.AppendUint32(o, z.MKey.Org)
	// string "Archive"
	o = append(o, 0xa7, 0x41, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65)
	o, err = z.Archive.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Archive")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *AMKey) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "MKey":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "MKey")
				return
			}
			for zb0002 > 0 {
				zb0002--
				field, bts, err = msgp.ReadMapKeyZC(bts)
				if err != nil {
					err = msgp.WrapError(err, "MKey")
					return
				}
				switch msgp.UnsafeString(field) {
				case "Key":
					bts, err = msgp.ReadExactBytes(bts, (z.MKey.Key)[:])
					if err != nil {
						err = msgp.WrapError(err, "MKey", "Key")
						return
					}
				case "Org":
					z.MKey.Org, bts, err = msgp.ReadUint32Bytes(bts)
					if err != nil {
						err = msgp.WrapError(err, "MKey", "Org")
						return
					}
				default:
					bts, err = msgp.Skip(bts)
					if err != nil {
						err = msgp.WrapError(err, "MKey")
						return
					}
				}
			}
		case "Archive":
			bts, err = z.Archive.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Archive")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *AMKey) Msgsize() (s int) {
	s = 1 + 5 + 1 + 4 + msgp.ArrayHeaderSize + (16 * (msgp.ByteSize)) + 4 + msgp.Uint32Size + 8 + z.Archive.Msgsize()
	return
}

// DecodeMsg implements msgp.Decodable
func (z *Key) DecodeMsg(dc *msgp.Reader) (err error) {
	err = dc.ReadExactBytes((z)[:])
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Key) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteBytes((z)[:])
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Key) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendBytes(o, (z)[:])
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Key) UnmarshalMsg(bts []byte) (o []byte, err error) {
	bts, err = msgp.ReadExactBytes(bts, (z)[:])
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Key) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize + (16 * (msgp.ByteSize))
	return
}

// DecodeMsg implements msgp.Decodable
func (z *MKey) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Key":
			err = dc.ReadExactBytes((z.Key)[:])
			if err != nil {
				err = msgp.WrapError(err, "Key")
				return
			}
		case "Org":
			z.Org, err = dc.ReadUint32()
			if err != nil {
				err = msgp.WrapError(err, "Org")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *MKey) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "Key"
	err = en.Append(0x82, 0xa3, 0x4b, 0x65, 0x79)
	if err != nil {
		return
	}
	err = en.WriteBytes((z.Key)[:])
	if err != nil {
		err = msgp.WrapError(err, "Key")
		return
	}
	// write "Org"
	err = en.Append(0xa3, 0x4f, 0x72, 0x67)
	if err != nil {
		return
	}
	err = en.WriteUint32(z.Org)
	if err != nil {
		err = msgp.WrapError(err, "Org")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *MKey) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "Key"
	o = append(o, 0x82, 0xa3, 0x4b, 0x65, 0x79)
	o = msgp.AppendBytes(o, (z.Key)[:])
	// string "Org"
	o = append(o, 0xa3, 0x4f, 0x72, 0x67)
	o = msgp.AppendUint32(o, z.Org)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *MKey) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Key":
			bts, err = msgp.ReadExactBytes(bts, (z.Key)[:])
			if err != nil {
				err = msgp.WrapError(err, "Key")
				return
			}
		case "Org":
			z.Org, bts, err = msgp.ReadUint32Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Org")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *MKey) Msgsize() (s int) {
	s = 1 + 4 + msgp.ArrayHeaderSize + (16 * (msgp.ByteSize)) + 4 + msgp.Uint32Size
	return
}
