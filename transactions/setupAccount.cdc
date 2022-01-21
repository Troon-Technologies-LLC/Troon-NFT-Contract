import NFTContractV01 from "../contracts/NFTContractV01.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


transaction {
    prepare(acct: AuthAccount) {

        let collection  <- NFTContractV01.createEmptyCollection()
        // store the empty NFT Collection in account storage
        acct.save( <- collection, to:NFTContractV01.CollectionStoragePath)
        log("Collection created for account".concat(acct.address.toString()))
        // create a public capability for the Collection
        acct.link<&{NonFungibleToken.CollectionPublic}>(NFTContractV01.CollectionPublicPath, target:NFTContractV01.CollectionStoragePath)
        
    }
}