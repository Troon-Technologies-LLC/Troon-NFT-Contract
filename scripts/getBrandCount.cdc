import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"


// Print the Collection owned by accounts 0x01
pub fun main() : Int {
  return  TroonAtomicStandardContract.getAllBrands().length
}