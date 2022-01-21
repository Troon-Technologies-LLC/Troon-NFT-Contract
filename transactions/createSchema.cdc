import NFTContractV01 from "../contracts/NFTContractV01.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"

transaction (schemaName:String){

   prepare(acct: AuthAccount) {
      let actorResource = acct.getCapability
            <&{NFTContractV01.NFTMethodsCapability}>
            (NFTContractV01.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")

         let format : {String: NFTContractV01.SchemaType} = {
            "artist" : NFTContractV01.SchemaType.String,
            "artistEmail"  :  NFTContractV01.SchemaType.String,
            "title":NFTContractV01.SchemaType.String,
            "mintType":  NFTContractV01.SchemaType.String,
            "nftType":  NFTContractV01.SchemaType.String,
            "rarity":  NFTContractV01.SchemaType.String,
            "contectType":  NFTContractV01.SchemaType.String,
            "contectValue":  NFTContractV01.SchemaType.String,
            "extras": NFTContractV01.SchemaType.Any
            }

         actorResource.createSchema(schemaName: schemaName, format: format)
   }
}