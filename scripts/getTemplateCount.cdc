import NFTContractV01 from "../contracts/NFTContractV01.cdc"


pub fun main() : Int {
  return  NFTContractV01.getAllTemplates().length
}