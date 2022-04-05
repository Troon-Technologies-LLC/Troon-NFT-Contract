import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


transaction() {
    prepare(signer: AuthAccount) {
        // save the resource to the signer's account storage
        if signer.getLinkTarget(TroonAtomicStandard.NFTMethodsCapabilityPrivatePath) == nil {
            let adminResouce <- TroonAtomicStandard.createAdminResource()
            signer.save(<- adminResouce, to: TroonAtomicStandard.AdminResourceStoragePath)
            // link the UnlockedCapability in private storage
            signer.link<&{TroonAtomicStandard.NFTMethodsCapability}>(
                TroonAtomicStandard.NFTMethodsCapabilityPrivatePath,
                target: TroonAtomicStandard.AdminResourceStoragePath
            )
        }

        signer.link<&{TroonAtomicStandard.UserSpecialCapability}>(
            /public/UserSpecialCapability,
            target: TroonAtomicStandard.AdminResourceStoragePath
        )

        let collection  <- TroonAtomicStandard.createEmptyCollection()
        // store the empty NFT Collection in account storage
        signer.save( <- collection, to:TroonAtomicStandard.CollectionStoragePath)
        // create a public capability for the Collection
        signer.link<&{NonFungibleToken.CollectionPublic}>(TroonAtomicStandard.CollectionPublicPath, target:TroonAtomicStandard.CollectionStoragePath)
    }
}