package main

import (
	"fmt"
	"strconv"
)

//go:export helloworld
func helloworld() {
	accountBalance := GetAccountBalance()

	LogString("Default format hex:" + accountBalance.String())

	// Convert Hi and Lo parts to strings using strconv.FormatUint
	hiStr := strconv.FormatUint(accountBalance.Hi, 10)
	loStr := strconv.FormatUint(accountBalance.Lo, 10)

	// Concatenate the strings
	decimalStr := hiStr + loStr

	LogString("Default format hex: " + accountBalance.String())
	LogString("Decimal value: " + fmt.Sprintf("%s", decimalStr))
	LogString("hiStr value: " + fmt.Sprintf("%s", hiStr))
	LogString("loStr value: " + fmt.Sprintf("%s", loStr))

}

func main() {
}
