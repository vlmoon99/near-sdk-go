package types

import (
	"fmt"
	"strconv"
	"testing"
)

func TestUint64ToString(t *testing.T) {
	testCases := []struct {
		input    uint64
		expected string
	}{
		{input: 0, expected: "0"},
		{input: 1, expected: "1"},
		{input: 10, expected: "10"},
		{input: 12345, expected: "12345"},
		{input: 9876543210, expected: "9876543210"},
		{input: 18446744073709551615, expected: "18446744073709551615"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			actual := Uint64ToString(tc.input)
			if actual != tc.expected {
				t.Errorf("Uint64ToString(%d) = %s; want %s", tc.input, actual, tc.expected)
			}
		})
	}
}

func TestBoolToUnit(t *testing.T) {
	testCases := []struct {
		input    bool
		expected uint64
	}{
		{input: true, expected: 1},
		{input: false, expected: 0},
	}

	for _, tc := range testCases {
		t.Run(strconv.FormatBool(tc.input), func(t *testing.T) {
			actual := BoolToUnit(tc.input)
			if actual != tc.expected {
				t.Errorf("BoolToUnit(%t) = %d; want %d", tc.input, actual, tc.expected)
			}
		})
	}
}

func TestLoadUint128BE(t *testing.T) {

	testCases := []struct {
		name        string
		input       []byte
		expected    Uint128
		expectedErr bool
	}{
		{
			name:        "Valid Big Endian",
			input:       []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02},
			expected:    Uint128{Hi: 1, Lo: 2},
			expectedErr: false,
		},
		{
			name:        "Invalid Length",
			input:       []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected:    Uint128{0, 0},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := LoadUint128BE(tc.input)
			if (err != nil) != tc.expectedErr {
				t.Errorf("LoadUint128BE(%v) error = %v, wantErr %v", tc.input, err, tc.expectedErr)
				return
			}
			if actual != tc.expected {
				t.Errorf("LoadUint128BE(%v) = %v, want %v", tc.input, actual, tc.expected)
			}
		})
	}
}

func TestLoadUint128LE(t *testing.T) {
	testCases := []struct {
		name        string
		input       []byte
		expected    Uint128
		expectedErr bool
	}{
		{
			name:        "Valid Little Endian",
			input:       []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected:    Uint128{Hi: 2, Lo: 1},
			expectedErr: false,
		},
		{
			name:        "Invalid Length",
			input:       []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected:    Uint128{0, 0},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := LoadUint128LE(tc.input)
			if (err != nil) != tc.expectedErr {
				t.Errorf("LoadUint128LE(%v) error = %v, wantErr %v", tc.input, err, tc.expectedErr)
				return
			}
			if actual != tc.expected {
				t.Errorf("LoadUint128LE(%v) = %v, want %v", tc.input, actual, tc.expected)
			}
		})
	}
}

func TestUint128_ToBE(t *testing.T) {
	u := Uint128{Hi: 1, Lo: 2}
	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02}
	actual := u.ToBE()
	if !bytesEqual(actual, expected) {
		t.Errorf("ToBE() = %v, want %v", actual, expected)
	}
}

func TestUint128_ToLE(t *testing.T) {
	u := Uint128{Hi: 2, Lo: 1}
	expected := []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	actual := u.ToLE()
	if !bytesEqual(actual, expected) {
		t.Errorf("ToLE() = %v, want %v", actual, expected)
	}
}

func TestU64ToBE(t *testing.T) {
	value := uint64(12345)
	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x30, 0x39}
	actual := U64ToBE(value)
	if !bytesEqual(actual, expected) {
		t.Errorf("U64ToBE(%d) = %v, want %v", value, actual, expected)
	}
}

