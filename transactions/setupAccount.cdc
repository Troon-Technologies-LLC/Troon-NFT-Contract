import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


transaction {
    prepare(acct: AuthAccount) {

        let collection  <- TroonAtomicStandard.createEmptyCollection()
        // store the empty NFT Collection in account storage
        acct.save( <- collection, to:TroonAtomicStandard.CollectionStoragePath)
        log("Collection created for account".concat(acct.address.toString()))
        // create a public capability for the Collection
        acct.link<&{TroonAtomicStandard.TroonAtomicStandardCollectionPublic}>(TroonAtomicStandard.CollectionPublicPath, target:TroonAtomicStandard.CollectionStoragePath)        
    }
}