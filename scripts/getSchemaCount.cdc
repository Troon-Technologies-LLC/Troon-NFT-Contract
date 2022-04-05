import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"

pub fun main() : Int {
  return  TroonAtomicStandard.getAllSchemas().length
}
