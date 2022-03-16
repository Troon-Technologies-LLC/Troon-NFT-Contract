import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"

pub fun main(brandId:UInt64): AnyStruct{
    return TroonAtomicStandard.getBrandById(brandId: brandId)
}