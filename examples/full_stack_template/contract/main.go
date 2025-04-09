package main

import (
	"encoding/hex"
	"fmt"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/json"
	"github.com/vlmoon99/near-sdk-go/types"
)

//go:export InitContract
func InitContract() {
	env.LogString("Init Smart Contract")
	env.ContractValueReturn([]byte("Init Smart Contract"))
}

//go:export WriteData
func WriteData() {
	options := types.ContractInputOptions{IsRawBytes: true}
	input, detectedType, err := env.ContractInput(options)
	if err != nil {
		env.PanicStr("Failed to get contract input: " + err.Error())
	}

	env.LogString("Contract input (JSON): " + string(input))
	env.LogString("Detected input type: " + detectedType)

	if detectedType != "object" {
		env.ContractValueReturn([]byte("Error : Incorrect type"))
	}
	parser := json.NewParser(input)

	keyResult, err := parser.GetString("key")
	if err != nil {
		env.ContractValueReturn([]byte("Error : Incorrect key"))
	}

	dataResult, err := parser.GetRawBytes("data")
	if err != nil {
		env.ContractValueReturn([]byte("Error : Incorrect data"))
	}

	resultStorageWrite, err := env.StorageWrite([]byte(keyResult), dataResult)
	if err != nil {
		env.ContractValueReturn([]byte("Error : " + err.Error()))
	}

	if resultStorageWrite {
		env.LogString("env.StorageWrite returned true")
	} else {
		env.LogString("env.StorageWrite returned false")
	}

	env.ContractValueReturn([]byte("WriteData was successful"))
}

//go:export ReadData
func ReadData() {
	options := types.ContractInputOptions{IsRawBytes: true}
	input, detectedType, err := env.ContractInput(options)
	if err != nil {
		env.PanicStr("Failed to get contract input: " + err.Error())
	}

	env.LogString("Contract input (JSON): " + string(input))
	env.LogString("Detected input type: " + detectedType)

	if detectedType != "object" {
		env.ContractValueReturn([]byte("Error : Incorrect type" + "detected type is " + detectedType))
	}
	parser := json.NewParser(input)

	keyResult, err := parser.GetString("key")
	if err != nil {
		env.ContractValueReturn([]byte("Error : Incorrect key"))
	}

	data, err := env.StorageRead([]byte(keyResult))
	if err != nil {
		env.ContractValueReturn([]byte("Error : Incorrect read from the storage by that key"))
	}
	env.LogString("ReadData was successful")

	env.ContractValueReturn(data)
}

//go:export AcceptPayment
func AcceptPayment() {
	attachedDeposit, err := env.GetAttachedDepoist()
	if err != nil {
		env.PanicStr("Failed to get attached deposit: " + err.Error())
	}
	env.LogString("Attachet Deposit :" + attachedDeposit.String())
	promiseIdx := env.PromiseBatchCreate([]byte("neargoclitest.testnet"))
	env.PromiseBatchActionTransfer(promiseIdx, attachedDeposit)
	//neargocli.testnet
	env.ContractValueReturn([]byte("AcceptPayment"))

}

//go:export ReadIncommingTxData
func ReadIncommingTxData() {
	env.LogString(`EVENT_JSON:{
  "standard": "nep999",
  "version": "1.0.0",
  "event": "ReadIncommingTxData",
  "data": [
    {"info": "ReadIncommingTxData", "test": ["test11"]}
  ]
}`)
	options := types.ContractInputOptions{IsRawBytes: true}
	input, detectedType, err := env.ContractInput(options)
	if err != nil {
		env.PanicStr("Failed to get contract input: " + err.Error())
	}
	env.LogString("Contract input (raw bytes): " + string(input))
	env.LogString("Detected input type: " + detectedType)

	attachedDeposit, err := env.GetAttachedDepoist()
	if err != nil {
		env.PanicStr("Failed to get attached deposit: " + err.Error())
	}
	env.LogString(fmt.Sprintf("Attached deposit: %s", attachedDeposit.String()))

	accountId, err := env.GetCurrentAccountId()
	if err != nil || accountId == "" {
		env.PanicStr("Failed to get current account ID: " + err.Error())
	}
	env.LogString("Current account ID: " + accountId)

	signerId, err := env.GetSignerAccountID()
	if err != nil || signerId == "" {
		env.PanicStr("Failed to get signer account ID: " + err.Error())
	}
	env.LogString("Signer account ID: " + signerId)

	signerPK, err := env.GetSignerAccountPK()
	if err != nil || signerPK == nil {
		env.PanicStr("Failed to get signer account PK: " + err.Error())
	}
	env.LogString("Signer account PK: " + hex.EncodeToString(signerPK))

	predecessorId, err := env.GetPredecessorAccountID()
	if err != nil || predecessorId == "" {
		env.PanicStr("Failed to get predecessor account ID: " + err.Error())
	}
	env.LogString("Predecessor account ID: " + predecessorId)

	blockHeight := env.GetCurrentBlockHeight()
	env.LogString("Current block height: " + fmt.Sprintf("%d", blockHeight))

	blockTimeMs := env.GetBlockTimeMs()
	env.LogString("Block time in ms: " + fmt.Sprintf("%d", blockTimeMs))

	epochHeight := env.GetEpochHeight()
	env.LogString("Epoch height: " + fmt.Sprintf("%d", epochHeight))

	storageUsage := env.GetStorageUsage()
	env.LogString("Storage usage: " + fmt.Sprintf("%d", storageUsage))

	accountBalance, err := env.GetAccountBalance()
	if err != nil {
		env.PanicStr("Failed to get account balance: " + err.Error())
	}
	env.LogString(fmt.Sprintf("Account balance: %s", accountBalance.String()))

	lockedBalance, err := env.GetAccountLockedBalance()
	if err != nil {
		env.PanicStr("Failed to get account locked balance: " + err.Error())
	}
	env.LogString(fmt.Sprintf("Account locked balance: %s", lockedBalance.String()))

	prepaidGas := env.GetPrepaidGas()
	env.LogString(fmt.Sprintf("Prepaid gas: %ds", prepaidGas.Inner))

	usedGas := env.GetUsedGas()
	env.LogString(fmt.Sprintf("Used gas: %d", usedGas.Inner))

	env.ContractValueReturn([]byte("ReadIncommingTxData"))
}

//go:export ReadBlockchainData
func ReadBlockchainData() {
	//neargocli.testnet
	env.ContractValueReturn([]byte("ReadBlockchainData"))
}
