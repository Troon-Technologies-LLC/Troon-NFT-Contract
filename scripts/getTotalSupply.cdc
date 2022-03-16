import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"


pub fun main(): UInt64 {
    return TroonAtomicStandardContract.totalSupply

}