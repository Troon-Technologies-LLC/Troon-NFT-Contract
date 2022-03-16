import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"


// Print the Collection owned by accounts 0x01
pub fun main() : Int {
  return  TroonAtomicStandard.getAllBrands().length
}