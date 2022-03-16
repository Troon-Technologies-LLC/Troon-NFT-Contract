import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"

transaction(templateId: UInt64, account:Address){

    prepare(acct: AuthAccount) {
        let actorResource = acct.getCapability
        <&{TroonAtomicStandardContract.NFTMethodsCapability}>
        (TroonAtomicStandardContract.NFTMethodsCapabilityPrivatePath)
        .borrow() ?? 
        panic("could not borrow a reference to the NFTMethodsCapability interface")
        actorResource.mintNFT(templateId: templateId, account: account) 
    }
}