import XGStudio from "../contracts/XGStudio.cdc"

pub fun main(brandId:UInt64): AnyStruct{
    return XGStudio.getBrandById(brandId: brandId)
}