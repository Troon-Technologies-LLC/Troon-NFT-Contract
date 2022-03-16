import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"


pub fun main(schemaId: UInt64): TroonAtomicStandard.Schema {
    return TroonAtomicStandard.getSchemaById(schemaId: schemaId)
}