func TestU64ToLE(t *testing.T) {
	value := uint64(12345)
	expected := []byte{0x39, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	actual := U64ToLE(value)
	if !bytesEqual(actual, expected) {
		t.Errorf("U64ToLE(%d) = %v, want %v", value, actual, expected)
	}
}

func TestU64ToUint128(t *testing.T) {
	value := uint64(12345)
	expected := Uint128{Hi: 0, Lo: value}
	actual := U64ToUint128(value)
	if actual != expected {
		t.Errorf("U64ToUint128(%d) = %v, want %v", value, actual, expected)
	}
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestU128FromString(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    Uint128
		expectedErr bool
	}{
		{
			name:        "Valid Input",
			input:       "12345678901234567890",
			expected:    Uint128{Hi: 0, Lo: 12345678901234567890},
			expectedErr: false,
		},
		{
			name:        "Zero Input",
			input:       "0",
			expected:    Uint128{Hi: 0, Lo: 0},
			expectedErr: false,
		},
		{
			name:        "Max Uint64 Input",
			input:       "18446744073709551615",
			expected:    Uint128{Hi: 0, Lo: 18446744073709551615},
			expectedErr: false,
		},
		{
			name:        "Small Hi, Large Lo",
			input:       "18446744073709551616",
			expected:    Uint128{Hi: 1, Lo: 0},
			expectedErr: false,
		},
		{
			name:        "Too Long Input",
			input:       "12345678901234567890123456789012345678901",
			expected:    Uint128{0, 0},
			expectedErr: true,
		},
		{
			name:        "Empty Input",
			input:       "",
			expected:    Uint128{0, 0},
			expectedErr: true,
		},
		{
			name:        "Invalid Character",
			input:       "123a456",
			expected:    Uint128{0, 0},
			expectedErr: true,
		},
		{
			name:        "Max Uint128 Input",
			input:       "340282366920938463463374607431768211455",
			expected:    Uint128{Hi: 0xFFFFFFFFFFFFFFFF, Lo: 0xFFFFFFFFFFFFFFFF},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := U128FromString(tc.input)
			if (err != nil) != tc.expectedErr {
				t.Errorf("U128FromString(%q) error = %v, wantErr %v", tc.input, err, tc.expectedErr)
				return
			}
			if actual != tc.expected {
				t.Errorf("U128FromString(%q) = %v, want %v", tc.input, actual, tc.expected)
			}
		})
	}
}

func TestUint128_String(t *testing.T) {
	testCases := []struct {
		name     string
		input    Uint128
		expected string
	}{
		{
			name:     "Zero",
			input:    Uint128{Hi: 0, Lo: 0},
			expected: "0",
		},
		{
			name:     "Small Lo",
			input:    Uint128{Hi: 0, Lo: 12345},
			expected: "12345",
		},
		{
			name:     "Max Uint64",
			input:    Uint128{Hi: 0, Lo: 18446744073709551615},
			expected: "18446744073709551615",
		},
		{
			name:     "Small Hi, large Lo",
			input:    Uint128{Hi: 1, Lo: 0},
			expected: "18446744073709551616",
		},
		{
			name:     "Large Number",
			input:    Uint128{Hi: 0xFFFFFFFFFFFFFFFF, Lo: 0xFFFFFFFFFFFFFFFF},
			expected: "340282366920938463463374607431768211455",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.input.String()
			if actual != tc.expected {
				t.Errorf("String() for %v = %q, want %q", tc.input, actual, tc.expected)
			}
		})
	}
}

// Uint128 Math tests

func TestMul64(t *testing.T) {
	testCases := []struct {
		x  uint64
		y  uint64
		lo uint64
		hi uint64
	}{
		{x: 2, y: 3, lo: 6, hi: 0},
		{x: 18446744073709551615, y: 2, lo: 18446744073709551614, hi: 1},
		{x: 4294967295, y: 4294967295, lo: 18446744065119617025, hi: 0},
		{x: 1, y: 0, lo: 0, hi: 0},
		{x: 0, y: 1, lo: 0, hi: 0},
		{x: 18446744073709551615, y: 18446744073709551615, lo: 1, hi: 18446744073709551614},
		{x: 123456789, y: 987654321, lo: 121932631112635269, hi: 0},
		{x: 9223372036854775808, y: 2, lo: 0, hi: 1},
	}

	for _, tc := range testCases {
		t.Run(strconv.FormatUint(tc.x, 10)+"*"+strconv.FormatUint(tc.y, 10), func(t *testing.T) {
			lo, hi := mul64(tc.x, tc.y)
			if lo != tc.lo || hi != tc.hi {
				t.Errorf("mul64(%d, %d) = (%d, %d), want (%d, %d)", tc.x, tc.y, lo, hi, tc.lo, tc.hi)
			}
		})
	}
}

