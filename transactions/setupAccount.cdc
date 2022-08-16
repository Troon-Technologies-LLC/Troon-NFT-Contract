import NFTContract from 0x179b6b1cb6755e31
import NonFungibleToken from 0x01cf0e2f2f715450
import MetadataViews from 0x01cf0e2f2f715450

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