package main

import (
	"github.com/vlmoon99/near-sdk-go/json"
	"github.com/vlmoon99/near-sdk-go/env"
)

func TestIdea() bool {
	builder := json.NewBuilder()
	expected := `{"age":30}`
	result := string(builder.AddInt("age", 30).Build())
	if result != expected {
		env.LogString("Error")
	}
	return true
}

type MyData struct {
    Name string
    Age  int
}

// @contract:public
func (c *Counter) ProcessData(data MyData) string {
    // Works!
	return data.Name
}
