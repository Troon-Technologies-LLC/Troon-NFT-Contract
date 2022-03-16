import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"


pub fun main(brandId:UInt64): TroonAtomicStandard.Brand {
    return TroonAtomicStandard.getBrandById(brandId: brandId)
}