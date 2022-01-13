import NFTContract from "../contracts/NFTContract.cdc"

pub fun main(brandId:UInt64): AnyStruct{
    return NFTContract.getBrandById(brandId: brandId)
}