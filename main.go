package main

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	contractInput, _, inputErr := ContractInput(true)
	if inputErr != nil {
		LogString("There are some error :" + inputErr.Error())

	}
	data, readErr := StorageRead(StorageKey)
	if readErr != nil {
		LogString("There are some error :" + readErr.Error())
	}
	state, deserializeErr := Deserialize(data)
	if deserializeErr != nil {
		LogString("There are some error :" + deserializeErr.Error())
	}
	accountId, errAccountId := GetPredecessorAccountID()
	if errAccountId != nil {
		LogString("There are some error :" + errAccountId.Error())
	}
	state.Data[accountId] = string(contractInput)
	StorageWrite(StorageKey, state.Serialize())
	ContractValueReturn([]byte(state.Data[accountId]))
}

//go:export GetStatus
func GetStatus() {
	data, readErr := StorageRead(StorageKey)
	if readErr != nil {
		LogString("There are some error :" + readErr.Error())
	}
	state, deserializeErr := Deserialize(data)
	if deserializeErr != nil {
		LogString("There are some error :" + deserializeErr.Error())
	}
	accountId, errAccountId := GetPredecessorAccountID()
	if errAccountId != nil {
		LogString("There are some error :" + errAccountId.Error())
	}
	ContractValueReturn([]byte(state.Data[accountId]))
}

//go:export InitContract
func InitContract() {
	LogString("Init Smart Contract")
	msg := NewStatusMessage()
	serialized := msg.Serialize()
	StorageWrite(StorageKey, serialized)

	balacne := GetAccountBalance()
	LogString("balacne.String() " + balacne.String())

	res, err := U128FromString("340282366920938463463374607431768211455")
	if err != nil {
		LogString("U128FromString Error " + err.Error())
	}

	LogString("U128FromString : " + res.String())

}
