package borsh

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vlmoon99/near-sdk-go/types"
)

func TestSerializeDeserializeBool(t *testing.T) {
	value := true

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized bool
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeInt8(t *testing.T) {
	value := int8(127)

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized int8
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeInt16(t *testing.T) {
	value := int16(32767)

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized int16
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeInt32(t *testing.T) {
	value := int32(2147483647)

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized int32
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeInt64(t *testing.T) {
	value := int64(9223372036854775807)

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized int64
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeInt(t *testing.T) {
	value := 123456

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized int
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeUint8(t *testing.T) {
	value := uint8(255)

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized uint8
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeUint16(t *testing.T) {
	value := uint16(65535)

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized uint16
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeUint32(t *testing.T) {
	value := uint32(4294967295)

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized uint32
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeUint64(t *testing.T) {
	value := uint64(18446744073709551615)

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized uint64
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeFloat32(t *testing.T) {
	value := float32(123.456)

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized float32
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeFloat64(t *testing.T) {
	value := float64(123.456789)

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized float64
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeComplex64(t *testing.T) {
	value := complex(float32(1.5), float32(2.5))

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized complex64
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeComplex128(t *testing.T) {
	value := complex(1.5, 2.5)

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized complex128
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeRune(t *testing.T) {
	value := rune(2147483647)

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized rune
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeByte(t *testing.T) {
	value := byte(255)

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized byte
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeString(t *testing.T) {
	value := "Hello, Borsh!"

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized string
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

// func TestSerializeDeserializeUintptr(t *testing.T) {
// 	value := uintptr(123456)

// 	serialized, err := Serialize(value)
// 	if err != nil {
// 		t.Fatalf("Serialization failed: %v", err)
// 	}

// 	var deserialized uintptr
// 	err = Deserialize(serialized, &deserialized)
// 	if err != nil {
// 		t.Fatalf("Deserialization failed: %v", err)
// 	}

// 	if value != deserialized {
// 		t.Fatalf("Expected %v, but got %v", value, deserialized)
// 	}
// }

func TestSerializeDeserializeUint128(t *testing.T) {
	value, _ := types.U128FromString("123456789012345678901234567890")

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized types.Uint128
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if !bytes.Equal(value.ToLE(), deserialized.ToLE()) {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

type SimpleStruct struct {
	BoolField   bool
	IntField    int
	StringField string
}

func TestSerializeDeserializeSimpleStruct(t *testing.T) {
	value := SimpleStruct{
		BoolField:   true,
		IntField:    123456,
		StringField: "Hello, Struct!",
	}

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized SimpleStruct
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

type InnerStruct struct {
	IntField int
}

type OuterStruct struct {
	BoolField   bool
	Inner       InnerStruct
	StringField string
}

func TestSerializeDeserializeOuterStruct(t *testing.T) {
	value := OuterStruct{
		BoolField:   true,
		Inner:       InnerStruct{IntField: 654321},
		StringField: "Nested Structs",
	}

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized OuterStruct
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeSlice(t *testing.T) {
	value := []int{1, 2, 3, 4, 5}

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized []int
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if !reflect.DeepEqual(value, deserialized) {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeArray(t *testing.T) {
	value := [5]int{1, 2, 3, 4, 5}

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized [5]int
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value != deserialized {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

func TestSerializeDeserializeMap(t *testing.T) {
	value := map[string]int{"one": 1, "two": 2, "three": 3}

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized map[string]int
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if !reflect.DeepEqual(value, deserialized) {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}
}

type PointerStruct struct {
	BoolField   bool
	InnerPtr    *InnerStruct
	StringField string
}

func TestSerializeDeserializePointerStruct(t *testing.T) {
	innerValue := InnerStruct{IntField: 987654}
	value := PointerStruct{
		BoolField:   true,
		InnerPtr:    &innerValue,
		StringField: "Pointer to Struct",
	}

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized PointerStruct
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value.BoolField != deserialized.BoolField || value.StringField != deserialized.StringField {
		t.Fatalf("Expected %v, but got %v", value, deserialized)
	}

	if deserialized.InnerPtr == nil || !reflect.DeepEqual(value.InnerPtr, deserialized.InnerPtr) {
		t.Fatalf("Expected %v, but got %v", value.InnerPtr, deserialized.InnerPtr)
	}
}

type StructC struct {
	Field1 int
	Field2 string
}

type StructB struct {
	Field3 float64
	Field4 *StructC
}

type StructA struct {
	Field5  bool
	Field6  StructB
	Field7  *StructC
	Field8  []StructB
	Field9  map[string]*StructC
	Field10 *StructB
}

func TestSerializeDeserializeComplexStruct(t *testing.T) {
	innerStructC := StructC{Field1: 789, Field2: "Inner Struct C"}
	innerStructC2 := StructC{Field1: 101112, Field2: "Another Inner Struct C"}

	value := StructA{
		Field5:  true,
		Field6:  StructB{Field3: 3.14, Field4: &innerStructC},
		Field7:  &innerStructC,
		Field8:  []StructB{{Field3: 6.28, Field4: &innerStructC}, {Field3: 9.42, Field4: &innerStructC2}},
		Field9:  map[string]*StructC{"first": &innerStructC, "second": &innerStructC2},
		Field10: &StructB{Field3: 1.23, Field4: &innerStructC2},
	}

	serialized, err := Serialize(value)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var deserialized StructA
	err = Deserialize(serialized, &deserialized)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if value.Field5 != deserialized.Field5 {
		t.Fatalf("Expected %v, but got %v", value.Field5, deserialized.Field5)
	}

	if !reflect.DeepEqual(value.Field6, deserialized.Field6) {
		t.Fatalf("Expected %v, but got %v", value.Field6, deserialized.Field6)
	}

	if !reflect.DeepEqual(value.Field7, deserialized.Field7) {
		t.Fatalf("Expected %v, but got %v", value.Field7, deserialized.Field7)
	}

	if !reflect.DeepEqual(value.Field8, deserialized.Field8) {
		t.Fatalf("Expected %v, but got %v", value.Field8, deserialized.Field8)
	}

	if !reflect.DeepEqual(value.Field9, deserialized.Field9) {
		t.Fatalf("Expected %v, but got %v", value.Field9, deserialized.Field9)
	}

	if !reflect.DeepEqual(value.Field10, deserialized.Field10) {
		t.Fatalf("Expected %v, but got %v", value.Field10, deserialized.Field10)
	}
}
