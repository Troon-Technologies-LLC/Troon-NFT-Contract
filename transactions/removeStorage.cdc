import NFTContract from 0xd4221a1979538992

// This transaction calls the "hello" method on the HelloAsset object
// that is stored in the account's storage by removing that object
// from storage, calling the method, and then putting it back in storage
//transaction{
//    prepare(account:AuthAccount){
//    let contract = account.contracts.remove(name: "NFTContract")
//    }
//}
transaction {

    prepare(acct: AuthAccount) {

        // load the resource from storage, specifying the type to load it as
        // and the path where it is stored
        let helloResource <- acct.load<@NFTContract.AdminResource>(from: NFTContract.CollectionStoragePath)

        destroy helloResource
    }
}