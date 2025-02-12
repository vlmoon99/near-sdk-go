// Package json represents a simple package for handling JSON parsing and encoding.
package json

import (
	"errors"
	"strconv"

	"github.com/vlmoon99/jsonparser"
)

const (
	ErrGettingRawBytesFromJson = "(JSON_ERROR): error while getting raw bytes from the json"
)

// The Builder type is used to create a JSON object by adding key-value pairs of different types.
type Builder struct {
	data []byte
}

// Creates and returns a new Builder instance.
func NewBuilder() *Builder {
	return &Builder{data: []byte{'{'}}
}

// Adds a string key-value pair to the JSON object.
//
// Parameters:
//
//	key: The key of the JSON element.
//	value: The string value of the JSON element.
func (b *Builder) AddString(key, value string) *Builder {
	return b.addKey(key).addValue(`"` + value + `"`).addComma()
}

// Adds an integer key-value pair to the JSON object.
//
// Parameters:
//
//	key: The key of the JSON element.
//	value: The integer value of the JSON element.
func (b *Builder) AddInt(key string, value int) *Builder {
	return b.addKey(key).addValue(strconv.Itoa(value)).addComma()
}

// Adds a 64-bit integer key-value pair to the JSON object.
//
// Parameters:
//
//	key: The key of the JSON element.
//	value: The 64-bit integer value of the JSON element.
func (b *Builder) AddInt64(key string, value int64) *Builder {
	return b.addKey(key).addValue(strconv.FormatInt(value, 10)).addComma()
}

// Adds a 64-bit unsigned integer key-value pair to the JSON object.
//
// Parameters:
//
//	key: The key of the JSON element.
//	value: The 64-bit unsigned integer value of the JSON element.
func (b *Builder) AddUint64(key string, value uint64) *Builder {
	return b.addKey(key).addValue(strconv.FormatUint(value, 10)).addComma()
}

// Adds a 64-bit floating-point key-value pair to the JSON object.
//
// Parameters:
//
//	key: The key of the JSON element.
//	value: The 64-bit floating-point value of the JSON element.
func (b *Builder) AddFloat64(key string, value float64) *Builder {
	return b.addKey(key).addValue(strconv.FormatFloat(value, 'f', -1, 64)).addComma()
}

// Adds a boolean key-value pair to the JSON object.
//
// Parameters:
//
//	key: The key of the JSON element.
//	value: The boolean value of the JSON element.
func (b *Builder) AddBool(key string, value bool) *Builder {
	return b.addKey(key).addValue(strconv.FormatBool(value)).addComma()
}

func (b *Builder) addKey(key string) *Builder {
	b.data = append(b.data, '"')
	b.data = append(b.data, key...)
	b.data = append(b.data, '"', ':')
	return b
}

func (b *Builder) addValue(value string) *Builder {
	b.data = append(b.data, value...)
	return b
}

func (b *Builder) addComma() *Builder {
	b.data = append(b.data, ',')
	return b
}

// Finalizes and returns the JSON object as a byte slice.
func (b *Builder) Build() []byte {
	if len(b.data) > 1 {
		b.data[len(b.data)-1] = '}'
	} else {
		b.data = append(b.data, '}')
	}
	return b.data
}

// The Parser type is used to parse a JSON object and retrieve values of different types by their keys.
type Parser struct {
	data []byte
}

// Creates and returns a new Parser instance with the provided JSON data.
//
// Parameters:
//
//	data: The JSON data as a byte slice.
func NewParser(data []byte) *Parser {
	return &Parser{data: data}
}

// Retrieves the raw byte slice associated with the specified key.
//
// Parameters:
//
//	key: The key of the JSON element.
//
// Returns:
//
//	[]byte: The raw byte slice of the JSON element.
//	error: An error if the key is not found.
func (p *Parser) GetRawBytes(key string) ([]byte, error) {
	data, _, _, err := jsonparser.Get(p.data, key)
	if err != nil {
		return nil, errors.New(ErrGettingRawBytesFromJson)
	}
	return data, nil
}

// Retrieves the string value associated with the specified key.
//
// Parameters:
//
//	key: The key of the JSON element.
//
// Returns:
//
//	string: The string value of the JSON element.
//	error: An error if the key is not found or the value is not a string.
func (p *Parser) GetString(key string) (string, error) {
	return jsonparser.GetString(p.data, key)
}

// Retrieves the 64-bit integer value associated with the specified key.
//
// Parameters:
//
//	key: The key of the JSON element.
//
// Returns:
//
//	int64: The 64-bit integer value of the JSON element.
//	error: An error if the key is not found or the value is not an integer.
func (p *Parser) GetInt(key string) (int64, error) {
	return jsonparser.GetInt(p.data, key)
}

// Retrieves the 64-bit floating-point value associated with the specified key.
//
// Parameters:
//
//	key: The key of the JSON element.
//
// Returns:
//
//	float64: The 64-bit floating-point value of the JSON element.
//	error: An error if the key is not found or the value is not a floating-point number.
func (p *Parser) GetFloat64(key string) (float64, error) {
	return jsonparser.GetFloat(p.data, key)
}

// Retrieves the boolean value associated with the specified key.
//
// Parameters:
//
//	key: The key of the JSON element.
//
// Returns:
//
//	bool: The boolean value of the JSON element.
//	error: An error if the key is not found or the value is not a boolean.
func (p *Parser) GetBoolean(key string) (bool, error) {
	return jsonparser.GetBoolean(p.data, key)
}

// Retrieves the byte slice associated with the specified key as a string.
//
// Parameters:
//
//	key: The key of the JSON element.
//
// Returns:
//
//	[]byte: The byte slice of the JSON element.
//	error: An error if the key is not found or the value is not a string.
func (p *Parser) GetBytes(key string) ([]byte, error) {
	s, err := jsonparser.GetString(p.data, key)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}
