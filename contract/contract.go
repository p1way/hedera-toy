package main

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

func main() {

	//start the client
	client := startClient()

	//Import and parse the compiled contract from the HelloHedera.json file
	//bytecode := getBytecode("./HelloHedera.json")

	//Create a file on Hedera and store the hex-encoded bytecode
	//bytecodeFileID := createFile(bytecode, client)
	bytecodeFileID, _ := hedera.FileIDFromString("0.0.34746510")
	fmt.Printf("contract bytecode file ID: %v\n", bytecodeFileID)

	//Deploy the contract
	//contractID := deploy(bytecodeFileID, client, hedera.NewContractFunctionParameters().
	//AddString("Hello from hedera"))
	contractID, _ := hedera.ContractIDFromString("0.0.34746511")
	fmt.Printf("contract ID: %v\n", contractID)

	//time.Sleep(1 * time.Second)

	//Query the contract message before change
	result := query(contractID, client, "getMessage", nil)
	// Get a string from the result at index 0
	message := result.GetString(0)
	fmt.Println("The contract message holds: ", message)

	//Send tx to change the contract message
	status := execute(contractID, client, "setMessage", hedera.NewContractFunctionParameters().
		AddString("Hello from Hedera again!"))
	fmt.Println("The transaction status is", status)

	time.Sleep(1 * time.Second)

	//Query the contract message before change
	//result2 := query(contractID, client, "getMessage", nil)
	result2 := query(contractID, client, "getMessage", nil)
	message2 := result2.GetString(0)
	fmt.Println("The contract message holds: ", message2)
}
