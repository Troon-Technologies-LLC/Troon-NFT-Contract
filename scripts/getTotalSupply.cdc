import NFTContractV01 from "../contracts/NFTContractV01.cdc"


pub fun main(): UInt64 {
    return NFTContractV01.totalSupply

}