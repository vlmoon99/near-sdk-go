package main

// @contract:public
// @contract:view
func (c *Counter) Test1() string {
	return ""
}

// @contract:public
// @contract:mutating
func (c *Counter) Test2(value int) int {
	return 0
}

// @contract:public
// @contract:payable min_deposit=0.0005
// @contract:mutating
func (c *Counter) Test3(msg string) string {
	return ""
}

// @contract:private
func (c *Counter) Test4() {
	// internal logic
}

// @contract:public
// @contract:view
func (c *Counter) Test5(key string) int {
	return 0
}

// @contract:public
// @contract:mutating
func (c *Counter) Test6(a, b int) int {
	return 0
}

// @contract:public
// @contract:payable min_deposit=1
// @contract:mutating
func (c *Counter) Test7() string {
	return ""
}

// @contract:private
func (c *Counter) Test8(flag bool) {
	// internal logic
}

// @contract:public
// @contract:view
func (c *Counter) Test9() bool {
	return false
}

// @contract:public
// @contract:mutating
func (c *Counter) Test10(name string, amount int) string {
	return ""
}
