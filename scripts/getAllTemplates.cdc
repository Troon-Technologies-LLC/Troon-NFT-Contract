import NFTContractV01 from "../contracts/NFTContractV01.cdc"

pub fun main():{UInt64:NFTContractV01.Template}  {
    return NFTContractV01.getAllTemplates()
}
