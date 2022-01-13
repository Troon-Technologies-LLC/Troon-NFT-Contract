import NFTContract from "../contracts/NFTContract.cdc"


pub fun main(): UInt64 {
    return NFTContract.totalSupply

}