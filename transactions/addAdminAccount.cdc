import NFTContractV01 from "../contracts/NFTContractV01.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


transaction(admin:Address) {
    prepare(signer: AuthAccount) {

        // get the public account object for the Admin
        let TemplateAdminAccount = getAccount(admin)

        // get the public capability from the Admin's public storage
        let TemplateAdminResource = TemplateAdminAccount.getCapability
            <&{NFTContractV01.UserSpecialCapability}>
            (/public/UserSpecialCapability)
            .borrow()
            ?? panic("could not borrow reference to UserSpecialCapability")

        // get the private capability from the Authorized owner of the AdminResource
        // this will be the signer of this transaction
        //
        let specialCapability = signer.getCapability
            <&{NFTContractV01.NFTMethodsCapability}>
            (NFTContractV01.NFTMethodsCapabilityPrivatePath) 

        // if the special capability is valid...
        if specialCapability.check() {
            // ...add it to the TemplateAdminResource
            TemplateAdminResource.addCapability(cap: specialCapability)
            log("capability added")
        } else {
            // ...let the people know we failed
            panic("special capability is invalid!")
        }
    }
}