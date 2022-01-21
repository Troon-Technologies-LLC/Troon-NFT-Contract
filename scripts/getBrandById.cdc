import NFTContractV01 from "../contracts/NFTContractV01.cdc"

pub fun main(brandId:UInt64): AnyStruct{
    return NFTContractV01.getBrandById(brandId: brandId)
}