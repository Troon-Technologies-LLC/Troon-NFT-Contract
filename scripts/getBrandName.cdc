import NFTContract from "./NFTContract.cdc"

pub fun main(brandId:UInt64): String? {
    return NFTContract.getBrandById(brandId: brandId).data["brandName"]
}