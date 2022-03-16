import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"


pub fun main(schemaId: UInt64): TroonAtomicStandardContract.Schema {
    return TroonAtomicStandardContract.getSchemaById(schemaId: schemaId)
}