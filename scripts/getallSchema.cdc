import NFTContract from "./NFTContract.cdc"
import NonFungibleToken from "./NonFungibleToken.cdc"

pub fun main(): {UInt64:NFTContract.Schema} {
    return NFTContract.getAllSchemas()
}