func TestAdd64(t *testing.T) {
	testCases := []struct {
		x        uint64
		y        uint64
		carryIn  uint64
		sum      uint64
		carryOut uint64
	}{
		{x: 1, y: 2, carryIn: 0, sum: 3, carryOut: 0},
		{x: 18446744073709551615, y: 1, carryIn: 0, sum: 0, carryOut: 1},
		{x: 18446744073709551615, y: 0, carryIn: 1, sum: 0, carryOut: 1},
		{x: 0, y: 0, carryIn: 0, sum: 0, carryOut: 0},
		{x: 0, y: 1, carryIn: 0, sum: 1, carryOut: 0},
		{x: 9223372036854775807, y: 1, carryIn: 1, sum: 9223372036854775809, carryOut: 0},
		{x: 123456789, y: 987654321, carryIn: 1, sum: 1111111111, carryOut: 0},
		{x: 0xFFFFFFFFFFFFFFFF, y: 0xFFFFFFFFFFFFFFFF, carryIn: 1, sum: 0xFFFFFFFFFFFFFFFF, carryOut: 1},
		{x: 0x8000000000000000, y: 0x8000000000000000, carryIn: 0, sum: 0, carryOut: 1},
	}

	for _, tc := range testCases {
		t.Run(strconv.FormatUint(tc.x, 10)+"+"+strconv.FormatUint(tc.y, 10)+"+"+strconv.FormatUint(tc.carryIn, 10), func(t *testing.T) {
			sum, carryOut := add64(tc.x, tc.y, tc.carryIn)
			if sum != tc.sum || carryOut != tc.carryOut {
				t.Errorf("add64(%d, %d, %d) = (%d, %d), want (%d, %d)", tc.x, tc.y, tc.carryIn, sum, carryOut, tc.sum, tc.carryOut)
			}
		})
	}
}

func TestSub64(t *testing.T) {
	testCases := []struct {
		x         uint64
		y         uint64
		borrowIn  uint64
		diff      uint64
		borrowOut uint64
	}{
		{x: 3, y: 2, borrowIn: 0, diff: 1, borrowOut: 0},
		{x: 0, y: 1, borrowIn: 0, diff: 18446744073709551615, borrowOut: 1},
		{x: 0, y: 0, borrowIn: 1, diff: 18446744073709551615, borrowOut: 1},
		{x: 1, y: 0, borrowIn: 0, diff: 1, borrowOut: 0},
		{x: 0, y: 0, borrowIn: 0, diff: 0, borrowOut: 0},
		{x: 18446744073709551615, y: 1, borrowIn: 0, diff: 18446744073709551614, borrowOut: 0},
		{x: 18446744073709551615, y: 18446744073709551615, borrowIn: 0, diff: 0, borrowOut: 0},
		{x: 18446744073709551615, y: 18446744073709551614, borrowIn: 0, diff: 1, borrowOut: 0},
		{x: 123456789, y: 987654321, borrowIn: 1, diff: 18446744072845354083, borrowOut: 1},
	}

	for _, tc := range testCases {
		t.Run(strconv.FormatUint(tc.x, 10)+"-"+strconv.FormatUint(tc.y, 10)+"-"+strconv.FormatUint(tc.borrowIn, 10), func(t *testing.T) {
			diff, borrowOut := sub64(tc.x, tc.y, tc.borrowIn)
			if diff != tc.diff || borrowOut != tc.borrowOut {
				t.Errorf("sub64(%d, %d, %d) = (%d, %d), want (%d, %d)", tc.x, tc.y, tc.borrowIn, diff, borrowOut, tc.diff, tc.borrowOut)
			}
		})
	}
}

func TestUint128_SafeAdd64(t *testing.T) {
	testCases := []struct {
		uStr        string
		v           uint64
		expectedStr string
		expectedErr bool
	}{
		{uStr: "100", v: 3, expectedStr: "103", expectedErr: false},
		{uStr: "18446744073709551615", v: 1, expectedStr: "18446744073709551616", expectedErr: false},
		{uStr: "9223372036854775808", v: 9223372036854775808, expectedStr: "18446744073709551616", expectedErr: false},
		{uStr: "0", v: 0, expectedStr: "0", expectedErr: false},
		//Result will be one because  "123456789012345678901234567890123456789000" will be parsed as 0 because its uStr > Uint128.Max
		{uStr: "123456789012345678901234567890123456789000", v: 1, expectedStr: "1", expectedErr: false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s+%d", tc.uStr, tc.v), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			expected, _ := U128FromString(tc.expectedStr)
			actual, err := u.SafeAdd64(tc.v)

			if (err != nil) != tc.expectedErr {
				t.Errorf("SafeAdd64(%s, %d) error = %v, wantErr %v", tc.uStr, tc.v, err, tc.expectedErr)
				return
			}
			if actual != expected {
				t.Errorf("SafeAdd64(%s, %d) = %s, want %s", tc.uStr, tc.v, actual.String(), expected.String())
			}
		})
	}
}

