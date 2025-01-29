package main

import (
	"encoding/base64"

	"github.com/mr-tron/base58"
)

func encodeTestBase58(inputBytes []byte) string {
	return base58.Encode(inputBytes)
}

func encodeTestBase64(inputBytes []byte) string {
	return base64.StdEncoding.EncodeToString(inputBytes)
}

//go:export helloworld
func helloworld() {
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

	// signerPK, err := GetSignerAccountPK()
	// if err != nil {
	// 	SmartContractLog("Error in GetSignerAccountPK: " + err.Error())
	// 	return
	// }

	// SmartContractLog(" GetSignerAccountPK - Len: " + fmt.Sprintf("%d", len(signerPK)))

}

func main() {
}
