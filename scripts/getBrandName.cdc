import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"


pub fun main(brandId:UInt64): TroonAtomicStandardContract.Brand {
    return TroonAtomicStandardContract.getBrandById(brandId: brandId)
}