package main

import (
	"fmt"
	"unsafe"

	"github.com/vlmoon99/near-sdk-go/env"
)

//go:export TestReadRegisterSafe
func TestReadRegisterSafe() {
	data := []byte("test data")
	registerId := uint64(1)

	env.WriteRegisterSafe(registerId, data)

	readData, err := env.ReadRegisterSafe(registerId)
	if err != nil {
		env.PanicStr("Failed to read from register:  " + err.Error())
	}

	if string(readData) != "test data" {
		env.ContractValueReturn([]byte("0"))
	} else {
		env.ContractValueReturn([]byte("1"))
	}
}

//go:export TestStorageAPIWithSimpleType
func TestStorageAPIWithSimpleType() {
	key := []byte("TEEEETSTs")
	value := []byte("Test111")

	_, err := env.StorageWrite(key, value)
	if err != nil {
		env.PanicStr("Failed to write simple type to storage: " + err.Error())
	}

	readValue, err := env.StorageRead(key)
	if err != nil {
		env.PanicStr("Failed to read simple type from storage: " + err.Error())
	}
	env.LogString("string(readValue)  : " + string(readValue))
	env.LogString("string(value)  : " + string(value))

	if string(readValue) != string(value) {
		env.ContractValueReturn([]byte("0"))
	} else {
		env.ContractValueReturn([]byte("1"))
	}
}
func GetKeyAndValuePointers(key, value []byte) (keyLen uint64, keyPtr uint64, valueLen uint64, valuePtr uint64, registerId uint64) {
	if len(key) == 0 {
		return 0, 0, 0, 0, env.EvictedRegister
	}

	if len(value) == 0 {
		return 0, 0, 0, 0, env.EvictedRegister
	}

	keyLen = uint64(len(key))
	keyPtr = uint64(uintptr(unsafe.Pointer(&key)))

	valueLen = uint64(len(value))
	valuePtr = uint64(uintptr(unsafe.Pointer(&value)))

	return keyLen, keyPtr, valueLen, valuePtr, env.EvictedRegister
}

//go:export InitContract
func main() {
	env.LogString("Init Smart Contract")
	key := []byte("TEEEETyhghSTs")
	value := []byte("Test111")

	result, _ := env.StorageWrite(key, value)

	// res, err := StorageWrite(Key, value)
	// if err != nil {
	// 	env.LogString("err.Error   : " + err.Error())

	// }
	env.LogString("result   : " + fmt.Sprintf("%d", result))

}

// package main

// import (
// 	"github.com/vlmoon99/near-sdk-go/collections"
// 	"github.com/vlmoon99/near-sdk-go/env"
// 	"github.com/vlmoon99/near-sdk-go/types"
// )

// type StatusMessage struct {
// 	Data *collections.LookupMap
// }

// func GetState() StatusMessage {

// 	return StatusMessage{
// 		Data: collections.NewLookupMap([]byte("b")),
// 	}
// }

// //go:export SetStatus
// func SetStatus() {

// 	options := types.ContractInputOptions{IsRawBytes: true}
// 	contractInput, _, inputErr := env.ContractInput(options)
// 	if inputErr != nil {
// 		env.PanicStr("Input error: " + inputErr.Error())
// 	}

// 	accountId, errAccountId := env.GetPredecessorAccountID()
// 	if errAccountId != nil {
// 		env.PanicStr("Account ID error: " + errAccountId.Error())
// 	}

// 	state := GetState()
// 	errInsert := state.Data.Insert([]byte(accountId), string(contractInput))
// 	if errInsert != nil {
// 		env.PanicStr("Error inserting into LookupMap : " + errInsert.Error())
// 	}

// 	env.ContractValueReturn([]byte(contractInput))
// }

// //go:export GetStatus
// func GetStatus() {

// 	accountId, errAccountId := env.GetPredecessorAccountID()
// 	if errAccountId != nil {
// 		env.PanicStr("Account ID error: " + errAccountId.Error())
// 	}

// 	state := GetState()

// 	val, err := state.Data.Get([]byte(accountId))
// 	if err != nil {
// 		env.PanicStr("Error getting from LookupMap : " + err.Error())
// 	}

// 	if val == nil {
// 		env.PanicStr("No status found for this account")
// 	}

// 	status, ok := val.(string)
// 	if !ok {
// 		env.PanicStr("Error: Value in LookupMap is not a string")
// 	}

// 	env.ContractValueReturn([]byte(status))
// }

// //go:export InitContract
// func InitContract() {
// 	env.LogString("Init Smart Contract")
// }
