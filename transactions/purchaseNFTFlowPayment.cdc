import MetadataViews from 0x631e88ae7f1d7c20
import FungibleToken from 0x9a0766d93b6608b7
import NFTContract from 0xeb1e6c075991c1a7
import DapperUtilityCoin from 0x82ec283f88a62e65

transaction(templateId: UInt64, reciptAddress: Address, price: UFix64){

    var temporaryVault: @FungibleToken.Vault
    
    prepare(account:AuthAccount){

        let vaultRef = account.borrow<&DapperUtilityCoin.Vault>(from:  /storage/dapperUtilityCoinVault)
                                    ??panic("could not borrow vault ref")

        self.temporaryVault <- vaultRef.withdraw(amount: price)
    }
    execute{
        let immutableData:{String:AnyStruct}? ={
            "name":"Nasir"
        }
        var royalties: [MetadataViews.Royalty] = []
        NFTContract.purchaseNFT(
            templateId: templateId,
            account: reciptAddress,
            immutableData:immutableData,
            name: "Dapper",
            description: "Dapper transaction",
            thumbnail: "Dapper.com",
            royalties: royalties,
            payment: <- self.temporaryVault,
            price: price)
    }
}