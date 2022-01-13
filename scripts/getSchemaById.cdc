import NFTContract from "../contracts/NFTContract.cdc"


pub fun main(schemaId: UInt64): NFTContract.Schema {
    return NFTContract.getSchemaById(schemaId: schemaId)
}