func TestUint128_Add(t *testing.T) {
	testCases := []struct {
		uStr        string
		vStr        string
		expectedStr string
		expectedErr bool
	}{
		{uStr: "100", vStr: "3", expectedStr: "103", expectedErr: false},
		{uStr: "18446744073709551615", vStr: "1", expectedStr: "18446744073709551616", expectedErr: false},
		{uStr: "9223372036854775808", vStr: "9223372036854775808", expectedStr: "18446744073709551616", expectedErr: false},
		{uStr: "0", vStr: "0", expectedStr: "0", expectedErr: false},
		// if input of some number (v or u) will be < than Uint128.Max than it wil return 0
		{uStr: "123456789012345678901234567890123456789000", vStr: "987654321098765432109876543210987654321000", expectedStr: "0", expectedErr: false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s+%s", tc.uStr, tc.vStr), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			v, _ := U128FromString(tc.vStr)
			expected, _ := U128FromString(tc.expectedStr)

			actual, err := u.Add(v)

			if (err != nil) != tc.expectedErr {
				t.Errorf("Add(%s, %s) error = %v, wantErr %v", tc.uStr, tc.vStr, err, tc.expectedErr)
				return
			}
			if actual != expected {
				t.Errorf("Add(%s, %s) = %v, want %v", u.String(), v.String(), actual.String(), expected.String())
			}
		})
	}
}

func TestUint128_Sub(t *testing.T) {
	testCases := []struct {
		uStr        string
		vStr        string
		expectedStr string
		expectedErr bool
	}{
		{uStr: "105", vStr: "3", expectedStr: "102", expectedErr: false},
		{uStr: "18446744073709551616", vStr: "1", expectedStr: "18446744073709551615", expectedErr: false},
		{uStr: "18446744073709551616", vStr: "9223372036854775808", expectedStr: "9223372036854775808", expectedErr: false},
		{uStr: "1", vStr: "2", expectedStr: "0", expectedErr: true},
		{uStr: "0", vStr: "0", expectedStr: "0", expectedErr: false},
		{uStr: "1000", vStr: "1", expectedStr: "999", expectedErr: false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s-%s", tc.uStr, tc.vStr), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			v, _ := U128FromString(tc.vStr)
			expected, _ := U128FromString(tc.expectedStr)

			actual, err := u.Sub(v)

			if (err != nil) != tc.expectedErr {
				t.Errorf("Sub(%s, %s) error = %v, wantErr %v", tc.uStr, tc.vStr, err, tc.expectedErr)
				return
			}
			if actual != expected {
				t.Errorf("Sub(%s, %s) = %v, want %v", tc.uStr, tc.vStr, actual, expected)
			}
		})
	}
}

func TestUint128_Mul(t *testing.T) {
	testCases := []struct {
		uStr        string
		vStr        string
		expectedStr string
		expectedErr bool
	}{
		{uStr: "100", vStr: "3", expectedStr: "300", expectedErr: false},
		{uStr: "33", vStr: "3", expectedStr: "99", expectedErr: false},
		{uStr: "18446744073709551615", vStr: "2", expectedStr: "36893488147419103230", expectedErr: false},
		{uStr: "9223372036854775808", vStr: "2", expectedStr: "18446744073709551616", expectedErr: false},
		{uStr: "123456789", vStr: "987654321", expectedStr: "121932631112635269", expectedErr: false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s * %s", tc.uStr, tc.vStr), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			v, _ := U128FromString(tc.vStr)
			expected, _ := U128FromString(tc.expectedStr)

			actual, err := u.Mul(v)
			if (err != nil) != tc.expectedErr {
				t.Errorf("Mul(%s, %s) error = %v, wantErr %v", tc.uStr, tc.vStr, err, tc.expectedErr)
				return
			}
			if actual != expected {
				t.Errorf("Mul(%s, %s) = %v, want %v", tc.uStr, tc.vStr, actual, expected)
			}
		})
	}
}

