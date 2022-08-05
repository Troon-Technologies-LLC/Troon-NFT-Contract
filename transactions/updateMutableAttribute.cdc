import XGStudio from "../contracts/XGStudio.cdc"
transaction (templateId:UInt64){
    prepare(acct: AuthAccount) {

        let actorResource = acct.getCapability
            <&{XGStudio.NFTMethodsCapability}>
            (XGStudio.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")
            let mutableData : {String: AnyStruct} = {   
                "Keyboard" : "Qwerty",
                "InputType" : "AlphaNumeric"
            }

    actorResource.updateTemplateMutableData(templateId: templateId, mutableData: mutableData)

    log("mutable data update")

    }
}