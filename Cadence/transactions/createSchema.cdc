import NFTContract from "./NFTContract.cdc"
import NonFungibleToken from "./NonFungibleToken.cdc"
transaction (schemaName:String, author:Address){

   prepare(acct: AuthAccount) {
      let actorResource = acct.getCapability
            <&{NFTContract.NFTMethodsCapability}>
            (NFTContract.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")

         let format : {String: NFTContract.SchemaType} = {
            "artist" : NFTContract.SchemaType.String,
            "artistEmail"  :  NFTContract.SchemaType.String,
            "title":NFTContract.SchemaType.String,
            "mintType":  NFTContract.SchemaType.String,
            "nftType":  NFTContract.SchemaType.String,
            "rarity":  NFTContract.SchemaType.String,
            "contectType":  NFTContract.SchemaType.String,
            "contectValue":  NFTContract.SchemaType.String,
            "extras": NFTContract.SchemaType.Any
            }

         actorResource.createSchema(schemaName: "series-1", format: format, author: 0xf3fcd2c1a78f5eee)
   }
}