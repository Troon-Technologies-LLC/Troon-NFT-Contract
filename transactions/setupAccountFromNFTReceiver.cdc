import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import FungibleToken from "../contracts/FungibleToken.cdc"
import ExampleNFT from "../contracts/ExampleNFT.cdc"
import MetadataViews from "../contracts/MetadataViews.cdc"
transaction(
    address: Address,
    publicPath: PublicPath,
    id: UInt64)
    {
    prepare(account:AuthAccount){
        let collection = getAccount(address)
                        .getCapability(ExampleNFT.CollectionPublicPath)
                        .borrow<&{NonFungibleToken.CollectionPublic, MetadataViews.ResolverCollection}>()
                        ??panic("could not borrow reference to collection")
        let resolver = collection.borrowViewResolver(id: id)!

        let nftCollectionView = resolver.resolveView(Type<MetadataViews.NFTCollectionData>())! as! MetadataViews.NFTCollectionData
        
        let emptyCollection <- nftCollectionView.createEmptyCollection()
        account.save(<- emptyCollection, to: nftCollectionView.storagePath)

        // create a public capability for the collection
        account.link<&{NonFungibleToken.CollectionPublic, ExampleNFT.ExampleNFTCollectionPublic, MetadataViews.ResolverCollection}>(
            nftCollectionView.publicPath,
            target: nftCollectionView.storagePath
        )
    
    }
    execute{
    }
}