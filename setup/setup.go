package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/joho/godotenv"
)

func main() {

	//Loads the .env file and throws an error if it cannot load the variables from that file correctly
	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("Unable to load environment variables from .env file. Error:\n%v\n", err))
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

	//Print your testnet account ID and private key to the console to make sure there was no error
	fmt.Printf("The account ID is = %v\n", myAccountId)
	fmt.Printf("The private key is = %v\n", myPrivateKey)

	//Create your testnet client
	client := hedera.ClientForTestnet()
	client.SetOperator(myAccountId, myPrivateKey)

	//Create the account balance query
	query := hedera.NewAccountBalanceQuery().
		SetAccountID(myAccountId)

	//Sign with client operator private key and submit the query to a Hedera network
	accountBalance, err := query.Execute(client)
	if err != nil {
		panic(err)
	}

	//Print the balance of tinybars
	fmt.Println("My account balance is ", accountBalance.Hbars)

	//Generate new keys for the account you will create
	newAccountPrivateKey, err := hedera.PrivateKeyGenerateEd25519()

	if err != nil {
		panic(err)
	}

	newAccountPublicKey := newAccountPrivateKey.PublicKey()
	fmt.Printf("The public key of the new account is = %v\n", newAccountPublicKey)
	fmt.Printf("The private key of the new account is = %v\n", newAccountPrivateKey)

	//Create new account and assign the public key
	newAccount, err := hedera.NewAccountCreateTransaction().
		SetKey(newAccountPublicKey).
		SetInitialBalance(hedera.HbarFrom(1000, hedera.HbarUnits.Tinybar)).
		Execute(client)

	//Request the receipt of the transaction
	receipt, err := newAccount.GetReceipt(client)
	if err != nil {
		panic(err)
	}

	//Get the new account ID from the receipt
	newAccountId := *receipt.AccountID
	fmt.Printf("The new account ID is %v\n", newAccountId)

	//Get the tx ID
	txId := newAccount.TransactionID
	fmt.Printf("The account creation tx ID is %v\n", txId)

	//Create the account balance query
	query = hedera.NewAccountBalanceQuery().
		SetAccountID(newAccountId)

	//Sign with client operator private key and submit the query to a Hedera network
	accountBalance, err = query.Execute(client)
	if err != nil {
		panic(err)
	}

	//Print the balance of tinybars
	fmt.Println("The account balance for the new account is ", accountBalance.Hbars.AsTinybar())

	//Transfer hbar from your testnet account to the new account
	transaction := hedera.NewTransferTransaction().
		AddHbarTransfer(myAccountId, hedera.HbarFrom(-100, hedera.HbarUnits.Hbar)).
		AddHbarTransfer(newAccountId, hedera.HbarFrom(100, hedera.HbarUnits.Hbar))

		//Submit the transaction to a Hedera network
	txResponse, err := transaction.Execute(client)

	if err != nil {
		panic(err)
	}

	//Request the receipt of the transaction
	transferReceipt, err := txResponse.GetReceipt(client)

	if err != nil {
		panic(err)
	}

	//Get the transaction consensus status
	transactionStatus := transferReceipt.Status

	fmt.Printf("The transaction consensus status is %v\n", transactionStatus)

	//Create the query that you want to submit
	balanceQuery := hedera.NewAccountBalanceQuery().
		SetAccountID(newAccountId)

		//Get the cost of the query
	cost, err := balanceQuery.GetCost(client)

	if err != nil {
		panic(err)
	}

	fmt.Println("The account balance query cost is:", cost.String())

	//Sign with client operator private key and submit the query to a Hedera network
	newAccountBalance, err := balanceQuery.Execute(client)
	if err != nil {
		panic(err)
	}

	//Print the balance of tinybars
	fmt.Println("The hbar account balance for this account is", newAccountBalance.Hbars.AsTinybar())
}
