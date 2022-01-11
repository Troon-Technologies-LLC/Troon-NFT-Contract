import NFTContract from "./NFTContract.cdc"
import NonFungibleToken from "./NonFungibleToken.cdc"

pub fun main():{UInt64:NFTContract.Template}  {
    return NFTContract.getAllTemplates()
}
