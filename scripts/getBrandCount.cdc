import XGStudio from "../contracts/XGStudio.cdc"

// Print the Collection owned by accounts 0x01
pub fun main() : Int {
  return  XGStudio.getAllBrands().length
}