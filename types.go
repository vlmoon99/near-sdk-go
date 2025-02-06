package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/bits"
)

const (
	ONE_NEAR      = 1_000_000_000_000_000_000_000_000
	ONE_MILI_NEAR = 1_000_000_000_000_000_000_000
	ONE_TERA_GAS  = 1_000_000_000_000
	ONE_GIGA_GAS  = 1_000_000_000
)

// Uint128
type Uint128 struct {
	Hi uint64
	Lo uint64
}

type Uint256 struct {
	Lo Uint128
	Hi Uint128
}

func LoadUint128BE(b []byte) Uint128 {
	if len(b) != 16 {
		PanicStr("byte slice must be exactly 16 bytes long")
	}

	hi := binary.BigEndian.Uint64(b[:8])
	lo := binary.BigEndian.Uint64(b[8:16])

	return Uint128{Hi: hi, Lo: lo}
}

func LoadUint128LE(b []byte) Uint128 {
	if len(b) != 16 {
		PanicStr("byte slice must be exactly 16 bytes long")
	}

	lo := binary.LittleEndian.Uint64(b[:8])
	hi := binary.LittleEndian.Uint64(b[8:16])

	return Uint128{Hi: hi, Lo: lo}
}

func (u Uint128) ToBE() []byte {
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b[:8], u.Hi)
	binary.BigEndian.PutUint64(b[8:16], u.Lo)
	return b
}

func (u Uint128) ToLE() []byte {
	b := make([]byte, 16)
	binary.LittleEndian.PutUint64(b[:8], u.Lo)
	binary.LittleEndian.PutUint64(b[8:16], u.Hi)
	return b
}

func (u Uint128) HexLE() string {
	return hex.EncodeToString(u.ToLE())
}

func (u Uint128) HexBE() string {
	return hex.EncodeToString(u.ToBE())
}

func U64ToBE(value uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, value)
	return b
}

func U64ToLE(value uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, value)
	return b
}

func U64ToUint128(value uint64) Uint128 {
	return Uint128{Hi: 0, Lo: value}
}

func BoolToUnit(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}

func (u Uint128) SafeQuoRem64(v uint64) (q Uint128, r uint64, err error) {
	if v == 0 {
		return Uint128{}, 0, errors.New("Divide by zero")
	}

	if u.Hi < v {
		q.Lo, r = bits.Div64(u.Hi, u.Lo, v)
	} else {
		q.Hi, r = bits.Div64(0, u.Hi, v)
		q.Lo, r = bits.Div64(r, u.Lo, v)
	}
	return
}

func (u Uint128) QuoRem64(v uint64) (q Uint128, r uint64) {
	if u.Hi < v {
		q.Lo, r = bits.Div64(u.Hi, u.Lo, v)
	} else {
		q.Hi, r = bits.Div64(0, u.Hi, v)
		q.Lo, r = bits.Div64(r, u.Lo, v)
	}
	return
}

func mul64(x, y uint64) (lo, hi uint64) {
	const mask32 = (1 << 32) - 1
	x0, x1 := x&mask32, x>>32
	y0, y1 := y&mask32, y>>32

	w0 := x0 * y0
	t := x1*y0 + w0>>32
	w1 := t & mask32
	w2 := t >> 32

	w1 += x0 * y1

	lo = x * y
	hi = x1*y1 + w2 + w1>>32
	return
}

func add64(x, y, carry uint64) (sum, carryOut uint64) {
	sum = x + y + carry
	if sum < x || (sum == x && carry != 0) {
		carryOut = 1
	}
	return
}

func (u Uint128) Mul64(v uint64) Uint128 {
	lo, hi := mul64(u.Lo, v)
	return Uint128{Lo: lo, Hi: hi}
}

func (u Uint128) SafeMul64(v uint64) (Uint128, error) {
	lo, hi := mul64(u.Lo, v)
	hi += u.Hi * v
	if hi < u.Hi {
		return Uint128{0, 0}, errors.New("Overflow")
	}
	return Uint128{Lo: lo, Hi: hi}, nil
}

func (u Uint128) SafeAdd64(v uint64) (Uint128, error) {
	lo, carry := add64(u.Lo, v, 0)
	hi, carry2 := add64(u.Hi, 0, carry)
	if carry2 != 0 {
		return Uint128{0, 0}, errors.New("Overflow")
	}
	return Uint128{Lo: lo, Hi: hi}, nil
}

func U128FromString(s string) (Uint128, error) {
	res := Uint128{0, 0}
	var err error

	if len(s) == 0 || len(s) > 40 {
		return Uint128{0, 0}, errors.New("Incorrect len")
	}

	res, err = processPart(s)

	if err != nil {
		return Uint128{0, 0}, err
	}
	return res, nil

}

func processPart(s string) (Uint128, error) {
	res := Uint128{0, 0}
	var err error

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch < '0' || ch > '9' {
			return Uint128{0, 0}, errors.New("Incorrect symbol in string")
		}

		val := uint64(ch - '0')
		res, err = res.SafeMul64(10)

		if err != nil {
			return Uint128{0, 0}, err
		}

		res, err = res.SafeAdd64(val)

		if err != nil {
			return Uint128{0, 0}, err
		}
	}
	return res, nil
}

func (u Uint128) String() string {
	if u.Hi == 0 && u.Lo == 0 {
		return "0"
	}
	buf := []byte("0000000000000000000000000000000000000000")
	for i := len(buf); ; i -= 19 {
		q, r := u.QuoRem64(1e19)
		var n int
		for ; r != 0; r /= 10 {
			n++
			buf[i-n] += byte(r % 10)
		}
		if q.Hi == 0 && q.Lo == 0 {
			return string(buf[i-n:])
		}
		u = q
	}
}

//Uint128

// Near Gas
type NearGas struct {
	inner uint64
}

//Near Gas
