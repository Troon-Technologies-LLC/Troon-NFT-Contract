
import DapperUtilityCoin from 0x82ec283f88a62e65
import FungibleToken from 0x9a0766d93b6608b7

transaction(){
    prepare(account:AuthAccount){

        if account.borrow<&DapperUtilityCoin.Vault>(from: /storage/dapperUtilityCoinVault) == nil{
            account.save(<-DapperUtilityCoin.createEmptyVault(), to:/storage/dapperUtilityCoinVault)
            account.link<&DapperUtilityCoin.Vault{FungibleToken.Receiver}>(
                /storage/dapperUtilityCoinVault,
                target:/public/dapperUtilityCoinReceiver)

            account.link<&DapperUtilityCoin.Vault{FungibleToken.Balance}>(
                /storage/dapperUtilityCoinVault,
                target:/public/dapperUtilityCoinBalance)
        }

    }
}