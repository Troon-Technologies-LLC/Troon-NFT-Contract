import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"

pub fun main(): {UInt64:TroonAtomicStandardContract.Brand} {
    return TroonAtomicStandardContract.getAllBrands()
}