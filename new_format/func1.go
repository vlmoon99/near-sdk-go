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
