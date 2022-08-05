import XGStudio from "../contracts/XGStudio.cdc"


pub fun main(address: Address):[UInt64]{
    let account1 = getAccount(address)
    let acct1Capability =  account1.getCapability(XGStudio.CollectionPublicPath)
                            .borrow<&{XGStudio.XGStudioCollectionPublic}>()
                            ??panic("could not borrow receiver Reference ")
    return  acct1Capability.getIDs()
}