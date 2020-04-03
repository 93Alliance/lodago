package lodago

// If 三元表达式 If(a > b, c, d)
func If(condition bool, trueValue, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}
