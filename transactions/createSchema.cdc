import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"

transaction (schemaName:String){

   prepare(acct: AuthAccount) {
      let actorResource = acct.getCapability
            <&{TroonAtomicStandard.NFTMethodsCapability}>
            (TroonAtomicStandard.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")

         let format : {String: TroonAtomicStandard.SchemaType} = {
            "artist" : TroonAtomicStandard.SchemaType.String,
            "artistEmail"  :  TroonAtomicStandard.SchemaType.String,
            "title":TroonAtomicStandard.SchemaType.String,
            "mintType":  TroonAtomicStandard.SchemaType.String,
            "nftType":  TroonAtomicStandard.SchemaType.String,
            "rarity":  TroonAtomicStandard.SchemaType.String,
            "contectType":  TroonAtomicStandard.SchemaType.String,
            "contectValue":  TroonAtomicStandard.SchemaType.String,
            "extras": TroonAtomicStandard.SchemaType.Any
            }

         actorResource.createSchema(schemaName: schemaName, format: format)
   }
}