import NonFungibleToken from 0x631e88ae7f1d7c20
import FungibleToken from 0x9a0766d93b6608b7

pub contract HelloWorld: NonFungibleToken {
    
    // Events
    pub event ContractInitialized()
    pub event Withdraw(id: UInt64, from: Address?)
    pub event Deposit(id: UInt64, to: Address?)
    pub event NFTMinted(nftId: UInt64, name: String,receiptAccount: Address)

    //Paths
    pub let CollectionStoragePath: StoragePath
    pub let CollectionPublicPath: PublicPath
    pub let AdminResourceStoragePath: StoragePath

    //global variables
    pub var totalSupply: UInt64

    pub resource NFT: NonFungibleToken.INFT{
        pub let id: UInt64
        pub let name: String


        init(name: String){
            self.id = HelloWorld.totalSupply
            self.name = name
            emit NFTMinted(nftId: self.id, name: name, receiptAccount: self.owner!.address)
        }
        destroy(){

        }
    }
    pub resource interface HelloWorldCollectionPublic {
        pub fun deposit(token: @NonFungibleToken.NFT)
        pub fun getIDs():[UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
        pub fun borrowHelloWorld_NFT(id: UInt64): &HelloWorld.NFT?
    }

    pub resource Collection: HelloWorldCollectionPublic, NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.CollectionPublic {
        pub var ownedNFTs: @{UInt64: NonFungibleToken.NFT}

        pub fun withdraw(withdrawID: UInt64): @NonFungibleToken.NFT{
            let token <- self.ownedNFTs.remove(key:withdrawID)
                ??panic("cannot withdraw: Id does not exist in collection")
                emit Withdraw(id: token.id, from: self.owner?.address)
                return <- token
        }
        pub fun getIDs():[UInt64]{
            return self.ownedNFTs.keys
        }
        pub fun deposit(token: @NonFungibleToken.NFT){
            let token <- token as! @HelloWorld.NFT
            let id = token.id
            let oldToken <-self.ownedNFTs[id]<- token

            emit Deposit(id: id, to: self.owner?.address)

            destroy <- oldToken
        }
        pub fun borrowNFT(id: UInt64):&NonFungibleToken.NFT{
            return (&self.ownedNFTs[id] as &NonFungibleToken.NFT?)!
        }
        pub fun borrowHelloWorld_NFT(id:UInt64): &HelloWorld.NFT?{
            if self.ownedNFTs[id] != nil{
                let ref = (&self.ownedNFTs[id] as auth &NonFungibleToken.NFT?)!
                return ref as! &HelloWorld.NFT
            }else {
                return nil
            }
        }
        init(){
            self.ownedNFTs <- {}
        }
    }
    pub resource Admin {
        pub fun mintNFT(receiverAddress: Address, name: String){
            pre {
                name.length>0:"name must be valid"
            }
            let receiverAccount = getAccount(receiverAddress)
            let receiverCollection = receiverAccount.getCapability(HelloWorld.CollectionPublicPath)
                .borrow<&{HelloWorld.HelloWorldCollectionPublic}>()
                ??panic("could not borrow reciver reference")
            var newNFT <- create NFT(name: name)
            receiverCollection.deposit(token: <- newNFT)
        }
    }
    pub fun purchaseNFT(payment: @FungibleToken.Vault){

    }
    //method to create empty Collection
    pub fun createEmptyCollection(): @NonFungibleToken.Collection {
        return <- create HelloWorld.Collection()
    }

    init() {
        self.totalSupply = 0
        self.CollectionStoragePath = /storage/HelloWorldCollection
        self.CollectionPublicPath = /public/HelloWorldCollection
        self.AdminResourceStoragePath = /storage/HelloWorldAdmin

        self.account.save(<- create Admin(), to: /storage/AdminResourceStoragePath)
    }
}