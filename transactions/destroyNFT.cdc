import ExampleNFT from "../contracts/ExampleNFT.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


transaction(id: UInt64) {
    let collectionRef : &ExampleNFT.Collection

    prepare(account: AuthAccount){

        self.collectionRef = account.borrow<&ExampleNFT.Collection>(from: ExampleNFT.CollectionStoragePath)
                            ??panic("could not borrow collection reference")
    }
    execute {
    
        let nft <- self.collectionRef.withdraw(withdrawID: id)
        destroy nft
    }
    post{
        !self.collectionRef.getIDs().contains(id): "The nft with the specific id should have been del"
    }

}