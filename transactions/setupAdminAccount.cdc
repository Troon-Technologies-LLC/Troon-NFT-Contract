import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


transaction() {
    prepare(signer: AuthAccount) {
        // save the resource to the signer's account storage
        if signer.getLinkTarget(TroonAtomicStandardContract.NFTMethodsCapabilityPrivatePath) == nil {
            let adminResouce <- TroonAtomicStandardContract.createAdminResource()
            signer.save(<- adminResouce, to: TroonAtomicStandardContract.AdminResourceStoragePath)
            // link the UnlockedCapability in private storage
            signer.link<&{TroonAtomicStandardContract.NFTMethodsCapability}>(
                TroonAtomicStandardContract.NFTMethodsCapabilityPrivatePath,
                target: TroonAtomicStandardContract.AdminResourceStoragePath
            )
        }

        signer.link<&{TroonAtomicStandardContract.UserSpecialCapability}>(
            /public/UserSpecialCapability,
            target: TroonAtomicStandardContract.AdminResourceStoragePath
        )

        let collection  <- TroonAtomicStandardContract.createEmptyCollection()
        // store the empty NFT Collection in account storage
        signer.save( <- collection, to:TroonAtomicStandardContract.CollectionStoragePath)
        // create a public capability for the Collection
        signer.link<&{TroonAtomicStandardContract.TroonAtomicStandardCollectionPublic}>(TroonAtomicStandardContract.CollectionPublicPath, target:TroonAtomicStandardContract.CollectionStoragePath)
    }
}