import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"
transaction(templateId: UInt64) {
    prepare(acct: AuthAccount) {
        let actorResource = acct.getCapability
            <&{TroonAtomicStandard.NFTMethodsCapability}>
            (TroonAtomicStandard.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")
        actorResource.removeTemplateById(templateId: templateId)
        log("ok")
    }
}