package main

import (
	"github.com/vlmoon99/near-sdk-go/collections"
)

type Counter struct {
	count  int
	owner  string
	values *collections.UnorderedMap[string, int]
}

func NewCounter() *Counter {
	return &Counter{
		count:  0,
		owner:  "",
		values: collections.NewUnorderedMap[string, int]("values"),
	}
}

// @contract:public
// @contract:view
func (c *Counter) GetCount() int {
	return c.count
}

// @contract:public
// @contract:payable min_deposit=0.0001
// @contract:mutating
func (c *Counter) AddValue(key string, value int) string {
	err := c.values.Insert(key, value)
	if err != nil {
		return "Error: " + err.Error()
	}
	return "Value added successfully"
}

// @contract:private
func (c *Counter) internalHelper() {
	// Private logic
}