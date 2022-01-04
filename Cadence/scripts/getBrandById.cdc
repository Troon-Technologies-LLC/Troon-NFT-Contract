import NFTContract from "./NFTContract.cdc"
import NonFungibleToken from "./NonFungibleToken.cdc"

pub fun main(brandId:UInt64): AnyStruct{
    return NFTContract.getBrandById(brandId: brandId)
}