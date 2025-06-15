// This package provides basic types for blockchain environment manipulation, such as Uint128, PublicKey, NearGas, and other helpful types for smart contract development.
package types

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/bits"
	"strconv"
	"strings"

	"github.com/mr-tron/base58"
)

const (
	ONE_NEAR      = 1_000_000_000_000_000_000_000_000
	ONE_MILI_NEAR = 1_000_000_000_000_000_000_000
	ONE_TERA_GAS  = 1_000_000_000_000
	ONE_GIGA_GAS  = 1_000_000_000
)

const (
	ErrLoadingUint128FromBEBytes = "(UINT128_ERROR): error while loading Uint128 from BE bytes"
	ErrLoadingUint128FromLEBytes = "(UINT128_ERROR): error while loading Uint128 from LE bytes"
	ErrOverflow                  = "(UINT128_ERROR): overflow"
	ErrUnderflow                 = "(UINT128_ERROR): underflow"
	ErrOverflowInMultiplication  = "(UINT128_ERROR): overflow in multiplication"
	ErrDivisionByZero            = "(UINT128_ERROR): division by zero"
	ErrDivideByZero              = "(UINT128_ERROR): divide by zero"
	ErrIncorrectLen              = "(UINT128_ERROR): incorrect length"
	ErrUint128Overflow           = "(UINT128_ERROR): uint128 overflow"
	ErrIncorrectSymbolInString   = "(UINT128_ERROR): incorrect symbol in string"
)

// ContractInputOptions represents the options for obtaining smart contract input from the user.
//
// Input can be provided in two formats:
// 1. Raw bytes
// 2. Structured as JSON
//
// For more details, refer to the `ContractInput` method in `env.go`.
type ContractInputOptions struct {
	IsRawBytes bool
}

// Uint128 represents a 128-bit unsigned integer.
//
// This struct is used extensively in various operations, particularly with NEAR native tokens.
// Additionally, it can be utilized to create and manage custom tokens.
type Uint128 struct {
	Hi uint64
	Lo uint64
}

func LoadUint128BE(b []byte) (Uint128, error) {
	if len(b) != 16 {
		return Uint128{0, 0}, errors.New(ErrLoadingUint128FromBEBytes)
	}

	hi := binary.BigEndian.Uint64(b[:8])
	lo := binary.BigEndian.Uint64(b[8:16])

	return Uint128{Hi: hi, Lo: lo}, nil
}

