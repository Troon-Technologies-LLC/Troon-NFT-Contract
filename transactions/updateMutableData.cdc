import XGStudio from "../contracts/XGStudio.cdc"
transaction (templateId:UInt64){
    prepare(acct: AuthAccount) {

        let actorResource = acct.getCapability
            <&{XGStudio.NFTMethodsCapability}>
            (XGStudio.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")
        var key = "game"
        var value = "cricket"

    actorResource.updateTemplateMutableAttribute(templateId: templateId, key: key, value: value)

    log("template mutable attribute updated")

    }
}