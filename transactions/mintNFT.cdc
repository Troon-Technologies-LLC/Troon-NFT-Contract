import NFTContract from 0x179b6b1cb6755e31
import NonFungibleToken from 0x01cf0e2f2f715450
import MetadataViews from 0x01cf0e2f2f715450
import FungibleToken from 0xee82856bf20e2aa6

transaction(templateId: UInt64, reciptAddress: Address){
    prepare(account: AuthAccount){
        let actorResource = account.getCapability
                                        <&{NFTContract.NFTMethodsCapability}>
                                        (NFTContract.NFTMethodsCapabilityPrivatePath)
                                        .borrow() ?? 
                                        panic("could not borrow a reference to the NFTMethodsCapability interface")
            let immutableData : {String: AnyStruct} = {
                "name" : "Nasir"  
        }
        var royalties: [MetadataViews.Royalty] = []
        actorResource.mintNFT( 
            templateId: templateId,
            account: reciptAddress,
            immutableData:immutableData,
            name: "Imsa",
            description: "Imsa is a fastlane",
            thumbnail: "www.google.com",
            royalties: royalties)
    }
}