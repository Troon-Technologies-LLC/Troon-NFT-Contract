## Directory Structure

The directories here are organized into contracts, scripts, and transactions.

Contracts contain the source code for the NFTContract and Nowwhere that are deployed to Flow.

Scripts contain read-only transactions to get information about
the state of someones Collection or about the state of the NFTcontract and Nowwhere.

Transactions contain the transactions that various users can use
to perform actions in the smart contract like creating Collection, Schema, Templates and Mint Templates.

 - `contracts/` : Where the NFTContract and Nowwhere smart contracts live.
 - `transactions/` : This directory contains all the transactions and scripts
 that are associated with these smart contracts.
 - `scripts/`  : This contains all the read-only Cadence scripts 
 that are used to read information from the smart contract
 or from a resource in account storage.
 - `go/` : This directory contains packages for golang programming languages
 to be able to read copies of the Nowwhere smart contracts, transaction templates,
 and scripts. Also contains automated tests written in those languages. Currently,
 Go is the only language that is supported, but we are hoping to add javascript
 and other languages soon. See the README in `go/` for more information
 about how to use the Go packages.