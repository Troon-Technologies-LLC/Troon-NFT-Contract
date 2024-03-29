import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


pub fun main(address: Address):[UInt64]{
    let account1 = getAccount(address)
    let acct1Capability =  account1.getCapability(TroonAtomicStandard.CollectionPublicPath)
                           .borrow<&{TroonAtomicStandard.TroonAtomicStandardCollectionPublic}>()
                            ??panic("could not borrow receiver Reference ")
    return  acct1Capability.getIDs()
}