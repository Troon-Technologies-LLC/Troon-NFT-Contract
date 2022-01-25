import NFTContract from "../contracts/NFTContract.cdc"


// Print the Collection owned by accounts 0x01
pub fun main() : Int {
  return  NFTContract.getAllBrands().length
}