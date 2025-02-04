package main

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"
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

func FromString(str string) Uint128 {
	value, err := strconv.ParseFloat(str, 64)

	if err != nil {
		return Uint128{Hi: 0, Lo: 0}
	}

	hi := uint64(value / (1 << 64))
	lo := uint64(value)
	return Uint128{Hi: hi, Lo: lo}
}

func FromFloat64(value float64) Uint128 {
	hi := uint64(value / (1 << 64))
	lo := uint64(value)
	return Uint128{Hi: hi, Lo: lo}
}

func BoolToUnit(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}

func (u Uint128) ToFloat64() float64 {
	value := float64(u.Hi)*(1<<64) + float64(u.Lo)
	return value
}

func (u Uint128) ToYoctoNear() float64 {
	value := float64(u.Hi)*(1<<64) + float64(u.Lo)
	return value
}

//Uint128

// Near Gas
type NearGas struct {
	inner uint64
}

//Near Gas

//Serializer

// type Person struct {
// 	Name  string
// 	Age   int
// 	Email string
// }

// // SerializeToJSON converts the struct to a JSON byte slice
// func (p *Person) SerializeToJSON() []byte {
// 	// Manually construct the JSON string as a byte slice
// 	jsonData := []byte(`{"name":"`)
// 	jsonData = append(jsonData, p.Name...)
// 	jsonData = append(jsonData, `","age":`...)
// 	jsonData = append(jsonData, intToBytes(p.Age)...)
// 	jsonData = append(jsonData, `,"email":"`...)
// 	jsonData = append(jsonData, p.Email...)
// 	jsonData = append(jsonData, `"}`...)
// 	return jsonData
// }

// // DeserializeFromJSON populates the struct from a JSON byte slice
// func (p *Person) DeserializeFromJSON(data []byte) error {
// 	// Use jsonparser to extract fields from the JSON data
// 	name, err := GetString(data, "name")
// 	if err != nil {
// 		return err
// 	}
// 	p.Name = name

// 	age, err := GetInt(data, "age")
// 	if err != nil {
// 		return err
// 	}
// 	p.Age = int(age)

// 	email, err := GetString(data, "email")
// 	if err != nil {
// 		return err
// 	}
// 	p.Email = email

// 	return nil
// }

// // ToRawBytes converts the struct to raw bytes (JSON in this case)
// func (p *Person) ToRawBytes() []byte {
// 	return p.SerializeToJSON()
// }

// // Helper function to convert an integer to a byte slice
