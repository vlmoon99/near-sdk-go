package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
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

	res, err := U128FromString("10000000000000000000000000")
	if err != nil {
		LogString("U128FromString res Error " + err.Error())
	}

	LogString("U128FromString res : " + res.String())

	res1, err1 := U128FromString("5")
	if err1 != nil {
		LogString("U128FromString res1 Error " + err1.Error())
	}

	LogString("U128FromString res1 : " + res1.String())

	resAddition, additionErr := res.Add(res1)
	if additionErr != nil {
		LogString("additionErr " + additionErr.Error())
	}

	LogString("U128FromString resAddition : " + resAddition.String())

	resSub, subErr := resAddition.Sub(res1)
	if subErr != nil {
		LogString("subErr " + subErr.Error())
	}

	LogString("U128FromString resSub : " + resSub.String())

	resMul, mulErr := resSub.Mul(res1)
	if mulErr != nil {
		LogString("mulErr " + mulErr.Error())
	}

	LogString("U128FromString resMul : " + resMul.String())

	resDiv, divErr := resMul.Div(res1)
	if divErr != nil {
		LogString("divErr " + divErr.Error())
	}

	LogString("U128FromString resDiv : " + resDiv.String())

	a := Uint128{Lo: 1, Hi: 0}
	b := Uint128{Lo: 0xFFFFFFFFFFFFFFFF, Hi: 0xFFFFFFFFFFFFFFFF}

	LogString("Bit (a, 0): " + fmt.Sprintf("%d", a.Bit(0)))   // Expected output: 1
	LogString("Bit (a, 1): " + fmt.Sprintf("%d", a.Bit(1)))   // Expected output: 0
	LogString("Bit (b, 63): " + fmt.Sprintf("%d", b.Bit(63))) // Expected output: 1
	LogString("Bit (b, 64): " + fmt.Sprintf("%d", b.Bit(64))) // Expected output: 1

	c := a.Lsh(1)
	LogString("Lsh (a, 1).Lo: " + fmt.Sprintf("%d", c.Lo)) // Expected output: 2
	LogString("Lsh (a, 1).Hi: " + fmt.Sprintf("%d", c.Hi)) // Expected output: 0

	d := b.Lsh(1)
	LogString("Lsh (b, 1).Lo: " + fmt.Sprintf("%d", d.Lo)) // Expected output: 0xFFFFFFFFFFFFFFFE
	LogString("Lsh (b, 1).Hi: " + fmt.Sprintf("%d", d.Hi)) // Expected output: 0xFFFFFFFFFFFFFFFF

	LogString("Cmp (a, b): " + fmt.Sprintf("%d", a.Cmp(b))) // Expected output: -1
	LogString("Cmp (b, a): " + fmt.Sprintf("%d", b.Cmp(a))) // Expected output: 1
	LogString("Cmp (a, a): " + fmt.Sprintf("%d", a.Cmp(a))) // Expected output: 0

	// Test Mod method
	modRes, modErr := res.Mod(res1)
	if modErr != nil {
		LogString("modErr " + modErr.Error())
	}
	LogString("U128FromString modRes : " + modRes.String()) // Expected output: 0

	// Test Bitwise AND method
	andRes := res.And(res1)
	LogString("U128FromString andRes : " + andRes.String()) // Expected output: 0

	// Test Bitwise OR method
	orRes := res.Or(res1)
	LogString("U128FromString orRes : " + orRes.String()) // Expected output: 10000000000000000000000005

	// Test Bitwise XOR method
	xorRes := res.Xor(res1)
	LogString("U128FromString xorRes : " + xorRes.String()) // Expected output: 10000000000000000000000005

	// Test ShiftRight method
	srRes := res.ShiftRight(1)
	LogString("U128FromString srRes : " + srRes.String()) // Expected output: 5000000000000000000000000
}
