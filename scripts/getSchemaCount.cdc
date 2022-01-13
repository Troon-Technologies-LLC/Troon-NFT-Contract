import NFTContract from "../contracts/NFTContract.cdc"

pub fun main() : Int {
  return  NFTContract.getAllSchemas().length
}
