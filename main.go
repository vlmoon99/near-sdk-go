package main

// 0.Create structs
type Person struct {
	Name  string
	Age   int
	Email string
}

//go:export helloworld
func helloworld() {
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
