import XGStudio from "../contracts/XGStudio.cdc"


// Print the NFTs owned by accounts 0x01 and 0x02.
pub fun main(address: Address) : [UInt64] {

    // Get both public account objects
    let account1 = getAccount(address)
    // Find the public Receiver capability for their Collections
    let acct1Capability =  account1.getCapability(XGStudio.CollectionPublicPath)

                            .borrow<&{XGStudio.XGStudioCollectionPublic}>()
                            ??panic("could not borrow receiver reference ")
    return  acct1Capability.getIDs()

}