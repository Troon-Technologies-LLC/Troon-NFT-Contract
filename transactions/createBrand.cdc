import NFTContract from 0x179b6b1cb6755e31
    transaction (brandName:String){
    prepare(acct: AuthAccount) {
        let actorResource = acct.getCapability
            <&{NFTContract.NFTMethodsCapability}>
            (NFTContract.NFTMethodsCapabilityPrivatePath)
            .borrow() ?? 
            panic("could not borrow a reference to the NFTMethodsCapability interface")
            let data  : {String:String} = {
                "name":"imsa",
                "description":"imsa rewards athletes’ real world sports participation with personalised digital collectibles and the xG® utility token.",
                "url":"https://imsa.io"   
        }
        actorResource.createNewBrand(
            brandName: brandName,
            data: data)
            log("brand created:")
    
    }
}