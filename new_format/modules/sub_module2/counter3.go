package main

// @contract:public
// @contract:view
func (c *Counter) Test13() string {
	return ""
}

// @contract:public
// @contract:mutating
func (c *Counter) Test23(value int) int {
	return 0
}

// @contract:public
// @contract:payable min_deposit=0.0005
// @contract:mutating
func (c *Counter) Test33(msg string) string {
	return ""
}

// @contract:private
func (c *Counter) Test43() {
	// internal logic
}

// @contract:public
// @contract:view
func (c *Counter) Test53(key string) int {
	return 0
}

// @contract:public
// @contract:mutating
func (c *Counter) Test63(a, b int) int {
	return 0
}

// @contract:public
// @contract:payable min_deposit=1
// @contract:mutating
func (c *Counter) Test73() string {
	return ""
}

// @contract:private
func (c *Counter) Test83(flag bool) {
	// internal logic
}

// @contract:public
// @contract:view
func (c *Counter) Test93() bool {
	return false
}

// @contract:public
// @contract:mutating
func (c *Counter) Test103(name string, amount int) string {
	return ""
}
