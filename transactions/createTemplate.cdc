import XGStudio from "../contracts/XGStudio.cdc"

transaction(brandId:UInt64, maxSupply:UInt64) {
    prepare(acct: AuthAccount) {

        let actorResource = acct.getCapability
            <&{XGStudio.NFTMethodsCapability}>
            (XGStudio.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")

        let mutableData : {String: AnyStruct} = {   
            "title": "racing car NFT",
            "description":  "wining moment of game"
        }
        let immutableData : {String: AnyStruct} = {
            "contentType" : "Image",
            "contectUrl"  : "https://xgstudios.io",
            "title"       : "Second NFT",
            "description" : "Second NFT for the xgStudio",
            "nftType"     : "AR",
            "gender"      : "Male",
            "raceName"    : "Lion",
            "raceDate":   1649994582.0 as Fix64,
            "raceDescription": "Lion race",
            "raceLocation" : "Mian Essa"   
        }
        actorResource.createTemplate(brandId: brandId, maxSupply: maxSupply, immutableData: immutableData, mutableData: mutableData)
        log("Template created")
    }
}
