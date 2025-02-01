package main

import (
	"fmt"
)

//go:export helloworld
func helloworld() {
	accountBalance := GetAccountBalance()

	LogString("accountBalance.Hi: " + fmt.Sprintf("%d", accountBalance.Hi) + " accountBalance.Lo: " + fmt.Sprintf("%d", accountBalance.Lo))
	LogString("accountBalance HexBE: " + accountBalance.HexBE())
	LogString("accountBalance ToFloat64: " + fmt.Sprintf("%d", accountBalance.ToFloat64()))
	LogString("accountBalance ToYoctoNear: " + fmt.Sprintf("%.0f", accountBalance.ToYoctoNear()))

}

func main() {
}
