import NFTContract from "./NFTContract.cdc"

pub fun main() : Int {
  return  NFTContract.getAllSchemas().length
}