func LoadUint128LE(b []byte) (Uint128, error) {
	if len(b) != 16 {
		return Uint128{0, 0}, errors.New(ErrLoadingUint128FromLEBytes)
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
	return hex.EncodeToString(u.ToLE())
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

func sub64(x, y, borrow uint64) (diff, borrowOut uint64) {
	diff = x - y - borrow
	if x < y+borrow {
		borrowOut = 1
	} else {
		borrowOut = 0
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
		return Uint128{0, 0}, errors.New(ErrOverflow)
	}
	return Uint128{Lo: lo, Hi: hi}, nil
}

func (u Uint128) SafeAdd64(v uint64) (Uint128, error) {
	lo, carry := add64(u.Lo, v, 0)
	hi, carry2 := add64(u.Hi, 0, carry)
	if carry2 != 0 {
		return Uint128{0, 0}, errors.New(ErrOverflow)
	}
	return Uint128{Lo: lo, Hi: hi}, nil
}

func (u Uint128) Add(v Uint128) (Uint128, error) {
	lo, carry := add64(u.Lo, v.Lo, 0)
	hi, carry2 := add64(u.Hi, v.Hi, carry)
	if carry2 != 0 {
		return Uint128{0, 0}, errors.New(ErrOverflow)
	}
	return Uint128{Lo: lo, Hi: hi}, nil
}

func (u Uint128) Sub(v Uint128) (Uint128, error) {
	lo, borrow := sub64(u.Lo, v.Lo, 0)
	hi, borrow2 := sub64(u.Hi, v.Hi, borrow)
	if borrow2 != 0 {
		return Uint128{0, 0}, errors.New(ErrUnderflow)
	}
	return Uint128{Lo: lo, Hi: hi}, nil
}

func (u Uint128) Mul(v Uint128) (Uint128, error) {
	lo, hi := mul128(u.Lo, v.Lo)
	hi1 := u.Hi * v.Lo
	hi2 := u.Lo * v.Hi

	newHi, overflow1 := add64(hi, hi1, 0)
	if overflow1 != 0 {
		return Uint128{0, 0}, errors.New(ErrOverflowInMultiplication)
	}

	newHi, overflow2 := add64(newHi, hi2, 0)
	if overflow2 != 0 {
		return Uint128{0, 0}, errors.New(ErrOverflowInMultiplication)
	}

	return Uint128{Lo: lo, Hi: newHi}, nil
}

func (u Uint128) Div(v Uint128) (Uint128, error) {
	if v.Lo == 0 && v.Hi == 0 {
		return Uint128{0, 0}, errors.New(ErrDivisionByZero)
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

func (u Uint128) SafeQuoRem64(v uint64) (q Uint128, r uint64, err error) {
	if v == 0 {
		return Uint128{0, 0}, 0, errors.New(ErrDivideByZero)
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
		return Uint128{0, 0}, errors.New(ErrDivisionByZero)
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

func isUint128Overflow(s string) bool {
	if len(s) > 39 {
		return true
	}

	if len(s) < 39 {
		return false
	}

	maxUint128Str := "340282366920938463463374607431768211455"
	for i := 0; i < 39; i++ {
		sDigit, _ := strconv.Atoi(string(s[i]))
		maxDigit, _ := strconv.Atoi(string(maxUint128Str[i]))

		if sDigit > maxDigit {
			return true
		} else if sDigit < maxDigit {
			return false
		}
	}
	return false
}

// U128FromString transforms a string into a Uint128 type.
//
// Returns an error if the string length is zero, exceeds 40 characters, or causes a Uint128 overflow.
func U128FromString(s string) (Uint128, error) {
	var res Uint128
	var err error

	if len(s) == 0 || len(s) > 40 {
		return Uint128{0, 0}, errors.New(ErrIncorrectLen)
	}

	if isUint128Overflow(s) {
		return Uint128{0, 0}, errors.New(ErrUint128Overflow)
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
			return Uint128{0, 0}, errors.New(ErrIncorrectLen)
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

func IntToString(n int) string {
	if n == 0 {
		return "0"
	}

	negative := n < 0
	if negative {
		n = -n
	}

	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		digits = append([]byte{'-'}, digits...)
	}

	return string(digits)
}

func BoolToUnit(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}

// NearGas represents the Gas consumed during smart contract execution.
//
// Each function call consumes a certain amount of gas. Optimizing and improving the efficiency of your code
// will reduce the amount of gas spent, resulting in faster transactions.
type NearGas struct {
	Inner uint64
}

// CurveType represents the type of cryptographic curve used.
type CurveType byte

const (
	// ED25519 represents the ed25519 curve.
	ED25519 CurveType = iota
	// SECP256K1 represents the secp256k1 curve.
	SECP256K1
)

// String returns the string representation of the CurveType.
func (c CurveType) String() string {
	switch c {
	case ED25519:
		return "ed25519"
	case SECP256K1:
		return "secp256k1"
	default:
		return "unknown"
	}
}

// DataLen returns the length of the data for the given CurveType.
func (c CurveType) DataLen() int {
	switch c {
	case ED25519:
		return 32
	case SECP256K1:
		return 64
	default:
		return 0
	}
}

// ParseCurveType parses a string to a CurveType.
func ParseCurveType(s string) (CurveType, error) {
	switch strings.ToLower(s) {
	case "ed25519":
		return ED25519, nil
	case "secp256k1":
		return SECP256K1, nil
	default:
		return 0, errors.New("unknown curve type")
	}
}

// PublicKey represents a public key with a specific curve type.
type PublicKey struct {
	Curve CurveType
	Data  []byte
}

// NewPublicKey creates a new PublicKey with the given curve type and data.
func NewPublicKey(curve CurveType, data []byte) (*PublicKey, error) {
	if len(data) != curve.DataLen() {
		return nil, errors.New("invalid data length for curve")
	}
	return &PublicKey{Curve: curve, Data: data}, nil
}

// PublicKeyFromString parses a string to create a PublicKey.
func PublicKeyFromString(s string) (*PublicKey, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid public key format")
	}

	curve, err := ParseCurveType(parts[0])
	if err != nil {
		return nil, err
	}

	data, err := base58.Decode(parts[1])
	if err != nil {
		return nil, errors.New("failed to decode Base58")
	}

	return NewPublicKey(curve, data)
}

// ToHexString returns the hexadecimal string representation of the PublicKey.
func (pk *PublicKey) ToHexString() string {
	return pk.Curve.String() + ":" + hex.EncodeToString(pk.Data)
}

// ToBase58String returns the Base58 string representation of the PublicKey.
func (pk *PublicKey) ToBase58String() string {
	return pk.Curve.String() + ":" + base58.Encode(pk.Data)
}

// Bytes returns the byte representation of the PublicKey.
func (pk *PublicKey) Bytes() []byte {
	curveByte := byte(pk.Curve)
	result := []byte{curveByte}
	result = append(result, pk.Data...)
	return result
}
