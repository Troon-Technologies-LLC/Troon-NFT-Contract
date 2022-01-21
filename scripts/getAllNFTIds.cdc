import NFTContractV01 from "../contracts/NFTContractV01.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


pub fun main(address: Address):[UInt64]{
    let account1 = getAccount(address)
    let acct1Capability =  account1.getCapability(NFTContractV01.CollectionPublicPath)
                           .borrow<&{NonFungibleToken.CollectionPublic}>()
                            ??panic("could not borrow receiver Reference ")
    return  acct1Capability.getIDs()
}