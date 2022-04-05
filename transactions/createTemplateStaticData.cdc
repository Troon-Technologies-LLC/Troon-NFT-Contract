import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


transaction(brandId:UInt64, schemaId:UInt64, maxSupply:UInt64) {
    prepare(acct: AuthAccount) {

        let actorResource = acct.getCapability
            <&{TroonAtomicStandard.NFTMethodsCapability}>
            (TroonAtomicStandard.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")


        let extra : {String: AnyStruct} = {
                "name":"alex" // string       
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
