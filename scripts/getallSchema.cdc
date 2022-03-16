import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"


pub fun main(): {UInt64:TroonAtomicStandardContract.Schema} {
    return TroonAtomicStandardContract.getAllSchemas()
}