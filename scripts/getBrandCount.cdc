import NFTContractV01 from "../contracts/NFTContractV01.cdc"


// Print the Collection owned by accounts 0x01
pub fun main() : Int {
  return  NFTContractV01.getAllBrands().length
}