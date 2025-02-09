package types

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

// ContractInputOptions

type ContractInputOptions struct {
	IsRawBytes bool
}

//

// Uint128
type Uint128 struct {
	Hi uint64
	Lo uint64
}

func LoadUint128BE(b []byte) (Uint128, error) {
	if len(b) != 16 {
		return Uint128{0, 0}, errors.New("Error while Loading Uint128 from BE Bytes")
	}

	hi := binary.BigEndian.Uint64(b[:8])
	lo := binary.BigEndian.Uint64(b[8:16])

	return Uint128{Hi: hi, Lo: lo}, nil
}

func LoadUint128LE(b []byte) (Uint128, error) {
	if len(b) != 16 {
		return Uint128{0, 0}, errors.New("Error while Loading Uint128 from LE Bytes")
	}

	lo := binary.LittleEndian.Uint64(b[:8])
	hi := binary.LittleEndian.Uint64(b[8:16])

	return Uint128{Hi: hi, Lo: lo}, nil
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

func (u Uint128) ShiftLeft(bits uint) Uint128 {
	if bits >= 64 {
		return Uint128{Lo: 0, Hi: u.Lo << (bits - 64)}
	}
	return Uint128{Lo: u.Lo << bits, Hi: (u.Hi << bits) | (u.Lo >> (64 - bits))}
}

func (u Uint128) ShiftRight(bits uint) Uint128 {
	if bits >= 64 {
		return Uint128{Lo: u.Hi >> (bits - 64), Hi: 0}
	}
	return Uint128{Lo: (u.Lo >> bits) | (u.Hi << (64 - bits)), Hi: u.Hi >> bits}
}

func (u Uint128) GreaterOrEqual(v Uint128) bool {
	if u.Hi > v.Hi {
		return true
	}
	if u.Hi < v.Hi {
		return false
	}
	return u.Lo >= v.Lo
}

func sub64(x, y, borrow uint64) (diff, borrowOut uint64) {
	diff = x - y - borrow
	if x < y || (x == y && borrow != 0) {
		borrowOut = 1
	}
	return
}

func mul128(x, y uint64) (lo, hi uint64) {
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

func (u Uint128) Add(v Uint128) (Uint128, error) {
	lo, carry := add64(u.Lo, v.Lo, 0)
	hi, carry2 := add64(u.Hi, v.Hi, carry)
	if carry2 != 0 {
		return Uint128{0, 0}, errors.New("Overflow")
	}
	return Uint128{Lo: lo, Hi: hi}, nil
}
func (u Uint128) Sub(v Uint128) (Uint128, error) {
	lo, borrow := sub64(u.Lo, v.Lo, 0)
	hi, borrow2 := sub64(u.Hi, v.Hi, borrow)
	if borrow2 != 0 {
		return Uint128{0, 0}, errors.New("Underflow")
	}
	return Uint128{Lo: lo, Hi: hi}, nil
}

func (u Uint128) Mul(v Uint128) (Uint128, error) {
	lo, hi := mul128(u.Lo, v.Lo)
	hi1 := u.Hi * v.Lo
	hi2 := u.Lo * v.Hi

	newHi, overflow1 := add64(hi, hi1, 0)
	if overflow1 != 0 {
		return Uint128{0, 0}, errors.New("Overflow in multiplication")
	}

	newHi, overflow2 := add64(newHi, hi2, 0)
	if overflow2 != 0 {
		return Uint128{0, 0}, errors.New("Overflow in multiplication")
	}

	return Uint128{Lo: lo, Hi: newHi}, nil
}

func (u Uint128) Div(v Uint128) (Uint128, error) {
	if v.Lo == 0 && v.Hi == 0 {
		return Uint128{0, 0}, errors.New("Division by zero")
	}

	var result Uint128
	var remainder Uint128 = u

	for remainder.GreaterOrEqual(v) {
		shift := v
		multiple := Uint128{Lo: 1, Hi: 0}

		for remainder.GreaterOrEqual(shift.ShiftLeft(1)) {
			shift = shift.ShiftLeft(1)
			multiple = multiple.ShiftLeft(1)
		}

		remainder, _ = remainder.Sub(shift)
		result, _ = result.Add(multiple)
	}

	return result, nil
}

func (u Uint128) Bit(i int) uint {
	if i < 64 {
		return uint((u.Lo >> i) & 1)
	}
	return uint((u.Hi >> (i - 64)) & 1)
}

func (u Uint128) Lsh(n uint) Uint128 {
	if n == 0 {
		return u
	} else if n < 64 {
		return Uint128{
			Lo: u.Lo << n,
			Hi: (u.Hi << n) | (u.Lo >> (64 - n)),
		}
	} else {
		return Uint128{
			Lo: 0,
			Hi: u.Lo << (n - 64),
		}
	}
}

func (u Uint128) Cmp(v Uint128) int {
	if u.Hi > v.Hi || (u.Hi == v.Hi && u.Lo > v.Lo) {
		return 1
	} else if u.Hi < v.Hi || (u.Hi == v.Hi && u.Lo < v.Lo) {
		return -1
	}
	return 0
}

func (u Uint128) Mod(v Uint128) (Uint128, error) {
	if v.Lo == 0 && v.Hi == 0 {
		return Uint128{0, 0}, errors.New("Division by zero")
	}
	_, remainder := u.QuoRem64(v.Lo)
	return U64ToUint128(remainder), nil
}

func (u Uint128) And(v Uint128) Uint128 {
	return Uint128{Lo: u.Lo & v.Lo, Hi: u.Hi & v.Hi}
}

func (u Uint128) Or(v Uint128) Uint128 {
	return Uint128{Lo: u.Lo | v.Lo, Hi: u.Hi | v.Hi}
}

func (u Uint128) Xor(v Uint128) Uint128 {
	return Uint128{Lo: u.Lo ^ v.Lo, Hi: u.Hi ^ v.Hi}
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
	Inner uint64
}

// Near Gas

func Uint64ToString(n uint64) string {
	if n == 0 {
		return "0"
	}

	var result []byte
	for n > 0 {
		result = append([]byte{byte('0' + n%10)}, result...)
		n /= 10
	}
	return string(result)
}

func BoolToUnit(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}
