## Directory Structure

The directories here are organized into contracts, scripts, and transactions.

Contracts directory contain the source code for the NFTContract and NonFungibleContract 
interface that are deployed to Flow Blockchain.

Scripts directory contain read-only transactions to get information about
the state of users Collection or about the state of the NFTContract.

Transactions directory contain the transactions that various users can use to
perform actions in the smart contract like Create Brand, Schema, Templates and Mint NFTs.

 - `contracts/` : Where the NFTContract smart contract live.
 - `transactions/` : This directory contains all the transactions that are associated 
 with these smart contracts.
 - `scripts/`  : This contains all the read-only Cadence scripts 
 that are used to read information from the smart contract
 or from a resource in account storage.
 - `test/` : This directory contains testcases in Golang and javascript. 'go' folder contain
 Golang testcases and 'js' folder contains Javascript testcases. This folder contains 
 automated tests written in both languages.  See the README in `go/` and `js/` for more information
 about how to run testcases.