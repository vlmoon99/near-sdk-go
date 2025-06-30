// Package contract provides functions for createing the basic structure ofthe smart contract.
package contract

import (
	"errors"

	"github.com/vlmoon99/near-sdk-go/borsh"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/json"
	"github.com/vlmoon99/near-sdk-go/promise"
	"github.com/vlmoon99/near-sdk-go/types"
)

func HandleClientJSONInput(fn func(*ContractInput) error) {
	input, err := GetJSONInput()
	if err != nil {
		env.PanicStr("failed to get input: " + err.Error())
	}
	if err := fn(input); err != nil {
		env.PanicStr(err.Error())
	}
}

func HandleClientRawBytesInput(fn func(*ContractInput) error) {
	input, err := GetRawBytesInput()
	if err != nil {
		env.PanicStr("failed to get input: " + err.Error())
	}
	if err := fn(input); err != nil {
		env.PanicStr(err.Error())
	}
}

func HandlePromiseResult(fn func(*promise.PromiseResult) error) {
	if err := promise.CallbackGuard(); err != nil {
		env.PanicStr("callback rejected: " + err.Error())
	}

	result, err := promise.GetPromiseResultSafe(0)
	if err != nil {
		env.PanicStr("failed to get promise result: " + err.Error())
	}

	if err := fn(&result); err != nil {
		env.PanicStr(err.Error())
	}
}

func ReturnValue(value interface{}) error {
	var data []byte
	var err error

	switch v := value.(type) {
	case []byte:
		data = v
	default:
		data, err = borsh.Serialize(v)
		if err != nil {
			return err
		}
	}

	env.ContractValueReturn(data)
	return nil
}

type ContractInput struct {
	Data []byte
	JSON *json.Parser
}

func GetJSONInput() (*ContractInput, error) {
	options := types.ContractInputOptions{IsRawBytes: false}
	data, _, err := env.ContractInput(options)
	if err != nil {
		return nil, err
	}

	parser := json.NewParser(data)
	return &ContractInput{
		Data: data,
		JSON: parser,
	}, nil
}

func GetRawBytesInput() (*ContractInput, error) {
	options := types.ContractInputOptions{IsRawBytes: true}
	data, _, err := env.ContractInput(options)
	if err != nil {
		return nil, err
	}

	return &ContractInput{
		Data: data,
		JSON: nil,
	}, nil
}

type ContractContext struct {
	AccountID       string
	SignerID        string
	PredecessorID   string
	AttachedDeposit types.Uint128
	PrepaidGas      uint64
}

func GetContext() *ContractContext {
	accountID, _ := env.GetCurrentAccountId()
	signerID, _ := env.GetSignerAccountID()
	predecessorID, _ := env.GetPredecessorAccountID()
	attachedDeposit, _ := env.GetAttachedDeposit()
	prepaidGas := env.GetPrepaidGas()

	return &ContractContext{
		AccountID:       accountID,
		SignerID:        signerID,
		PredecessorID:   predecessorID,
		AttachedDeposit: attachedDeposit,
		PrepaidGas:      prepaidGas.Inner,
	}
}

func RequireDeposit(minDeposit types.Uint128) error {
	context := GetContext()
	if context.AttachedDeposit.Lo < minDeposit.Lo ||
		(context.AttachedDeposit.Lo == minDeposit.Lo && context.AttachedDeposit.Hi < minDeposit.Hi) {
		return errors.New("insufficient deposit")
	}
	return nil
}

type PromiseResult struct {
	Success    bool
	Data       []byte
	StatusCode int
}
