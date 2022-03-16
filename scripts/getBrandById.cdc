import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"

pub fun main(brandId:UInt64): AnyStruct{
    return TroonAtomicStandardContract.getBrandById(brandId: brandId)
}