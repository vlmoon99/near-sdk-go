package main

import (
	"fmt"
)

//go:export helloworld
func helloworld() {
	value, dataType, err := GetSmartContractInput()
	if err != nil {
		SmartContractLog("Error in input: " + err.Error())
		return
	}
	SmartContractLog("dataType: " + dataType)
	SmartContractLog("value: " + string(value))

	accountID, err := GetCurrentAccountID()
	if err != nil {
		SmartContractLog("Error in GetCurrentAccountID: " + err.Error())
		return
	}
	SmartContractLog("CurrentAccountID: " + accountID)

	signerID, err := GetSignerAccountID()
	if err != nil {
		SmartContractLog("Error in GetSignerAccountID: " + err.Error())
		return
	}
	SmartContractLog("SignerAccountID: " + signerID)

	predecessorID, err := GetPredecessorAccountID()
	if err != nil {
		SmartContractLog("Error in GetPredecessorAccountID: " + err.Error())
		return
	}
	SmartContractLog("PredecessorAccountID: " + predecessorID)

	signerPK, err := GetSignerAccountPK()
	if err != nil {
		SmartContractLog("Error in GetSignerAccountPK: " + err.Error())
		return
	}

	SmartContractLog(" GetSignerAccountPK - Len: " + fmt.Sprintf("%d", len(signerPK)))

}

func main() {
}
