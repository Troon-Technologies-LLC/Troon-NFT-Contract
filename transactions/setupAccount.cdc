import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


transaction {
    prepare(acct: AuthAccount) {

        let collection  <- TroonAtomicStandardContract.createEmptyCollection()
        // store the empty NFT Collection in account storage
        acct.save( <- collection, to:TroonAtomicStandardContract.CollectionStoragePath)
        log("Collection created for account".concat(acct.address.toString()))
        // create a public capability for the Collection
        acct.link<&{TroonAtomicStandardContract.TroonAtomicStandardCollectionPublic}>(TroonAtomicStandardContract.CollectionPublicPath, target:TroonAtomicStandardContract.CollectionStoragePath)
        
    }
}