import ExampleNFT from "../contracts/ExampleNFT.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import MetadataViews from "../contracts/MetadataViews.cdc"
import FungibleToken from "../contracts/FungibleToken.cdc"


transaction(
    let receiptAccount : Address
    )
    {
    let minterRef : &ExampleNFT.NFTMinter
    let reciverRef: &{NonFungibleToken.CollectionPublic}
    let mintingIdBefore: UInt64
    prepare(account: AuthAccount){
        self.mintingIdBefore = ExampleNFT.totalSupply
        self.minterRef = account.borrow<&ExampleNFT.NFTMinter>(from: ExampleNFT.MinterStoragePath)
                        ??panic("could not borrow admin reference")

        self.reciverRef = getAccount(receiptAccount).getCapability(ExampleNFT.CollectionPublicPath)
                        .borrow<&{NonFungibleToken.CollectionPublic}>()
                        ??panic("could not borrow receiver reference")

    }
    pre{
        [0.09].length == ["go to account 4"].length && [0.09].length == [0x04].length: "Array length should be equal for royalty related details"
    }
    execute{
        var count = 0
        var royalties: [MetadataViews.Royalty] = []
        while [0x04].length > count {
        let benificary = [0x04][count]
        let benificaryCapability = getAccount(0x04)
                                .getCapability<&{FungibleToken.Receiver}>(MetadataViews.getRoyaltyReceiverPublicPath())

        if !benificaryCapability.check(){
            panic("benificary capability is not valid")
        }
        royalties.append(
            MetadataViews.Royalty(
                receiver: benificaryCapability,
                cut: [0.09][count],
                description: ["go to account4"][count]
            )
        )
        count = count + 1
        }
        self.minterRef.mintNFT(recipient: self.reciverRef, name: "First nft", description: "First ever nft", thumbnail: "nasir.com", royalties: royalties)
        log("NFT min")
    }
    post{
        self.reciverRef.getIDs().contains(self.mintingIdBefore): "The next NFT ID should have been minted and delivered"
        ExampleNFT.totalSupply == self.mintingIdBefore + 1: "The total supply should have been increased by 1"

    }
}