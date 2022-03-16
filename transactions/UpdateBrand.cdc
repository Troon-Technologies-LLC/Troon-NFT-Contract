import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


transaction (brandId:UInt64,brandName:String){
  prepare(acct: AuthAccount) {

     let actorResource = acct.getCapability
              <&{TroonAtomicStandardContract.NFTMethodsCapability}>
              (TroonAtomicStandardContract.NFTMethodsCapabilityPrivatePath)
              .borrow() ?? 
              panic("could not borrow a reference to the NFTMethodsCapability interface")

    actorResource.updateBrandData(
      brandId: brandId,
      data:  {
          "brandName":brandName
      })
  }
}