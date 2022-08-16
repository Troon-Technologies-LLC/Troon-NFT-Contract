
import FungibleToken from 0xee82856bf20e2aa6
import NFTContract from 0x179b6b1cb6755e31
pub fun main(account:Address): Bool{
    let account1 = getAccount(account)
    let cap = account1.getCapability<&{NFTContract.NFTContractCollectionPublic}>(NFTContract.CollectionPublicPath)
    return cap.check()
}