import NFTContractV01 from "../contracts/NFTContractV01.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


transaction() {
    prepare(signer: AuthAccount) {
        // save the resource to the signer's account storage
        if signer.getLinkTarget(NFTContractV01.NFTMethodsCapabilityPrivatePath) == nil {
            let adminResouce <- NFTContractV01.createAdminResource()
            signer.save(<- adminResouce, to: NFTContractV01.AdminResourceStoragePath)
            // link the UnlockedCapability in private storage
            signer.link<&{NFTContractV01.NFTMethodsCapability}>(
                NFTContractV01.NFTMethodsCapabilityPrivatePath,
                target: NFTContractV01.AdminResourceStoragePath
            )
        }

        signer.link<&{NFTContractV01.UserSpecialCapability}>(
            /public/UserSpecialCapability,
            target: NFTContractV01.AdminResourceStoragePath
        )

        let collection  <- NFTContractV01.createEmptyCollection()
        // store the empty NFT Collection in account storage
        signer.save( <- collection, to:NFTContractV01.CollectionStoragePath)
        // create a public capability for the Collection
        signer.link<&{NonFungibleToken.CollectionPublic}>(NFTContractV01.CollectionPublicPath, target:NFTContractV01.CollectionStoragePath)
    }
}