import XGStudio from "../contracts/XGStudio.cdc"

transaction (brandId:UInt64,brandName:String){
  prepare(acct: AuthAccount) {

    let actorResource = acct.getCapability
              <&{XGStudio.NFTMethodsCapability}>
              (XGStudio.NFTMethodsCapabilityPrivatePath)
              .borrow() ?? 
              panic("could not borrow a reference to the NFTMethodsCapability interface")

    actorResource.updateBrandData(
      brandId: brandId,
      data:  {
          "brandName":brandName
      })
  }
}