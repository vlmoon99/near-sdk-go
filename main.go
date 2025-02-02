package main

// 0.Create structs
type Person struct {
	Name  string
	Age   int
	Email string
}

//go:export helloworld
func helloworld() {
	//Get input as json
	input, inputType, err := ContractInput()
	if err != nil {
		LogString("Error while getting Smart Contract Input")
	}
	LogString("Detected current : " + inputType)
	LogString(string(input))

	//Basic Serialization into bytes
	parser := NewParser(input)
	LogString("Parsed Data:" + string(parser.data))
	//Use parser.GetType(GetRawBytes,GetString,GetInt ... etc) for getting the correct type from json

	name, nameErr := parser.GetRawBytes("name")
	if nameErr != nil {
		LogString("Error while getting name key")
	}
	LogString("name:" + (string(name)))

	//Save raw bytes into storage
	key := []byte("person")
	ContractStorageWrite(key, parser.data)
	//Read
	data, readErr := ContractStorageRead(key)
	LogString("Write Data:" + (string(data)))

	if readErr != nil {
		LogString("Error while reading name key")
	}
	LogString("Read Data:" + (string(data)))

	//Deserialize
	deserializeParser := NewParser(data)
	LogString("Read Data:" + (string(deserializeParser.data)))

	age, ageErr := parser.GetRawBytes("age")
	if ageErr != nil {
		LogString("Error while getting age key")
	}
	LogString("age:" + (string(age)))

	//Return Result
	ContractValueReturn(age)
}

func main() {
}
