package main

// @contract:public
// @contract:view
func (c *Counter) Test12() string {
	return ""
}

// @contract:public
// @contract:mutating
func (c *Counter) Test22(value int) int {
	return 0
}

// @contract:public
// @contract:payable min_deposit=0.0005
// @contract:mutating
func (c *Counter) Test32(msg string) string {
	return ""
}

// @contract:private
func (c *Counter) Test42() {
	// internal logic
}

// @contract:public
// @contract:view
func (c *Counter) Test52(key string) int {
	return 0
}

// @contract:public
// @contract:mutating
func (c *Counter) Test62(a, b int) int {
	return 0
}

// @contract:public
// @contract:payable min_deposit=1
// @contract:mutating
func (c *Counter) Test72() string {
	return ""
}

// @contract:private
func (c *Counter) Test82(flag bool) {
	// internal logic
}

// @contract:public
// @contract:view
func (c *Counter) Test92() bool {
	return false
}

// @contract:public
// @contract:mutating
func (c *Counter) Test102(name string, amount int) string {
	return ""
}
