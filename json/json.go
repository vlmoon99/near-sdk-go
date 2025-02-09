package json

import (
	"errors"

	"github.com/vlmoon99/jsonparser"
)

type Builder struct {
	data []byte
}

func NewBuilder() *Builder {
	return &Builder{data: []byte{'{'}}
}

func (b *Builder) AddString(key, value string) *Builder {
	b.data = append(b.data, '"')
	b.data = append(b.data, key...)
	b.data = append(b.data, '"', ':')
	b.data = append(b.data, '"')
	b.data = append(b.data, value...)
	b.data = append(b.data, '"', ',')
	return b
}

func (b *Builder) AddInt(key string, value int) *Builder {
	b.data = append(b.data, '"')
	b.data = append(b.data, key...)
	b.data = append(b.data, '"', ':')
	b.data = append(b.data, intToBytes(value)...)
	b.data = append(b.data, ',')
	return b
}

func intToBytes(n int) []byte {
	if n == 0 {
		return []byte("0")
	}

	length := 0
	for temp := n; temp != 0; temp /= 10 {
		length++
	}

	result := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		result[i] = byte('0' + n%10)
		n /= 10
	}
	return result
}

func (b *Builder) Build() []byte {
	if len(b.data) > 1 {
		b.data[len(b.data)-1] = '}'
	} else {
		b.data = append(b.data, '}')
	}
	return b.data
}

type Parser struct {
	data []byte
}

func NewParser(data []byte) *Parser {
	return &Parser{data: data}
}

func (p *Parser) GetRawBytes(key string) ([]byte, error) {
	data, _, _, err := jsonparser.Get(p.data, key)
	if err != nil {
		return nil, errors.New("Error while getting raw bytes from the json")
	}
	return data, nil
}

func (p *Parser) GetString(key string) (string, error) {
	return jsonparser.GetString(p.data, key)
}

func (p *Parser) GetInt(key string) (int64, error) {
	return jsonparser.GetInt(p.data, key)
}

func (p *Parser) GetBoolean(key string) (bool, error) {
	return jsonparser.GetBoolean(p.data, key)
}

func (p *Parser) GetFloat(key string) (float64, error) {
	return jsonparser.GetFloat(p.data, key)
}
