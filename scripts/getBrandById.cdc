import NFTContract from 0x179b6b1cb6755e31
pub fun main(brandId:UInt64): AnyStruct{
    return NFTContract.getBrandById(brandId: brandId)
}