package main

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/vlmoon99/near-sdk-go/sdk"
	"github.com/vlmoon99/near-sdk-go/types"
)

type StatusMessage struct {
	Data map[string]string
}

var StorageKey = []byte("State")

func NewStatusMessage() *StatusMessage {
	return &StatusMessage{
		Data: make(map[string]string),
	}
}

func (sm *StatusMessage) Serialize() []byte {
	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, uint16(len(sm.Data)))

	for key, value := range sm.Data {
		keyBytes := []byte(key)
		valueBytes := []byte(value)

		binary.Write(&buffer, binary.BigEndian, uint16(len(keyBytes)))
		buffer.Write(keyBytes)

		binary.Write(&buffer, binary.BigEndian, uint16(len(valueBytes)))
		buffer.Write(valueBytes)
	}

	return buffer.Bytes()
}

func Deserialize(data []byte) (*StatusMessage, error) {
	if len(data) < 2 {
		return nil, errors.New("invalid data: too short")
	}

	buffer := bytes.NewReader(data)

	var numPairs uint16
	binary.Read(buffer, binary.BigEndian, &numPairs)

	dataMap := make(map[string]string)

	for i := 0; i < int(numPairs); i++ {
		var keyLen, valueLen uint16

		if err := binary.Read(buffer, binary.BigEndian, &keyLen); err != nil {
			return nil, errors.New("failed to read key length")
		}

		keyBytes := make([]byte, keyLen)
		if _, err := buffer.Read(keyBytes); err != nil {
			return nil, errors.New("failed to read key")
		}
		key := string(keyBytes)

		if err := binary.Read(buffer, binary.BigEndian, &valueLen); err != nil {
			return nil, errors.New("failed to read value length")
		}

		valueBytes := make([]byte, valueLen)
		if _, err := buffer.Read(valueBytes); err != nil {
			return nil, errors.New("failed to read value")
		}
		value := string(valueBytes)

		dataMap[key] = value
	}

	return &StatusMessage{Data: dataMap}, nil
}

//go:export SetStatus
func SetStatus() {
	options := types.ContractInputOptions{IsRawBytes: true}
	contractInput, _, inputErr := sdk.ContractInput(options)
	if inputErr != nil {
		sdk.LogString("There are some error :" + inputErr.Error())

	}
	data, readErr := sdk.StorageRead(StorageKey)
	if readErr != nil {
		sdk.LogString("There are some error :" + readErr.Error())
	}
	state, deserializeErr := Deserialize(data)
	if deserializeErr != nil {
		sdk.LogString("There are some error :" + deserializeErr.Error())
	}
	accountId, errAccountId := sdk.GetPredecessorAccountID()
	if errAccountId != nil {
		sdk.LogString("There are some error :" + errAccountId.Error())
	}
	state.Data[accountId] = string(contractInput)
	sdk.StorageWrite(StorageKey, state.Serialize())
	sdk.ContractValueReturn([]byte(state.Data[accountId]))
}

//go:export GetStatus
func GetStatus() {
	data, readErr := sdk.StorageRead(StorageKey)
	if readErr != nil {
		sdk.LogString("There are some error :" + readErr.Error())
	}
	state, deserializeErr := Deserialize(data)
	if deserializeErr != nil {
		sdk.LogString("There are some error :" + deserializeErr.Error())
	}
	accountId, errAccountId := sdk.GetPredecessorAccountID()
	if errAccountId != nil {
		sdk.LogString("There are some error :" + errAccountId.Error())
	}
	sdk.ContractValueReturn([]byte(state.Data[accountId]))
}

//go:export InitContract
func InitContract() {
	input, dataType, err := sdk.ContractInput(types.ContractInputOptions{IsRawBytes: true})
	if err != nil {
		sdk.LogString("Input error :" + err.Error())
	}
	sdk.LogString("Init Smart Contract dataType :" + dataType)
	sdk.LogString("Init Smart Contract input len :" + string(input))

	sdk.LogString("Init Smart Contract")
	msg := NewStatusMessage()
	serialized := msg.Serialize()
	sdk.StorageWrite(StorageKey, serialized)
}
