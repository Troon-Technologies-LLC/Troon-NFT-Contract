import NFTContract from 0xeb1e6c075991c1a7
    transaction (brandName:String){
        // let adminRef : &NFTContract.Admin
    prepare(acct: AuthAccount) {
        let adminRef =  acct.borrow<&NFTContract.BrandAdmin>(from: /storage/BrandAdmin)
                            ??panic("could not borrow admin reference")

        let data  : {String:String} = {
                "name":"imsa",
                "description":"imsa rewards athletes’ real world sports participation with personalised digital collectibles and the xG® utility token.",
                "url":"https://imsa.io"   
        }
        adminRef.createNewBrand(
            brandName: brandName,
            data: data)
            log("brand created:")
    }
    execute{
    
    }
}