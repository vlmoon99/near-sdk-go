// Package borsh provides functions and types for serializing and deserializing data
// using Binary Object Representation Serialization (BORSH).
package borsh

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"

	"github.com/vlmoon99/near-sdk-go/types"
)

const (
	ErrEOF             = "(BORSH_ERROR): EOF"
	ErrPointerRequired = "(BORSH_ERROR): passed struct must be pointer"
	ErrUnsupportedType = "(BORSH_ERROR): unsupported type: "
)

// The ByteReader type is used to read bytes from a byte slice.
type ByteReader struct {
	data []byte
	pos  int
}

// NewByteReader creates and returns a new ByteReader instance.
//
// Parameters:
//
//	data: The byte slice to read from.
func NewByteReader(data []byte) *ByteReader {
	return &ByteReader{data: data}
}

// Read reads up to len(p) bytes into p.
//
// Parameters:
//
//	p: The byte slice to read into.
//
// Returns:
//
//	n: The number of bytes read.
//	err: An error if reading fails.
func (r *ByteReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, errors.New(ErrEOF)
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// ReadByte reads and returns a single byte.
//
// Returns:
//
//	b: The byte read.
//	err: An error if reading fails.
func (r *ByteReader) ReadByte() (byte, error) {
	if r.pos >= len(r.data) {
		return 0, errors.New(ErrEOF)
	}
	b := r.data[r.pos]
	r.pos++
	return b, nil
}

// The ByteWriter type is used to write bytes to a buffer.
type ByteWriter struct {
	buf []byte
}

// NewByteWriter creates and returns a new ByteWriter instance.
func NewByteWriter() *ByteWriter {
	return &ByteWriter{}
}

// Write writes the byte slice p to the buffer.
//
// Parameters:
//
//	p: The byte slice to write.
//
// Returns:
//
//	n: The number of bytes written.
//	err: An error if writing fails.
func (w *ByteWriter) Write(p []byte) (n int, err error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}

// WriteByte writes a single byte to the buffer.
//
// Parameters:
//
//	c: The byte to write.
//
// Returns:
//
//	err: An error if writing fails.
func (w *ByteWriter) WriteByte(c byte) error {
	w.buf = append(w.buf, c)
	return nil
}

// Bytes returns the buffer as a byte slice.
func (w *ByteWriter) Bytes() []byte {
	return w.buf
}

// Deserialize deserializes the byte slice data into the provided struct.
//
// Parameters:
//
//	data: The byte slice to deserialize.
//	s: The struct to deserialize into.
//
// Returns:
//
//	err: An error if deserialization fails.
func Deserialize(data []byte, s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return errors.New(ErrPointerRequired)
	}

	r := NewByteReader(data)
	result, err := deserialize(reflect.TypeOf(s).Elem(), r)
	if err != nil {
		return err
	}

	v.Elem().Set(reflect.ValueOf(result))
	return nil
}

