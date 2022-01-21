import NFTContractV01 from "../contracts/NFTContractV01.cdc"


pub fun main(schemaId: UInt64): NFTContractV01.Schema {
    return NFTContractV01.getSchemaById(schemaId: schemaId)
}