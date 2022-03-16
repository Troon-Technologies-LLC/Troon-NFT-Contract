import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"


pub fun main() : Int {
  return  TroonAtomicStandardContract.getAllTemplates().length
}