package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/joho/godotenv"
)

func startClient() *hedera.Client {
	//Loads the .env file and throws an error if it cannot load the variables from that file corectly
	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("Unable to load enviroment variables from .env file. Error:\n%v\n", err))
	}

	//Grab your testnet account ID and private key from the .env file
	myAccountId, err := hedera.AccountIDFromString(os.Getenv("MY_ACCOUNT_ID"))
	if err != nil {
		panic(err)
	}

	myPrivateKey, err := hedera.PrivateKeyFromString(os.Getenv("MY_PRIVATE_KEY"))
	if err != nil {
		panic(err)
	}

	//Create your testnet client
	client := hedera.ClientForTestnet()
	client.SetOperator(myAccountId, myPrivateKey)

	return client
}

func getBytecode(path string) string {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		println("Error reading HelloHedera.json")
		panic(err)
	}

	var contract map[string]interface{}

	err = json.Unmarshal([]byte(file), &contract)
	if err != nil {
		println("Error unmarshaling the json file")
		panic(err)
	}

	data := contract["data"].(map[string]interface{})
	bytecode_field := data["bytecode"].(map[string]interface{})
	bytecode := bytecode_field["object"].(string)
	return bytecode
}

func createFile(bytecode string, client *hedera.Client) hedera.FileID {
	fileTx, err := hedera.NewFileCreateTransaction().
		//Set the bytecode of the contract
		SetContents([]byte(bytecode)).
		//Submit the transaction to a Hedera network
		Execute(client)

	if err != nil {
		println(err.Error(), ": error creating file")
		panic(err)
	}

	//Get the receipt of the file create transaction
	fileReceipt, err := fileTx.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error getting file create transaction receipt")
		panic(err)
	}

	//Get the file ID
	bytecodeFileID := *fileReceipt.FileID

	return bytecodeFileID
}

func deploy(id hedera.FileID, client *hedera.Client, params *hedera.ContractFunctionParameters) hedera.ContractID {
	// Instantiate the contract instance
	contractTx, err := hedera.NewContractCreateTransaction().
		//Set the file ID of the Hedera file storing the bytecode
		SetBytecodeFileID(id).
		//Set the gas to instantiate the contract
		SetGas(100000).
		//Provide the constructor parameters for the contract
		SetConstructorParameters(params).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error creating contract")
		panic(err)
	}

	//Get the receipt of the contract create transaction
	contractReceipt, err := contractTx.GetReceipt(client)
	if err != nil {
		println(err.Error(), ": error retrieving contract creation receipt")
		panic(err)
	}

	//Get the contract ID from the receipt
	contractID := *contractReceipt.ContractID
	return contractID
}

func query(id hedera.ContractID, client *hedera.Client, method string, params *hedera.ContractFunctionParameters) hedera.ContractFunctionResult {
	// Calls a function of the smart contract
	contractQuery, err := hedera.NewContractCallQuery().
		//Set the contract ID to return the request for
		SetContractID(id).
		//Set the gas for the query
		SetGas(100000).
		//Set the query payment for the node returning the request
		//This value must cover the cost of the request otherwise will fail
		SetQueryPayment(hedera.NewHbar(1)).
		//Set the contract function to call
		SetFunction(method, params). // nil -> no parameters
		//Submit the query to a Hedera network
		Execute(client)

	if err != nil {
		println(err.Error(), ": error executing contract call query")
		panic(err)
	}

	return contractQuery
}

func execute(id hedera.ContractID, client *hedera.Client, method string, params *hedera.ContractFunctionParameters) hedera.Status {
	//Create the transaction to update the contract message
	contractExecTx, err := hedera.NewContractExecuteTransaction().
		//Set the ID of the contract
		SetContractID(id).
		//Set the gas to execute the call
		SetGas(100000).
		//Set the contract function to call
		SetFunction(method, params).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error executing contract")
		panic(err)
	}

	//Get the receipt of the transaction
	receipt, err := contractExecTx.GetReceipt(client)

	return receipt.Status
}
