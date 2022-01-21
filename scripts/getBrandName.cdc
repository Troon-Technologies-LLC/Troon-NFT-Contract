import NFTContractV01 from "./NFTContractV01.cdc"

pub fun main(brandId:UInt64): NFTContractV01.Brand {
    return NFTContractV01.getBrandById(brandId: brandId)
}