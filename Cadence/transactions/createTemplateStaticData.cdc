import NFTContract from "./NFTContract.cdc"
import NonFungibleToken from "./NonFungibleToken.cdc"

transaction(brandId:UInt64, schemaId:UInt64, maxSupply:UInt64) {
    prepare(acct: AuthAccount) {

       let actorResource = acct.getCapability
            <&{NFTContract.NFTMethodsCapability}>
            (NFTContract.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")


        let extra : {String: AnyStruct} = {
                "name":"alex", // string
                "age":21,// integer
                "percentage":2.1 as Fix64, // address
                "owner": 0x01 as Address, // bool
                "burnable":false,
                "startDate":"",
                "endDate":""             
        }
        
        let immutableData : {String: AnyStruct} = {
            "artist" : "Nasir And Sham",
            "artistEmail"  :  "sham&nasir@gmai.com",
            "title":"First NFT",
            "mintType":  "MintOnSale",
            "nftType":  "AR",
            "rarity":  "Epic",
            "contectType":  "Image",
            "contectValue": "https://troontechnologies.com/",
            "extras": extra        
        }
        actorResource.createTemplate(brandId: brandId, schemaId: schemaId, maxSupply: maxSupply, immutableData: immutableData)
        log("Template created")
    }
}
