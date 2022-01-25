import NFTContract from "./NFTContract.cdc"

pub fun main(brandId:UInt64): NFTContract.Brand {
    return NFTContract.getBrandById(brandId: brandId)
}