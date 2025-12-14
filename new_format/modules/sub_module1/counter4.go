package main

// @contract:public
// @contract:view
func (c *Counter) Test14() string {
	return ""
}

// @contract:public
// @contract:mutating
func (c *Counter) Test24(value int) int {
	return 0
}

// @contract:public
// @contract:payable min_deposit=0.0005
// @contract:mutating
func (c *Counter) Test34(msg string) string {
	return ""
}

// @contract:private
func (c *Counter) Test44() {
	// internal logic
}

// @contract:public
// @contract:view
func (c *Counter) Test54(key string) int {
	return 0
}

// @contract:public
// @contract:mutating
func (c *Counter) Test64(a, b int) int {
	return 0
}

// @contract:public
// @contract:payable min_deposit=1
// @contract:mutating
func (c *Counter) Test74() string {
	return ""
}

// @contract:private
func (c *Counter) Test84(flag bool) {
	// internal logic
}

// @contract:public
// @contract:view
func (c *Counter) Test94() bool {
	return false
}

// @contract:public
// @contract:mutating
func (c *Counter) Test104(name string, amount int) string {
	return ""
}
