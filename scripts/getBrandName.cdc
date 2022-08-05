import XGStudio from "../contracts/XGStudio.cdc"

pub fun main(brandId:UInt64): XGStudio.Brand {
    return XGStudio.getBrandById(brandId: brandId)
}