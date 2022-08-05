import XGStudio from "../contracts/XGStudio.cdc"

transaction (brandName:String){
    prepare(acct: AuthAccount) {
        let actorResource = acct.getCapability
            <&{XGStudio.NFTMethodsCapability}>
            (XGStudio.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")
            let data  : {String:String} = {
                "name":"XGStudio",
                "description":"xG® rewards athletes’ real world sports participation with personalised digital collectibles and the xG® utility token.",
                "url":"https://xgstudios.io"   
        }
        actorResource.createNewBrand(
            brandName: brandName,
            data: data)
            log("brand created:")
    
    }
}