This script should:
1.  start a client for testnet
2.  store the bytecode of an EVM smart-contract in a Hedera File
3.  deploy the smart-contract on Hedera
4.  query the contract
5.  update the contract
6.  query the contract (to check the change)

...but the script currently panic at step 4 because of INVALID_SIGNATURE.
