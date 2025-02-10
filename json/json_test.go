package json

import (
	"testing"
)

func TestBuilder_AddString(t *testing.T) {
	builder := NewBuilder()
	expected := `{"name":"John"}`
	result := string(builder.AddString("name", "John").Build())
	if result != expected {
		t.Errorf("AddString failed, expected %s, got %s", expected, result)
	}
}

func TestBuilder_AddInt(t *testing.T) {
	builder := NewBuilder()
	expected := `{"age":30}`
	result := string(builder.AddInt("age", 30).Build())
	if result != expected {
		t.Errorf("AddInt failed, expected %s, got %s", expected, result)
	}
}

func TestBuilder_AddInt64(t *testing.T) {
	builder := NewBuilder()
	expected := `{"population":10000000000}`
	result := string(builder.AddInt64("population", 10000000000).Build())
	if result != expected {
		t.Errorf("AddInt64 failed, expected %s, got %s", expected, result)
	}
}

func TestBuilder_AddUint64(t *testing.T) {
	builder := NewBuilder()
	expected := `{"total":18446744073709551615}`
	result := string(builder.AddUint64("total", 18446744073709551615).Build())
	if result != expected {
		t.Errorf("AddUint64 failed, expected %s, got %s", expected, result)
	}
}

func TestBuilder_AddFloat64(t *testing.T) {
	builder := NewBuilder()
	expected := `{"pi":3.141592653589793}`
	result := string(builder.AddFloat64("pi", 3.141592653589793).Build())
	if result != expected {
		t.Errorf("AddFloat64 failed, expected %s, got %s", expected, result)
	}
}

func TestBuilder_AddBool(t *testing.T) {
	builder := NewBuilder()
	expected := `{"valid":true}`
	result := string(builder.AddBool("valid", true).Build())
	if result != expected {
		t.Errorf("AddBool failed, expected %s, got %s", expected, result)
	}
}

func TestParser_GetRawBytes(t *testing.T) {
	parser := NewParser([]byte(`{"data":"example data"}`))
	expected := []byte("example data")
	result, err := parser.GetRawBytes("data")
	if err != nil || string(result) != string(expected) {
		t.Errorf("GetRawBytes failed, expected %s, got %s, err %v", expected, result, err)
	}
}

func TestParser_GetString(t *testing.T) {
	parser := NewParser([]byte(`{"name":"John"}`))
	expected := "John"
	result, err := parser.GetString("name")
	if err != nil || result != expected {
		t.Errorf("GetString failed, expected %s, got %s, err %v", expected, result, err)
	}
}

func TestParser_GetInt(t *testing.T) {
	parser := NewParser([]byte(`{"age":30}`))
	expected := int64(30)
	result, err := parser.GetInt("age")
	if err != nil || result != expected {
		t.Errorf("GetInt failed, expected %d, got %d, err %v", expected, result, err)
	}
}

func TestParser_GetFloat64(t *testing.T) {
	parser := NewParser([]byte(`{"pi":3.141592653589793}`))
	expected := 3.141592653589793
	result, err := parser.GetFloat64("pi")
	if err != nil || result != expected {
		t.Errorf("GetFloat64 failed, expected %f, got %f, err %v", expected, result, err)
	}
}

func TestParser_GetBoolean(t *testing.T) {
	parser := NewParser([]byte(`{"valid":true}`))
	expected := true
	result, err := parser.GetBoolean("valid")
	if err != nil || result != expected {
		t.Errorf("GetBoolean failed, expected %v, got %v, err %v", expected, result, err)
	}
}

func TestParser_GetBytes(t *testing.T) {
	parser := NewParser([]byte(`{"data":"example data"}`))
	expected := []byte("example data")
	result, err := parser.GetBytes("data")
	if err != nil || string(result) != string(expected) {
		t.Errorf("GetBytes failed, expected %s, got %s, err %v", expected, result, err)
	}
}
