import NFTContract from "./NFTContract.cdc"
import NonFungibleToken from "./NonFungibleToken.cdc"

pub fun main(schemaId: UInt64): NFTContract.Schema {
    return NFTContract.getSchemaById(schemaId: schemaId)
}