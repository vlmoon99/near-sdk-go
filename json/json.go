package json

import (
	"errors"
	"strconv"

	"github.com/vlmoon99/jsonparser"
)

type Builder struct {
	data []byte
}

func NewBuilder() *Builder {
	return &Builder{data: []byte{'{'}}
}

func (b *Builder) AddString(key, value string) *Builder {
	return b.addKey(key).addValue(`"` + value + `"`).addComma()
}

func (b *Builder) AddInt(key string, value int) *Builder {
	return b.addKey(key).addValue(strconv.Itoa(value)).addComma()
}

func (b *Builder) AddInt64(key string, value int64) *Builder {
	return b.addKey(key).addValue(strconv.FormatInt(value, 10)).addComma()
}

func (b *Builder) AddUint64(key string, value uint64) *Builder {
	return b.addKey(key).addValue(strconv.FormatUint(value, 10)).addComma()
}

func (b *Builder) AddFloat64(key string, value float64) *Builder {
	return b.addKey(key).addValue(strconv.FormatFloat(value, 'f', -1, 64)).addComma() // 'f' format, -1 precision
}

func (b *Builder) AddBool(key string, value bool) *Builder {
	return b.addKey(key).addValue(strconv.FormatBool(value)).addComma()
}

func (b *Builder) AddBytes(key string, value []byte) *Builder {
	return b.addKey(key).addValue(string(value)).addComma() // Bytes are added as strings
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

func (b *Builder) Build() []byte {
	if len(b.data) > 1 {
		b.data[len(b.data)-1] = '}' // Replace last comma with closing brace
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

func (p *Parser) GetFloat64(key string) (float64, error) {
	return jsonparser.GetFloat(p.data, key)
}

func (p *Parser) GetBoolean(key string) (bool, error) {
	return jsonparser.GetBoolean(p.data, key)
}

func (p *Parser) GetBytes(key string) ([]byte, error) {
	s, err := jsonparser.GetString(p.data, key)
	if err != nil {
		return nil, err
	}
	return []byte(s), nil
}
