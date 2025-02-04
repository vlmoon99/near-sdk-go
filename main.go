package main

// 0.Create structs
type Person struct {
	Name  string
	Age   int
	Email string
}

//go:export testpromisethen
func testpromisethen() {
	value, dataType, err := ContractInput(true)

	if err != nil {
		LogString("Error while getting contract input")
	}

	LogString("Test Promise Then")
	LogString("Test Promise Then value " + string(value))
	LogString("Test Promise Then dataType " + dataType)

}

//go:export testpromise
func testpromise() {
	value, dataType, err := ContractInput(true)

	if err != nil {
		LogString("Error while getting contract input")
	}

	LogString("Test Promise")
	LogString("Test Promise value " + string(value))
	LogString("Test Promise dataType " + dataType)

}

//go:export testpromise1
func testpromise1() {
	value, dataType, err := ContractInput(true)

	if err != nil {
		LogString("Error while getting contract input")
	}

	LogString("Test Promise1")
	LogString("Test Promise1 value " + string(value))
	LogString("Test Promise1 dataType " + dataType)

}

//go:export helloworld
func helloworld() {
	// promiseBatchCreate := PromiseBatchCreate([]byte("testiktinygo.testnet"))
	// LogString("promiseBatchCreate : " + fmt.Sprintf("%d", promiseBatchCreate))

	// promiseValue := PromiseCreate([]byte("testiktinygo.testnet"), []byte("testpromise"), []byte("test"), Uint128{0, 0}, 5*ONE_TERA_GAS)
	// LogString("promiseValue : " + fmt.Sprintf("%d", promiseValue))

	// promiseBatchThen := PromiseBatchThen(promiseValue, []byte("testiktinygo.testnet"))
	// LogString("promiseBatchThen : " + fmt.Sprintf("%d", promiseBatchThen))

	// promiseValue1 := PromiseCreate([]byte("testiktinygo.testnet"), []byte("testpromise1"), []byte("test"), Uint128{0, 0}, 5*ONE_TERA_GAS)

	// promiseThenValue := PromiseThen(promiseValue, []byte("testiktinygo.testnet"), []byte("testpromisethen"), []byte("test"), Uint128{0, 0}, 5*ONE_TERA_GAS)
	// LogString("promiseValue : " + fmt.Sprintf("%d", promiseValue))
	// LogString("promiseValue1 : " + fmt.Sprintf("%d", promiseValue1))
	// LogString("promiseThenValue : " + fmt.Sprintf("%d", promiseThenValue))

	// PromiseAnd([]uint64{promiseValue1, promiseValue, promiseThenValue})

	// // Test AltBn128G1MultiExp
	// _, err := AltBn128G1MultiExp(data)
	// if err != nil {
	// 	LogString("Error while executing AltBn128G1MultiExp: " + err.Error())
	// } else {
	// 	LogString("AltBn128G1MultiExp executed successfully: ")
	// }

	// // Test AltBn128G1Sum
	// _, err = AltBn128G1Sum(data)
	// if err != nil {
	// 	LogString("Error while executing AltBn128G1Sum: " + err.Error())
	// } else {
	// 	LogString("AltBn128G1Sum executed successfully: ")
	// }

	// // Test AltBn128PairingCheck
	// if AltBn128PairingCheck(data) {
	// 	LogString("AltBn128PairingCheck verified successfully.")
	// } else {
	// 	LogString("AltBn128PairingCheck verification failed.")
	// }

	// totalStake := ValidatorTotalStakeAmount()
	// LogString("Total validator stake: " + fmt.Sprintf("%v", totalStake))

	// data := []byte("sample_data")

	// randomSeed, err := GetRandomSeed()
	// if err != nil {
	// 	LogString("Error while generating random seed: " + err.Error())
	// } else {
	// 	LogString("Random seed generated successfully : " + fmt.Sprintf("%v", randomSeed))
	// }

	// sha256Hash, err := Sha256Hash(data)
	// if err != nil {
	// 	LogString("Error while hashing with SHA-256: " + err.Error())
	// } else {
	// 	LogString("SHA-256 Hash: " + fmt.Sprintf("%v", sha256Hash))
	// }

	// keccak256Hash, err := Keccak256Hash(data)
	// if err != nil {
	// 	LogString("Error while hashing with Keccak-256: " + err.Error())
	// } else {
	// 	LogString("Keccak-256 Hash: " + fmt.Sprintf("%v", keccak256Hash))
	// }

	// keccak512Hash, err := Keccak512Hash(data)
	// if err != nil {
	// 	LogString("Error while hashing with Keccak-512: " + err.Error())
	// } else {
	// 	LogString("Keccak-512 Hash: " + fmt.Sprintf("%v", keccak512Hash))
	// }

	// ripemd160Hash, err := Ripemd160Hash(data)
	// if err != nil {
	// 	LogString("Error while hashing with RIPEMD-160: " + err.Error())
	// } else {
	// 	LogString("RIPEMD-160 Hash: " + fmt.Sprintf("%v", ripemd160Hash))
	// }

	// state, err := StateRead()
	// if err != nil {
	// 	LogString("Error reading state after write:" + err.Error())
	// } else {
	// 	LogString("Final State Read: " + string(state))
	// }

	// // Generate an ED25519 key pair
	// priv, pub := GenerateEd25519Key()

	// // Message to be signed
	// message := []byte("Test message for ED25519")

	// // Sign the message
	// signature := SignMessageEd25519(priv, message)

	// // Verify the signature
	// if VerifyMessageEd25519(pub, message, signature) {
	// 	LogString("Signature verified successfully.")
	// } else {
	// 	LogString("Signature verification failed.")
	// }

	// // Print private key, public key, and signature for testing purposes
	// LogString(hex.EncodeToString(priv.Seed()))
	// LogString(hex.EncodeToString(signature))

	// LogString("Contract Execution Flow Completed.")

}

// //go:export helloworld
// func helloworld() {
// //Get input as json
// input, inputType, err := ContractInput()
// if err != nil {
// 	LogString("Error while getting Smart Contract Input")
// }
// LogString("Detected current : " + inputType)
// LogString(string(input))
// //Basic Serialization into bytes
// parser := NewParser(input)
// LogString("Parsed Data:" + string(parser.data))
// //Use parser.GetType(GetRawBytes,GetString,GetInt ... etc) for getting the correct type from json
// name, nameErr := parser.GetRawBytes("name")
// if nameErr != nil {
// 	LogString("Error while getting name key")
// }
// LogString("name:" + string(name))
// LogString("name hex:" + hex.EncodeToString(name))
// LogString("name base64:" + base64.RawStdEncoding.EncodeToString(name))
// LogString("name base58:" + base58.Encode(name))
// //Save raw bytes into storage
// key := []byte("person")
// result, writeErr := ContractStorageWrite(key, parser.data)
// if writeErr != nil {
// 	LogString("Error while write to the name key")
// }
// LogString("result:" + strconv.FormatBool(result))
// //Read
// data, readErr := ContractStorageRead(key)
// LogString("Write Data:" + string(data))

// if readErr != nil {
// 	LogString("Error while reading name key")
// }
// LogString("Read Data:" + string(data))
// //Deserialize
// deserializeParser := NewParser(data)
// LogString("Deserialize Data:" + string(deserializeParser.data))
// age, ageErr := parser.GetRawBytes("age")
// if ageErr != nil {
// 	LogString("Error while getting age key")
// }
// LogString("age:" + string(age))

// //Return Result
// ContractValueReturn(age)
// }

func main() {
}
