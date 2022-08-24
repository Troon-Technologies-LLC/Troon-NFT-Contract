import MetadataViews from 0x631e88ae7f1d7c20
import NonFungibleToken from 0x631e88ae7f1d7c20
import NFTContract from 0xeb1e6c075991c1a7
transaction(){
    prepare(account: AuthAccount){
        
        if account.borrow<&NFTContract.Collection>(from: NFTContract.CollectionStoragePath) !=nil {
            return 
        }
        let collection <- NFTContract.createEmptyCollection() as! @NFTContract.Collection

        account.save(<- collection, to: NFTContract.CollectionStoragePath)

        account.link<&{NonFungibleToken.CollectionPublic, NFTContract.NFTContractCollectionPublic, MetadataViews.ResolverCollection}>(
            NFTContract.CollectionPublicPath,
            target: NFTContract.CollectionStoragePath)
    }
}