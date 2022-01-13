import NFTContract from "../contracts/NFTContract.cdc"
import NonFungibleToken from "../contracts/NonFungibleToken.cdc"


transaction (brandName:String, data:{String:String}){
    prepare(acct: AuthAccount) {
        let actorResource = acct.getCapability
            <&{NFTContract.NFTMethodsCapability}>
            (NFTContract.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")

        actorResource.createNewBrand(
        brandName: brandName,
        data: data)
    
    }
}