import TroonAtomicStandard from 0x3a57788afdda9ea7
import NonFungibleToken from 0x631e88ae7f1d7c20

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
        signer.save( <- collection, to: TroonAtomicStandard.CollectionStoragePath)
        // create a public capability for the Collection
        signer.link<&{TroonAtomicStandard.TroonAtomicStandardCollectionPublic}>(TroonAtomicStandard.CollectionPublicPath, target:TroonAtomicStandard.CollectionStoragePath)
        log("ok")
    }
}