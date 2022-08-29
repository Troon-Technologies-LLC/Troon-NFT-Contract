import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import FungibleToken from "../contracts/FungibleToken.cdc"
import NFTContract from "../contracts/NFTContract.cdc"
import MetadataViews from "../contracts/MetadataViews.cdc"

transaction(vaultPath: StoragePath){
    prepare(account: AuthAccount){

        if account.borrow<&FungibleToken.Vault>(from: vaultPath) == nil{
            panic("A vault from the specific path does not exist")
        }

        let capability = account.link<&{FungibleToken.Balance, FungibleToken.Receiver}>
                        (MetadataViews.getRoyaltyReceiverPublicPath(), target: vaultPath)!

        // Make sure the capability is valid
        if !capability.check() {
            panic("Beneficiary capability is not valid!") 
        }
        

    }
    execute{
    }
}
 