import NonFungibleToken from 0x01cf0e2f2f715450
import NFTContract from 0x179b6b1cb6755e31

// This script borrows an NFT from a collection
pub fun main(address: Address, id: UInt64) :UInt64{
    let account = getAccount(address)

    let collectionRef = account
        .getCapability(NFTContract.CollectionPublicPath)
        .borrow<&{NonFungibleToken.CollectionPublic}>()
        ?? panic("Could not borrow capability from public collection")

    // Borrow a reference to a specific NFT in the collection
    let _ = collectionRef.borrowNFT(id: id)
    return _.id
}