func TestUint128_Div(t *testing.T) {
	testCases := []struct {
		uStr        string
		vStr        string
		expectedStr string
		expectedErr bool
	}{
		{uStr: "6", vStr: "2", expectedStr: "3", expectedErr: false},
		{uStr: "18446744073709551616", vStr: "2", expectedStr: "9223372036854775808", expectedErr: false},
		{uStr: "1", vStr: "1", expectedStr: "1", expectedErr: false},
		{uStr: "18446744073709551616", vStr: "18446744073709551616", expectedStr: "1", expectedErr: false},
		{uStr: "1", vStr: "0", expectedStr: "", expectedErr: true},
		{uStr: "18446744073709551617", vStr: "18446744073709551617", expectedStr: "1", expectedErr: false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s / %s", tc.uStr, tc.vStr), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			v, _ := U128FromString(tc.vStr)
			expected, _ := U128FromString(tc.expectedStr)

			actual, err := u.Div(v)
			if (err != nil) != tc.expectedErr {
				t.Errorf("Div(%s, %s) error = %v, wantErr %v", tc.uStr, tc.vStr, err, tc.expectedErr)
				return
			}
			if actual != expected {
				t.Errorf("Div(%s, %s) = %v, want %v", tc.uStr, tc.vStr, actual, expected)
			}
		})
	}
}

func TestUint128_SafeQuoRem64(t *testing.T) {
	testCases := []struct {
		uStr         string
		v            uint64
		expectedQStr string
		expectedR    uint64
		expectedErr  bool
	}{
		{uStr: "6", v: 2, expectedQStr: "3", expectedR: 0, expectedErr: false},
		{uStr: "18446744073709551616", v: 2, expectedQStr: "9223372036854775808", expectedR: 0, expectedErr: false},
		{uStr: "1", v: 1, expectedQStr: "1", expectedR: 0, expectedErr: false},
		{uStr: "18446744073709551617", v: 1, expectedQStr: "18446744073709551617", expectedR: 0, expectedErr: false},
		{uStr: "1", v: 0, expectedQStr: "", expectedR: 0, expectedErr: true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s / %d", tc.uStr, tc.v), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			expectedQ, _ := U128FromString(tc.expectedQStr)

			actualQ, actualR, err := u.SafeQuoRem64(tc.v)

			if tc.expectedErr {
				if err == nil {
					t.Errorf("SafeQuoRem64(%s, %d) expected error, but got nil", tc.uStr, tc.v)
				}
				return
			} else if err != nil {
				t.Errorf("SafeQuoRem64(%s, %d) returned error: %v, but did not expect one", tc.uStr, tc.v, err)
				return
			}

			if actualQ != expectedQ {
				t.Errorf("SafeQuoRem64(%s, %d) quotient = %v, want %v", tc.uStr, tc.v, actualQ, expectedQ)
			}
			if actualR != tc.expectedR {
				t.Errorf("SafeQuoRem64(%s, %d) remainder = %d, want %d", tc.uStr, tc.v, actualR, tc.expectedR)
			}
		})
	}
}

func TestUint128_QuoRem64(t *testing.T) {
	testCases := []struct {
		uStr         string
		v            uint64
		expectedQStr string
		expectedR    uint64
	}{
		{uStr: "6", v: 2, expectedQStr: "3", expectedR: 0},
		{uStr: "18446744073709551616", v: 2, expectedQStr: "9223372036854775808", expectedR: 0},
		{uStr: "1", v: 1, expectedQStr: "1", expectedR: 0},
		{uStr: "18446744073709551617", v: 1, expectedQStr: "18446744073709551617", expectedR: 0},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s / %d", tc.uStr, tc.v), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			expectedQ, _ := U128FromString(tc.expectedQStr)

			actualQ, actualR := u.QuoRem64(tc.v)

			if actualQ != expectedQ {
				t.Errorf("QuoRem64(%s, %d) quotient = %v, want %v", tc.uStr, tc.v, actualQ, expectedQ)
			}
			if actualR != tc.expectedR {
				t.Errorf("QuoRem64(%s, %d) remainder = %d, want %d", tc.uStr, tc.v, actualR, tc.expectedR)
			}
		})
	}
}

