package main

import (
	"encoding/base64"
	"fmt"
	"unsafe"

	"github.com/buger/jsonparser"
	"github.com/mr-tron/base58"
)

func encodeTestBase58(inputBytes []byte) string {
	return base58.Encode(inputBytes)
}

func encodeTestBase64(inputBytes []byte) string {
	return base64.StdEncoding.EncodeToString(inputBytes)
}

// Custom function to manually build a JSON string from a map
func buildJSONFromMap(parsedMap map[string]string) string {
	jsonStr := "{"
	first := true

	for key, value := range parsedMap {
		if !first {
			jsonStr += ","
		}
		first = false
		jsonStr += fmt.Sprintf("\"%s\":\"%s\"", key, value)
	}

	jsonStr += "}"
	return jsonStr
}

//go:export helloworld
func helloworld() {
	data, err := SmartContractInput()
	if err != nil {
		LogString([]byte("Error"))
		return
	}
	LogString(data)

	helloValue, err := jsonparser.GetString(data, "hello")
	if err != nil {
		LogString([]byte("Error extracting hello value"))
		return
	}

	data, errBase58 := jsonparser.Set(data, []byte(`"`+encodeTestBase58([]byte(helloValue))+`"`), "world_base58")
	if errBase58 != nil {
		LogString([]byte("Error setting hello_base58"))
		return
	}

	data, errBase64 := jsonparser.Set(data, []byte(`"`+encodeTestBase64([]byte(helloValue))+`"`), "world_base64")
	if errBase64 != nil {
		LogString([]byte("Error setting hello_base64"))
		return
	}

	LogString(data)

	dataLen := uint64(len(data))
	dataPtr := uint64(uintptr(unsafe.Pointer(&data[0])))

	ValueReturn(dataLen, dataPtr)
}

func main() {
}
