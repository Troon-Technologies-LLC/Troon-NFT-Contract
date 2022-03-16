import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"

transaction (schemaName:String){

   prepare(acct: AuthAccount) {
      let actorResource = acct.getCapability
            <&{TroonAtomicStandardContract.NFTMethodsCapability}>
            (TroonAtomicStandardContract.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")

         let format : {String: TroonAtomicStandardContract.SchemaType} = {
            "artist" : TroonAtomicStandardContract.SchemaType.String,
            "artistEmail"  :  TroonAtomicStandardContract.SchemaType.String,
            "title":TroonAtomicStandardContract.SchemaType.String,
            "mintType":  TroonAtomicStandardContract.SchemaType.String,
            "nftType":  TroonAtomicStandardContract.SchemaType.String,
            "rarity":  TroonAtomicStandardContract.SchemaType.String,
            "contectType":  TroonAtomicStandardContract.SchemaType.String,
            "contectValue":  TroonAtomicStandardContract.SchemaType.String,
            "extras": TroonAtomicStandardContract.SchemaType.Any
            }

         actorResource.createSchema(schemaName: schemaName, format: format)
   }
}