func TestUint128_ShiftLeft(t *testing.T) {
	testCases := []struct {
		uStr        string
		bits        uint
		expectedStr string
	}{
		{uStr: "1", bits: 1, expectedStr: "2"},
		{uStr: "1", bits: 64, expectedStr: "18446744073709551616"},
		{uStr: "2", bits: 1, expectedStr: "4"},
		{uStr: "3", bits: 2, expectedStr: "12"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s<<%d", tc.uStr, tc.bits), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			expected, _ := U128FromString(tc.expectedStr)

			actual := u.ShiftLeft(tc.bits)
			if actual != expected {
				t.Errorf("ShiftLeft(%s, %d) = %v, want %v", tc.uStr, tc.bits, actual.String(), expected.String())
			}
		})
	}
}

func TestUint128_ShiftRight(t *testing.T) {
	testCases := []struct {
		uStr        string
		bits        uint
		expectedStr string
	}{
		{uStr: "2", bits: 1, expectedStr: "1"},
		{uStr: "4", bits: 1, expectedStr: "2"},
		{uStr: "18446744073709551616", bits: 64, expectedStr: "1"},
		{uStr: "18446744073709551616", bits: 1, expectedStr: "9223372036854775808"},
		{uStr: "18446744073709551616", bits: 63, expectedStr: "2"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s>>%d", tc.uStr, tc.bits), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			expected, _ := U128FromString(tc.expectedStr)

			actual := u.ShiftRight(tc.bits)
			if actual != expected {
				t.Errorf("ShiftRight(%s, %d) = %v, want %v", u.String(), tc.bits, actual.String(), expected.String())
			}
		})
	}
}

func TestUint128_GreaterOrEqual(t *testing.T) {
	testCases := []struct {
		uStr     string
		vStr     string
		expected bool
	}{
		{uStr: "100", vStr: "99", expected: true},
		{uStr: "99", vStr: "100", expected: false},
		{uStr: "100", vStr: "100", expected: true},
		{uStr: "18446744073709551615", vStr: "9223372036854775808", expected: true},
		{uStr: "9223372036854775808", vStr: "18446744073709551615", expected: false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s >= %s", tc.uStr, tc.vStr), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			v, _ := U128FromString(tc.vStr)

			actual := u.GreaterOrEqual(v)
			if actual != tc.expected {
				t.Errorf("GreaterOrEqual(%s, %s) = %v, want %v", u.String(), v.String(), actual, tc.expected)
			}
		})
	}
}

func TestUint128_Bit(t *testing.T) {
	testCases := []struct {
		uStr        string
		index       int
		expectedBit uint
	}{
		{uStr: "4", index: 2, expectedBit: 1},
		{uStr: "2", index: 1, expectedBit: 1},
		{uStr: "1", index: 0, expectedBit: 1},
		{uStr: "18446744073709551616", index: 64, expectedBit: 1},
		{uStr: "9223372036854775808", index: 63, expectedBit: 1},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s[%d]", tc.uStr, tc.index), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)

			actualBit := u.Bit(tc.index)
			if actualBit != tc.expectedBit {
				t.Errorf("Bit(%s, %d) = %v, want %v", u.String(), tc.index, actualBit, tc.expectedBit)
			}
		})
	}
}

func TestUint128_Lsh(t *testing.T) {
	testCases := []struct {
		uStr        string
		n           uint
		expectedStr string
	}{
		{uStr: "1", n: 1, expectedStr: "2"},
		{uStr: "9223372036854775808", n: 1, expectedStr: "18446744073709551616"},
		{uStr: "1", n: 64, expectedStr: "18446744073709551616"},
		{uStr: "18446744073709551615", n: 1, expectedStr: "36893488147419103230"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s Lsh %d", tc.uStr, tc.n), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			expected, _ := U128FromString(tc.expectedStr)

			actual := u.Lsh(tc.n)
			if actual != expected {
				t.Errorf("Lsh(%s, %d) = %v, want %v", u.String(), tc.n, actual.String(), expected.String())
			}
		})
	}
}

