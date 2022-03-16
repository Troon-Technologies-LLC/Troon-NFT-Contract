import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"

transaction(brandId:UInt64, schemaId:UInt64, maxSupply:UInt64,immutableData:{String: AnyStruct}) {
     prepare(acct: AuthAccount) {

       let actorResource = acct.getCapability
            <&{TroonAtomicStandard.NFTMethodsCapability}>
            (TroonAtomicStandard.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")

        actorResource.createTemplate(brandId: brandId, schemaId: schemaId, maxSupply: maxSupply, immutableData: immutableData)
    }
}