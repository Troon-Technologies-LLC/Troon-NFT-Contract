import NFTContract from "../contracts/NFTContract.cdc"

pub fun main():{UInt64:NFTContract.Template}  {
    return NFTContract.getAllTemplates()
}
