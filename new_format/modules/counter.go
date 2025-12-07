package main

// @contract:public
// @contract:view
func (c *Counter) GetCount() int {
	return c.count
}

// @contract:public
// @contract:mutating
func (c *Counter) Increment(amount int) int {
	c.count += amount
	return c.count
}

// @contract:public
// @contract:mutating
func (c *Counter) Decrement(amount int) int {
	c.count -= amount
	return c.count
}

// @contract:public
// @contract:mutating
// @contract:payable min_deposit=0.001
func (c *Counter) Reset() string {
	c.count = 0
	return "Counter reset"
}