func TestUint128_Cmp(t *testing.T) {
	testCases := []struct {
		uStr     string
		vStr     string
		expected int
	}{
		{uStr: "100", vStr: "99", expected: 1},
		{uStr: "99", vStr: "100", expected: -1},
		{uStr: "100", vStr: "100", expected: 0},
		{uStr: "18446744073709551615", vStr: "9223372036854775808", expected: 1},
		{uStr: "9223372036854775808", vStr: "18446744073709551615", expected: -1},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s Cmp %s", tc.uStr, tc.vStr), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			v, _ := U128FromString(tc.vStr)

			actual := u.Cmp(v)
			if actual != tc.expected {
				t.Errorf("Cmp(%s, %s) = %v, want %v", u.String(), v.String(), actual, tc.expected)
			}
		})
	}
}

func TestUint128_Mod(t *testing.T) {
	testCases := []struct {
		uStr        string
		vStr        string
		expectedStr string
		expectedErr bool
	}{
		{uStr: "10", vStr: "3", expectedStr: "1", expectedErr: false},
		{uStr: "18446744073709551616", vStr: "2", expectedStr: "0", expectedErr: false},
		{uStr: "1", vStr: "1", expectedStr: "0", expectedErr: false},
		{uStr: "1", vStr: "0", expectedStr: "", expectedErr: true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s %% %s", tc.uStr, tc.vStr), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			v, _ := U128FromString(tc.vStr)
			expected, _ := U128FromString(tc.expectedStr)

			actual, err := u.Mod(v)
			if (err != nil) != tc.expectedErr {
				t.Errorf("Mod(%s, %s) error = %v, wantErr %v", u.String(), v.String(), err, tc.expectedErr)
				return
			}
			if actual != expected {
				t.Errorf("Mod(%s, %s) = %v, want %v", u.String(), v.String(), actual.String(), expected.String())
			}
		})
	}
}

func TestUint128_And(t *testing.T) {
	testCases := []struct {
		uStr        string
		vStr        string
		expectedStr string
	}{
		{uStr: "5", vStr: "3", expectedStr: "1"},
		{uStr: "18446744073709551616", vStr: "1", expectedStr: "0"},
		{uStr: "255", vStr: "15", expectedStr: "15"},
		{uStr: "9223372036854775808", vStr: "1", expectedStr: "0"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s & %s", tc.uStr, tc.vStr), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			v, _ := U128FromString(tc.vStr)
			expected, _ := U128FromString(tc.expectedStr)

			actual := u.And(v)
			if actual != expected {
				t.Errorf("And(%s, %s) = %v, want %v", u.String(), v.String(), actual.String(), expected.String())
			}
		})
	}
}

func TestUint128_Or(t *testing.T) {
	testCases := []struct {
		uStr        string
		vStr        string
		expectedStr string
	}{
		{uStr: "5", vStr: "3", expectedStr: "7"},
		{uStr: "18446744073709551616", vStr: "1", expectedStr: "18446744073709551617"},
		{uStr: "255", vStr: "15", expectedStr: "255"},
		{uStr: "9223372036854775808", vStr: "1", expectedStr: "9223372036854775809"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s | %s", tc.uStr, tc.vStr), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			v, _ := U128FromString(tc.vStr)
			expected, _ := U128FromString(tc.expectedStr)

			actual := u.Or(v)
			if actual != expected {
				t.Errorf("Or(%s, %s) = %v, want %v", u.String(), v.String(), actual.String(), expected.String())
			}
		})
	}
}

func TestUint128_Xor(t *testing.T) {
	testCases := []struct {
		uStr        string
		vStr        string
		expectedStr string
	}{
		{uStr: "5", vStr: "3", expectedStr: "6"},
		{uStr: "18446744073709551616", vStr: "1", expectedStr: "18446744073709551617"},
		{uStr: "255", vStr: "15", expectedStr: "240"},
		{uStr: "9223372036854775808", vStr: "1", expectedStr: "9223372036854775809"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s ^ %s", tc.uStr, tc.vStr), func(t *testing.T) {
			u, _ := U128FromString(tc.uStr)
			v, _ := U128FromString(tc.vStr)
			expected, _ := U128FromString(tc.expectedStr)

			actual := u.Xor(v)
			if actual != expected {
				t.Errorf("Xor(%s, %s) = %v, want %v", u.String(), v.String(), actual.String(), expected.String())
			}
		})
	}
}

// Uint128 Math tests
