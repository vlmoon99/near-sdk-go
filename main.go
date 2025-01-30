package main

//go:export helloworld
func helloworld() {
	accountBalance := GetAccountBalance()

	LogString("Default format hex: " + accountBalance.String())

	decimalStr := hexToDecimal(accountBalance.String())

	LogString("Decimal value: " + decimalStr)
}

func main() {
}
