package main

import (
	"fmt"

	"github.com/vlmoon99/near-sdk-go/collections"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
	"github.com/vlmoon99/near-sdk-go/types"
)

type StatusMessage struct {
	Data *collections.LookupMap
}

func GetState() StatusMessage {
	return StatusMessage{
		Data: collections.NewLookupMap([]byte("b")),
	}
}

//go:export SetStatus
func SetStatus() {
	options := types.ContractInputOptions{IsRawBytes: true}
	contractInput, _, inputErr := env.ContractInput(options)
	if inputErr != nil {
		env.PanicStr("Input error: " + inputErr.Error())
	}

	accountId, errAccountId := env.GetPredecessorAccountID()
	if errAccountId != nil {
		env.PanicStr("Account ID error: " + errAccountId.Error())
	}

	state := GetState()
	errInsert := state.Data.Insert([]byte(accountId), string(contractInput))
	if errInsert != nil {
		env.PanicStr("Error inserting into LookupMap : " + errInsert.Error())
	}

	env.ContractValueReturn([]byte(contractInput))
}

//go:export GetStatus
func GetStatus() {
	accountId, errAccountId := env.GetPredecessorAccountID()
	if errAccountId != nil {
		env.PanicStr("Account ID error: " + errAccountId.Error())
	}

	state := GetState()

	val, err := state.Data.Get([]byte(accountId))
	if err != nil {
		env.PanicStr("Error getting from LookupMap : " + err.Error())
	}

	if val == nil {
		env.PanicStr("No status found for this account")
	}

	status, ok := val.(string)
	if !ok {
		env.PanicStr("Error: Value in LookupMap is not a string")
	}

	env.ContractValueReturn([]byte(status))
}

//go:export InitContract
func InitContract() {
	env.LogString("Init Smart Contract")
	realSys := system.RealSystem{}

	len := realSys.RegisterLen(env.AtomicOpRegister)

	env.LogString("len  : " + fmt.Sprintf("%d", len))

}
