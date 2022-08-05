import XGStudio from "../contracts/XGStudio.cdc"

transaction() {
    prepare(signer: AuthAccount) {
    
        if signer.borrow<&XGStudio.Collection>(from:XGStudio.CollectionStoragePath) == nil{
            let collection  <- XGStudio.createEmptyCollection() as! @XGStudio.Collection
            // store the empty NFT Collection in account storage
            signer.save( <- collection, to: XGStudio.CollectionStoragePath)
            // create a public capability for the Collection
            signer.link<&{XGStudio.XGStudioCollectionPublic}>(XGStudio.CollectionPublicPath, target:XGStudio.CollectionStoragePath)
            log("Collection create")
        }
        else {
            log("capability already created")
        }
        
    }
}