package main

import (
	"encoding/binary"
	"encoding/hex"
	"strings"
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

func hexToDecimal(hex string) string {
	// Remove any leading zeros
	hex = strings.TrimLeft(hex, "0")
	if hex == "" {
		return "0"
	}

	// Initialize the decimal result as "0"
	decimal := "0"

	// Iterate over each character in the hex string
	for _, char := range hex {
		// Convert the hex character to its decimal value
		var value int
		if char >= '0' && char <= '9' {
			value = int(char - '0')
		} else if char >= 'a' && char <= 'f' {
			value = int(char - 'a' + 10)
		} else if char >= 'A' && char <= 'F' {
			value = int(char - 'A' + 10)
		} else {
			// Handle invalid characters (though the input should be valid hex)
			panic("invalid hex character")
		}

		// Multiply the current decimal result by 16 and add the new value
		decimal = multiplyBy16(decimal)
		decimal = addDecimal(decimal, value)
	}

	return decimal
}

func multiplyBy16(decimal string) string {
	result := ""
	carry := 0

	// Iterate over the decimal string from right to left
	for i := len(decimal) - 1; i >= 0; i-- {
		digit := int(decimal[i] - '0')
		product := digit*16 + carry
		carry = product / 10
		result = string('0'+(product%10)) + result
	}

	// If there's a carry left, add it to the result
	if carry > 0 {
		result = string('0'+carry) + result
	}

	return result
}

func addDecimal(decimal string, value int) string {
	result := ""
	carry := value

	// Iterate over the decimal string from right to left
	for i := len(decimal) - 1; i >= 0; i-- {
		digit := int(decimal[i] - '0')
		sum := digit + carry
		carry = sum / 10
		result = string('0'+(sum%10)) + result
	}

	// If there's a carry left, add it to the result
	if carry > 0 {
		result = string('0'+carry) + result
	}

	return result
}
