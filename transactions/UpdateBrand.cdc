import NFTContractV01 from "../contracts/NFTContractV01.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


transaction (brandId:UInt64,brandName:String){
  prepare(acct: AuthAccount) {

     let actorResource = acct.getCapability
              <&{NFTContractV01.NFTMethodsCapability}>
              (NFTContractV01.NFTMethodsCapabilityPrivatePath)
              .borrow() ?? 
              panic("could not borrow a reference to the NFTMethodsCapability interface")

    actorResource.updateBrandData(
      brandId: brandId,
      data:  {
          "brandName":brandName
      })
  }
}