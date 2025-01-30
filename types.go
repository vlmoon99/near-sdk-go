package main

import (
	"encoding/binary"
	"encoding/hex"
)

type Uint128 struct {
	Hi, Lo uint64
}

func (u Uint128) GetBytes() []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[:8], u.Lo)
	binary.BigEndian.PutUint64(buf[8:], u.Hi)
	return buf
}

func (u Uint128) String() string {
	return hex.EncodeToString(u.GetBytes())
}

func FromBytes(b []byte) Uint128 {
	lo := binary.LittleEndian.Uint64(b[8:])
	hi := binary.LittleEndian.Uint64(b[:8])
	return Uint128{hi, lo}
}
