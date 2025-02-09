package borsh

import (
	"encoding/binary"
	"errors"
	"reflect"

	"github.com/vlmoon99/near-sdk-go/types"
)

type ByteReader struct {
	data []byte
	pos  int
}

func NewByteReader(data []byte) *ByteReader {
	return &ByteReader{data: data}
}

func (r *ByteReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, errors.New("EOF")
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

func (r *ByteReader) ReadByte() (byte, error) {
	if r.pos >= len(r.data) {
		return 0, errors.New("EOF")
	}
	b := r.data[r.pos]
	r.pos++
	return b, nil
}

type ByteWriter struct {
	buf []byte
}

func NewByteWriter() *ByteWriter {
	return &ByteWriter{}
}

func (w *ByteWriter) Write(p []byte) (n int, err error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}

func (w *ByteWriter) WriteByte(c byte) error {
	w.buf = append(w.buf, c)
	return nil
}

func (w *ByteWriter) Bytes() []byte {
	return w.buf
}

func Deserialize(data []byte, s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return errors.New("passed struct must be pointer")
	}

	r := NewByteReader(data)
	result, err := deserialize(reflect.TypeOf(s).Elem(), r)
	if err != nil {
		return err
	}

	v.Elem().Set(reflect.ValueOf(result))
	return nil
}

func Serialize(s interface{}) ([]byte, error) {
	b := NewByteWriter()
	err := serialize(reflect.ValueOf(s), b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

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
	case reflect.Uint64:
		tmp := make([]byte, 8)
		_, err := r.Read(tmp)
		if err != nil {
			return nil, err
		}
		return uint64(binary.LittleEndian.Uint64(tmp)), nil
	case reflect.Uint:
		tmp := make([]byte, 8)
		_, err := r.Read(tmp)
		if err != nil {
			return nil, err
		}
		return uint(binary.LittleEndian.Uint64(tmp)), nil
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
		} else {
			return nil, nil
		}
	case reflect.Array:
		return nil, nil
	case reflect.Slice:
		return nil, nil
	case reflect.Map:
		return nil, nil
	case reflect.Ptr:
		return nil, nil
	default:
		return nil, errors.New("unsupported type:" + t.Name())
	}
}

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
		} else {
			return nil
		}
	case reflect.Array:
		return nil
	case reflect.Slice:
		return nil
	case reflect.Map:
		return nil
	case reflect.Ptr:
		return nil
	default:
		return errors.New("unsupported type:" + v.Type().String())

	}
}