// Serialize serializes the provided struct into a byte slice.
//
// Parameters:
//
//	s: The struct to serialize.
//
// Returns:
//
//	b: The serialized byte slice.
//	err: An error if serialization fails.
func Serialize(s interface{}) ([]byte, error) {
	b := NewByteWriter()
	err := serialize(reflect.ValueOf(s), b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// deserialize deserializes the supported types.
//
// Parameters:
//
//	t: The type to deserialize.
//	r: The ByteReader instance to read the serialized data.
//
// Returns:
//
//	v: The deserialized value.
//	err: An error if deserialization fails.
func deserialize(t reflect.Type, r *ByteReader) (interface{}, error) {
	switch t.Kind() {
	case reflect.Bool:
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		return b == 1, nil
	case reflect.Int8:
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		return int8(b), nil
	case reflect.Int16:
		tmp := make([]byte, 2)
		_, err := r.Read(tmp)
		if err != nil {
			return nil, err
		}
		return int16(binary.LittleEndian.Uint16(tmp)), nil
	case reflect.Int32:
		tmp := make([]byte, 4)
		_, err := r.Read(tmp)
		if err != nil {
			return nil, err
		}
		return int32(binary.LittleEndian.Uint32(tmp)), nil
	case reflect.Int64:
		tmp := make([]byte, 8)
		_, err := r.Read(tmp)
		if err != nil {
			return nil, err
		}
		return int64(binary.LittleEndian.Uint64(tmp)), nil
	case reflect.Int:
		tmp := make([]byte, 8)
		_, err := r.Read(tmp)
		if err != nil {
			return nil, err
		}
		return int(binary.LittleEndian.Uint64(tmp)), nil
	case reflect.Uint8:
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		return b, nil
	case reflect.Uint16:
		tmp := make([]byte, 2)
		_, err := r.Read(tmp)
		if err != nil {
			return nil, err
		}
		return uint16(binary.LittleEndian.Uint16(tmp)), nil
	case reflect.Uint32:
		tmp := make([]byte, 4)
		_, err := r.Read(tmp)
		if err != nil {
			return nil, err
		}
		return uint32(binary.LittleEndian.Uint32(tmp)), nil
	case reflect.Uint64, reflect.Uint:
		tmp := make([]byte, 8)
		_, err := r.Read(tmp)
		if err != nil {
			return nil, err
		}
		return uint64(binary.LittleEndian.Uint64(tmp)), nil
	case reflect.Float32:
		tmp := make([]byte, 4)
		_, err := r.Read(tmp)
		if err != nil {
			return nil, err
		}
		return math.Float32frombits(binary.LittleEndian.Uint32(tmp)), nil
	case reflect.Float64:
		tmp := make([]byte, 8)
		_, err := r.Read(tmp)
		if err != nil {
			return nil, err
		}
		return math.Float64frombits(binary.LittleEndian.Uint64(tmp)), nil
	case reflect.Complex64:
		tmp := make([]byte, 8)
		_, err := r.Read(tmp)
		if err != nil {
			return nil, err
		}
		realPart := math.Float32frombits(binary.LittleEndian.Uint32(tmp))
		imagPart := math.Float32frombits(binary.LittleEndian.Uint32(tmp[4:]))
		return complex(realPart, imagPart), nil
	case reflect.Complex128:
		tmp := make([]byte, 16)
		_, err := r.Read(tmp)
		if err != nil {
			return nil, err
		}
		realPart := math.Float64frombits(binary.LittleEndian.Uint64(tmp))
		imagPart := math.Float64frombits(binary.LittleEndian.Uint64(tmp[8:]))
		return complex(realPart, imagPart), nil
	case reflect.String:
		lenBytes := make([]byte, 4)
		_, err := r.Read(lenBytes)
		if err != nil {
			return nil, err
		}
		strLen := int(binary.LittleEndian.Uint32(lenBytes))
		strBytes := make([]byte, strLen)
		_, err = r.Read(strBytes)
		if err != nil {
			return nil, err
		}
		return string(strBytes), nil
	case reflect.Struct:
		if t == reflect.TypeOf(types.Uint128{}) {
			data := make([]byte, 16)
			_, err := r.Read(data)
			if err != nil {
				return nil, err
			}
			u128, err := types.LoadUint128LE(data)
			if err != nil {
				return nil, err
			}
			return u128, nil
		}
		result := reflect.New(t).Elem()
		for i := 0; i < t.NumField(); i++ {
			field, err := deserialize(t.Field(i).Type, r)
			if err != nil {
				return nil, err
			}
			result.Field(i).Set(reflect.ValueOf(field))
		}
		return result.Interface(), nil
	case reflect.Slice:
		lenBytes := make([]byte, 4)
		_, err := r.Read(lenBytes)
		if err != nil {
			return nil, err
		}
		sliceLen := int(binary.LittleEndian.Uint32(lenBytes))
		slice := reflect.MakeSlice(t, sliceLen, sliceLen)
		for i := 0; i < sliceLen; i++ {
			elem, err := deserialize(t.Elem(), r)
			if err != nil {
				return nil, err
			}
			slice.Index(i).Set(reflect.ValueOf(elem))
		}
		return slice.Interface(), nil
	case reflect.Array:
		array := reflect.New(t).Elem()
		for i := 0; i < t.Len(); i++ {
			elem, err := deserialize(t.Elem(), r)
			if err != nil {
				return nil, err
			}
			array.Index(i).Set(reflect.ValueOf(elem))
		}
		return array.Interface(), nil
	case reflect.Map:
		lenBytes := make([]byte, 4)
		_, err := r.Read(lenBytes)
		if err != nil {
			return nil, err
		}
		mapLen := int(binary.LittleEndian.Uint32(lenBytes))
		mapType := reflect.MakeMapWithSize(t, mapLen)
		for i := 0; i < mapLen; i++ {
			key, err := deserialize(t.Key(), r)
			if err != nil {
				return nil, err
			}
			value, err := deserialize(t.Elem(), r)
			if err != nil {
				return nil, err
			}
			mapType.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
		}
		return mapType.Interface(), nil
	case reflect.Ptr:
		tmp, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if tmp == 0 {
			return reflect.Zero(t).Interface(), nil
		}
		elem, err := deserialize(t.Elem(), r)
		if err != nil {
			return nil, err
		}
		ptr := reflect.New(t.Elem())
		ptr.Elem().Set(reflect.ValueOf(elem))
		return ptr.Interface(), nil
	default:
		return nil, errors.New(ErrUnsupportedType + t.Name())
	}
}

// serialize serializes the supported types.
//
// Parameters:
//
//	v: The value to serialize.
//	b: The ByteWriter instance to write the serialized data.
//
// Returns:
//
//	err: An error if serialization fails.
func serialize(v reflect.Value, b *ByteWriter) error {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return b.WriteByte(1)
		}
		return b.WriteByte(0)
	case reflect.Int8:
		return b.WriteByte(byte(v.Int()))
	case reflect.Int16:
		tmp := make([]byte, 2)
		binary.LittleEndian.PutUint16(tmp, uint16(v.Int()))
		_, err := b.Write(tmp)
		return err
	case reflect.Int32:
		tmp := make([]byte, 4)
		binary.LittleEndian.PutUint32(tmp, uint32(v.Int()))
		_, err := b.Write(tmp)
		return err
	case reflect.Int64:
		tmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(tmp, uint64(v.Int()))
		_, err := b.Write(tmp)
		return err
	case reflect.Int:
		tmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(tmp, uint64(v.Int()))
		_, err := b.Write(tmp)
		return err
	case reflect.Uint8:
		return b.WriteByte(byte(v.Uint()))
	case reflect.Uint16:
		tmp := make([]byte, 2)
		binary.LittleEndian.PutUint16(tmp, uint16(v.Uint()))
		_, err := b.Write(tmp)
		return err
	case reflect.Uint32:
		tmp := make([]byte, 4)
		binary.LittleEndian.PutUint32(tmp, uint32(v.Uint()))
		_, err := b.Write(tmp)
		return err
	case reflect.Uint64, reflect.Uint:
		tmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(tmp, v.Uint())
		_, err := b.Write(tmp)
		return err
	case reflect.Float32:
		tmp := make([]byte, 4)
		binary.LittleEndian.PutUint32(tmp, math.Float32bits(float32(v.Float())))
		_, err := b.Write(tmp)
		return err
	case reflect.Float64:
		tmp := make([]byte, 8)
		binary.LittleEndian.PutUint64(tmp, math.Float64bits(v.Float()))
		_, err := b.Write(tmp)
		return err
	case reflect.Complex64:
		tmp := make([]byte, 8)
		binary.LittleEndian.PutUint32(tmp, math.Float32bits(float32(real(v.Complex()))))
		binary.LittleEndian.PutUint32(tmp[4:], math.Float32bits(float32(imag(v.Complex()))))
		_, err := b.Write(tmp)
		return err
	case reflect.Complex128:
		tmp := make([]byte, 16)
		binary.LittleEndian.PutUint64(tmp, math.Float64bits(real(v.Complex())))
		binary.LittleEndian.PutUint64(tmp[8:], math.Float64bits(imag(v.Complex())))
		_, err := b.Write(tmp)
		return err
	case reflect.String:
		str := v.String()
		lenBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(lenBytes, uint32(len(str)))
		_, err := b.Write(lenBytes)
		if err != nil {
			return err
		}
		_, err = b.Write([]byte(str))
		return err
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(types.Uint128{}) {
			u128 := v.Interface().(types.Uint128)
			leBytes := u128.ToLE()
			_, err := b.Write(leBytes)
			return err
		}
		for i := 0; i < v.NumField(); i++ {
			err := serialize(v.Field(i), b)
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Slice:
		lenBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(lenBytes, uint32(v.Len()))
		_, err := b.Write(lenBytes)
		if err != nil {
			return err
		}
		for i := 0; i < v.Len(); i++ {
			err := serialize(v.Index(i), b)
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			err := serialize(v.Index(i), b)
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		lenBytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(lenBytes, uint32(v.Len()))
		_, err := b.Write(lenBytes)
		if err != nil {
			return err
		}
		for _, key := range v.MapKeys() {
			err := serialize(key, b)
			if err != nil {
				return err
			}
			err = serialize(v.MapIndex(key), b)
			if err != nil {
				return err
			}
		}
		return nil
	case reflect.Ptr:
		if v.IsNil() {
			return b.WriteByte(0)
		}
		err := b.WriteByte(1)
		if err != nil {
			return err
		}
		return serialize(v.Elem(), b)
	default:
		return errors.New(ErrUnsupportedType + v.Type().String())
	}
}
