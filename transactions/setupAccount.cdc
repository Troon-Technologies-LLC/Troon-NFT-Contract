import ExampleNFT from "../contracts/ExampleNFT.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import MetadataViews from "../contracts/MetadataViews.cdc"

transaction(){
    prepare(account: AuthAccount){
        
        if account.borrow<&ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath) !=nil {
            return 
        }
        let collection <- ExampleNFT.createEmptyCollection()

        account.save(<- collection, to: ExampleNFT.CollectionStoragePath)

        account.link<&{NonFungibleToken.CollectionPublic, ExampleNFT.ExampleNFTCollectionPublic, MetadataViews.ResolverCollection}>(
            ExampleNFT.CollectionPublicPath,
            target: ExampleNFT.CollectionStoragePath)
    }
}