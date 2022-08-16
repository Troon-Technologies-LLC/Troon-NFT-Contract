import NonFungibleToken from 0x01cf0e2f2f715450
import MetadataViews from 0x01cf0e2f2f715450
pub contract NFTContract: NonFungibleToken {

    // Events
    pub event ContractInitialized()
    pub event Withdraw(id: UInt64, from: Address?)
    pub event Deposit(id: UInt64, to: Address?)
    pub event NFTDestroyed(id: UInt64)
    pub event NFTMinted(nftId: UInt64, templateId: UInt64, mintNumber: UInt64)
    pub event BrandCreated(brandId: UInt64, brandName: String, author: Address, data:{String: String})
    pub event BrandUpdated(brandId: UInt64, brandName: String, author: Address, data:{String: String})
    pub event TemplateCreated(templateId: UInt64, brandId: UInt64, maxSupply: UInt64)
    pub event TemplateRemoved(templateId: UInt64)
    pub event TemplateUpdated(templateId: UInt64)

    // Paths
    pub let AdminResourceStoragePath: StoragePath
    pub let NFTMethodsCapabilityPrivatePath: PrivatePath
    pub let CollectionStoragePath: StoragePath
    pub let CollectionPublicPath: PublicPath
    pub let AdminStorageCapability: StoragePath
    pub let AdminCapabilityPrivate: PrivatePath

    // Latest brand-id
    pub var lastIssuedBrandId: UInt64
    
    // Latest brand-id
    pub var lastIssuedTemplateId: UInt64

    // Total supply of all NFTs that are minted using this contract
    pub var totalSupply: UInt64
    
    // A dictionary that stores all Brands against it's brand-id.
    access(self) var allBrands: {UInt64: Brand}
    access(self) var allTemplates: {UInt64: Template}
    access(self) var allNFTs: {UInt64: NFTData}

    // Accounts ability to add capability
    access(self) var whiteListedAccounts: [Address]



    // A structure that contain all the data related to a Brand
    pub struct Brand {
        pub let brandId: UInt64
        pub let brandName: String
        pub let author: Address
        access(contract) var data: {String: String}
        
        init(brandName: String, author: Address, data: {String: String}) {
            pre {
                brandName.length > 0: "Brand name is required";
            }

            let newBrandId = NFTContract.lastIssuedBrandId
            self.brandId = newBrandId
            self.brandName = brandName
            self.author = author
            self.data = data
        }
        pub fun update(data: {String: String}) {
            self.data = data
        }
    }

    // A structure that contain all the data and methods related to Template
    pub struct Template {
        pub let templateId: UInt64
        pub let brandId: UInt64
        pub var maxSupply: UInt64
        pub var issuedSupply: UInt64
        pub var locked: Bool
        access(contract) var immutableData: {String: AnyStruct}
        access(contract) var mutableData: {String: AnyStruct}?
        init(brandId: UInt64,  maxSupply: UInt64, immutableData: {String: AnyStruct}, mutableData: {String: AnyStruct}?) {
            pre {
                NFTContract.allBrands[brandId] != nil:"Brand Id must be valid"
                maxSupply > 0 : "MaxSupply must be greater than zero"
                immutableData != nil: "ImmutableData must not be nil"
            }

            self.templateId = NFTContract.lastIssuedTemplateId
            self.brandId = brandId
            self.maxSupply = maxSupply
            self.immutableData = immutableData
            self.mutableData = mutableData
            self.issuedSupply = 0
            self.locked = false
        }
        // a method to update entire MutableData field of Template
        pub fun updateMutableData(mutableData: {String: AnyStruct}) {     
                self.mutableData = mutableData
        }

        // a method to update or add particular pair in MutableData field of Template
        pub fun updateMutableAttribute(key: String, value: AnyStruct){
            pre{
                self.mutableData != nil: "Mutable data is nil, update complete mutable data of template instead!"
                key != "": "Can't update invalid key"
            }
            self.mutableData?.insert(key: key, value)
        }
        // a method to get ImmutableData field of Template
        pub fun getImmutableData(): {String:AnyStruct} {
            return self.immutableData
        }
        
        // a method to get MutableData field of Template
        pub fun getMutableData(): {String: AnyStruct}? {
            return self.mutableData
        }
        // a method to increment issued supply for template
        access(contract) fun incrementIssuedSupply(): UInt64 {
            pre {
                self.issuedSupply < self.maxSupply: "Template reached max supply"
            }   

            self.issuedSupply = self.issuedSupply + 1
            return self.issuedSupply
        }

        // A method to lock the template
        pub fun lockTemplate(){
            if !self.locked {
                self.locked= true
            }
        }
    }

    // A structure that link template and mint-no of NFT
    pub struct NFTData {
        pub let templateID: UInt64
        pub let mintNumber: UInt64
        access(contract) var immutableData: {String: AnyStruct}?

        init(templateID: UInt64, mintNumber: UInt64, immutableData: {String: AnyStruct}?) {
            self.templateID = templateID
            self.mintNumber = mintNumber
            self.immutableData = immutableData
        }
        // a method to get the immutable data of the NFT
        pub fun getImmutableData(): {String:AnyStruct}? {
            return self.immutableData
        }
    }

    // The resource that represents the NFTContract NFTs
    // 
    pub resource NFT: NonFungibleToken.INFT,  MetadataViews.Resolver{
        pub let id: UInt64
        pub let name: String
        pub let description: String
        pub let thumbnail: String
        access(self) let royalties: [MetadataViews.Royalty]
        access(self) let metadata: {String: AnyStruct}
        access(contract) let data: NFTData

        init(
            templateID: UInt64,
            mintNumber: UInt64,
            immutableData: {String:AnyStruct}?,
            name: String,
            description: String,
            thumbnail: String,
            royalties: [MetadataViews.Royalty],
            metadata: {String: AnyStruct},) {
            NFTContract.totalSupply = NFTContract.totalSupply + 1
            self.id = NFTContract.totalSupply
            self.name = name
            self.description = description
            self.thumbnail = thumbnail
            self.royalties = royalties
            self.metadata = metadata
            NFTContract.allNFTs[self.id] = NFTData(templateID: templateID, mintNumber: mintNumber, immutableData: immutableData)
            self.data = NFTContract.allNFTs[self.id]!
            emit NFTMinted(nftId: self.id, templateId: templateID, mintNumber: mintNumber)
        }
        
        pub fun getViews(): [Type] {
            return [
                Type<MetadataViews.Display>(),
                Type<MetadataViews.Royalties>(),
                Type<MetadataViews.Editions>(),
                Type<MetadataViews.ExternalURL>(),
                Type<MetadataViews.NFTCollectionData>(),
                Type<MetadataViews.NFTCollectionDisplay>(),
                Type<MetadataViews.Serial>(),
                Type<MetadataViews.Traits>()
            ]
        }

        pub fun resolveView(_ view: Type): AnyStruct? {
            switch view {
                case Type<MetadataViews.Display>():
                    return MetadataViews.Display(
                        name: self.name,
                        description: self.description,
                        thumbnail: MetadataViews.HTTPFile(
                            url: self.thumbnail
                        )
                    )
                case Type<MetadataViews.Editions>():
                    // There is no max number of NFTs that can be minted from this contract
                    // so the max edition field value is set to nil
                    let editionInfo = MetadataViews.Edition(name: "Example NFT Edition", number: self.id, max: nil)
                    let editionList: [MetadataViews.Edition] = [editionInfo]
                    return MetadataViews.Editions(
                        editionList
                    )
                case Type<MetadataViews.Serial>():
                    return MetadataViews.Serial(
                        self.id
                    )
                case Type<MetadataViews.Royalties>():
                    return MetadataViews.Royalties(
                        self.royalties
                    )
                case Type<MetadataViews.ExternalURL>():
                    return MetadataViews.ExternalURL("https://example-nft.onflow.org/".concat(self.id.toString()))
                case Type<MetadataViews.NFTCollectionData>():
                    return MetadataViews.NFTCollectionData(
                        storagePath: NFTContract.CollectionStoragePath,
                        publicPath: NFTContract.CollectionPublicPath,
                        providerPath: /private/NFTContractCollection,
                        publicCollection: Type<&NFTContract.Collection{NFTContract.NFTContractCollectionPublic}>(),
                        publicLinkedType: Type<&NFTContract.Collection{NFTContract.NFTContractCollectionPublic,NonFungibleToken.CollectionPublic,NonFungibleToken.Receiver,MetadataViews.ResolverCollection}>(),
                        providerLinkedType: Type<&NFTContract.Collection{NFTContract.NFTContractCollectionPublic,NonFungibleToken.CollectionPublic,NonFungibleToken.Provider,MetadataViews.ResolverCollection}>(),
                        createEmptyCollectionFunction: (fun (): @NonFungibleToken.Collection {
                            return <-NFTContract.createEmptyCollection()
                        })
                    )
                case Type<MetadataViews.NFTCollectionDisplay>():
                    let media = MetadataViews.Media(
                        file: MetadataViews.HTTPFile(
                            url: "https://assets.website-files.com/5f6294c0c7a8cdd643b1c820/5f6294c0c7a8cda55cb1c936_Flow_Wordmark.svg"
                        ),
                        mediaType: "image/svg+xml"
                    )
                    return MetadataViews.NFTCollectionDisplay(
                        name: "The Example Collection",
                        description: "This collection is used as an example to help you develop your next Flow NFT.",
                        externalURL: MetadataViews.ExternalURL("https://example-nft.onflow.org"),
                        squareImage: media,
                        bannerImage: media,
                        socials: {
                            "twitter": MetadataViews.ExternalURL("https://twitter.com/flow_blockchain")
                        }
                    )
                case Type<MetadataViews.Traits>():
                    // exclude mintedTime and foo to show other uses of Traits
                    let excludedTraits = ["mintedTime", "foo"]
                    let traitsView = MetadataViews.dictToTraits(dict: self.metadata, excludedNames: excludedTraits)

                    // mintedTime is a unix timestamp, we should mark it with a displayType so platforms know how to show it.
                    let mintedTimeTrait = MetadataViews.Trait(name: "mintedTime", value: self.metadata["mintedTime"]!, displayType: "Date", rarity: nil)
                    traitsView.addTrait(mintedTimeTrait)

                    // foo is a trait with its own rarity
                    let fooTraitRarity = MetadataViews.Rarity(score: 10.0, max: 100.0, description: "Common")
                    let fooTrait = MetadataViews.Trait(name: "foo", value: self.metadata["foo"], displayType: nil, rarity: fooTraitRarity)
                    traitsView.addTrait(fooTrait)
                    
                    return traitsView

            }
            return nil
        }
        destroy(){
            emit NFTDestroyed(id: self.id)
        }
    }


    pub resource interface NFTContractCollectionPublic {
        pub fun deposit(token: @NonFungibleToken.NFT)
        pub fun getIDs(): [UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
        pub fun borrowNFTContract_NFT(id: UInt64): &NFTContract.NFT? {
            // If the result isn't nil, the id of the returned reference
            // should be the same as the argument to the function
            post {
                (result == nil) || (result?.id == id):
                    "Cannot borrow Reward reference: The ID of the returned reference is incorrect"
            }
        }
    }

    // Collection is a resource that every user who owns NFTs 
    // will store in their account to manage their NFTS
    //
    pub resource Collection: NFTContractCollectionPublic,NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.CollectionPublic, MetadataViews.ResolverCollection {
        pub var ownedNFTs: @{UInt64: NonFungibleToken.NFT}

        pub fun withdraw(withdrawID: UInt64): @NonFungibleToken.NFT {
            let token <- self.ownedNFTs.remove(key: withdrawID) 
                ?? panic("Cannot withdraw: template does not exist in the collection")
            emit Withdraw(id: token.id, from: self.owner?.address)
            return <-token
        }

        pub fun getIDs(): [UInt64] {
            return self.ownedNFTs.keys
        }

        pub fun deposit(token: @NonFungibleToken.NFT) {
            let token <- token as! @NFTContract.NFT
            let id = token.id
            let oldToken <- self.ownedNFTs[id] <- token
            if self.owner?.address != nil {
                emit Deposit(id: id, to: self.owner?.address)
            }
            destroy oldToken
        }

        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
            return (&self.ownedNFTs[id] as &NonFungibleToken.NFT?)!
        }

         // borrowNFTContract_NFT returns a borrowed reference to a NFTContract
        // so that the caller can read data and call methods from it.
        //
        // Parameters: id: The ID of the NFT to get the reference for
        //
        // Returns: A reference to the NFT
        pub fun borrowNFTContract_NFT(id: UInt64): &NFTContract.NFT? {
            if self.ownedNFTs[id] != nil {
                let ref = (&self.ownedNFTs[id] as auth &NonFungibleToken.NFT?)!
                return ref as! &NFTContract.NFT
            } else {
                return nil
            }
        }
        pub fun borrowViewResolver(id: UInt64): &AnyResource{MetadataViews.Resolver} {
            let nft = (&self.ownedNFTs[id] as auth &NonFungibleToken.NFT?)!
            let NFTContract = nft as! &NFTContract.NFT
            return NFTContract as &AnyResource{MetadataViews.Resolver}
        }
        init() {
            self.ownedNFTs <- {}
        }
        
        
        destroy () {
            destroy self.ownedNFTs
        }
    }


    // Interface, which contains all the methods that are called by any user to mint NFT and manage brand, and template funtionality
    pub resource interface NFTMethodsCapability {
        pub fun createNewBrand(brandName: String, data: {String: String})
        pub fun updateBrandData(brandId: UInt64, data: {String: String})
        pub fun createTemplate(brandId: UInt64, maxSupply: UInt64, immutableData: {String: AnyStruct}, mutableData: {String: AnyStruct}?)
        pub fun updateTemplateMutableData(templateId: UInt64, mutableData: {String: AnyStruct})
        pub fun updateTemplateMutableAttribute(templateId: UInt64, key: String, value: AnyStruct)
        pub fun mintNFT(
            templateId: UInt64,
            account: Address,
            immutableData:{String:AnyStruct}?,
            name: String,
            description: String,
            thumbnail: String,
            royalties: [MetadataViews.Royalty])
        pub fun removeTemplateById(templateId: UInt64): Bool
    }
    
    //AdminCapability to add whiteListedAccounts
    pub resource AdminCapability {
        pub fun addwhiteListedAccount(_user: Address) {
            pre{
                NFTContract.whiteListedAccounts.contains(_user) == false: "user already exist"
            }
            NFTContract.whiteListedAccounts.append(_user)
        }

        pub fun isWhiteListedAccount(_user: Address): Bool {
            return NFTContract.whiteListedAccounts.contains(_user)
        }

        init(){}
    }

    // AdminResource, where are defining all the methods related to Brands, Template and NFTs
    pub resource AdminResource: NFTMethodsCapability {
        // a variable which stores all Brands owned by a user
        priv var ownedBrands: {UInt64: Brand}
         // a variable which stores all Templates owned by a user
        priv var ownedTemplates: {UInt64: Template}
        // a variable that store user capability to utilize methods 
        access(contract) var capability: Capability<&{NFTMethodsCapability}>?

        //method to create new Brand, only access by the verified user
        pub fun createNewBrand(brandName: String, data: {String: String}) {
            pre {
                NFTContract.whiteListedAccounts.contains(self.owner!.address): "you are not authorized for this action"
            }
            
            let newBrand = Brand(brandName: brandName, author: self.owner?.address!, data: data)
            NFTContract.allBrands[NFTContract.lastIssuedBrandId] = newBrand
            emit BrandCreated(brandId: NFTContract.lastIssuedBrandId ,brandName: brandName, author: self.owner?.address!, data: data)
            self.ownedBrands[NFTContract.lastIssuedBrandId] = newBrand 
            NFTContract.lastIssuedBrandId = NFTContract.lastIssuedBrandId + 1
        }

        //method to update the existing Brand, only author of brand can update this brand
        pub fun updateBrandData(brandId: UInt64, data: {String: String}) {
            pre{
                NFTContract.whiteListedAccounts.contains(self.owner!.address): "you are not authorized for this action"
                NFTContract.allBrands[brandId] != nil: "brand Id does not exists"
            }

            let oldBrand = NFTContract.allBrands[brandId]
            if self.owner?.address! != oldBrand!.author {
                panic("No permission to update others brand")
            }

            NFTContract.allBrands[brandId]!.update(data: data)
            emit BrandUpdated(brandId: brandId, brandName: oldBrand!.brandName, author: oldBrand!.author, data: data)
        }

        //method to create new Template, only access by the verified user
        pub fun createTemplate(brandId: UInt64, maxSupply: UInt64, immutableData: {String: AnyStruct}, mutableData:{String: AnyStruct}?) {
            pre { 
                NFTContract.whiteListedAccounts.contains(self.owner!.address): "you are not authorized for this action"
                self.ownedBrands[brandId] != nil: "Collection Id Must be valid"
                }

            let newTemplate = Template(brandId: brandId,  maxSupply: maxSupply, immutableData: immutableData, mutableData: mutableData)
            NFTContract.allTemplates[NFTContract.lastIssuedTemplateId] = newTemplate
            emit TemplateCreated(templateId: NFTContract.lastIssuedTemplateId, brandId: brandId, maxSupply: maxSupply)
            self.ownedTemplates[NFTContract.lastIssuedTemplateId] = newTemplate
            NFTContract.lastIssuedTemplateId = NFTContract.lastIssuedTemplateId + 1
        }
          //method to update the existing template's mutable data, only author of brand can update this template
        pub fun updateTemplateMutableData(templateId: UInt64, mutableData: {String: AnyStruct}) {
            pre{
                NFTContract.whiteListedAccounts.contains(self.owner!.address): "you are not authorized for this action"
                NFTContract.allTemplates[templateId] != nil: "Template Id does not exists"        
            }

            let oldTemplate = NFTContract.allTemplates[templateId]
            if self.owner?.address! != NFTContract.allBrands[oldTemplate!.brandId]!.author {
                panic("No permission to update others Template's Mutable Data")
            }

            NFTContract.allTemplates[templateId]!.updateMutableData(mutableData: mutableData)
            emit TemplateUpdated(templateId: templateId)
        }

        //method to update or add particular key-value pair in Template's mutable data, only author of brand can update this template
        pub fun updateTemplateMutableAttribute(templateId: UInt64, key: String, value: AnyStruct) {
            pre{
                NFTContract.whiteListedAccounts.contains(self.owner!.address): "you are not authorized for this action"
                NFTContract.allTemplates[templateId] != nil: "Template Id does not exists"   
            }

            let oldTemplate = NFTContract.allTemplates[templateId]
            if self.owner?.address! != NFTContract.allBrands[oldTemplate!.brandId]!.author {
                panic("No permission to update others Template's Mutable Data")
            }

            NFTContract.allTemplates[templateId]!.updateMutableAttribute(key: key, value: value)
            emit TemplateUpdated(templateId: templateId)
        }
        //method to mint NFT, only access by the verified user
        pub fun mintNFT(
            templateId: UInt64,
            account: Address,
            immutableData:{String:AnyStruct}?,
            name: String,
            description: String,
            thumbnail: String,
            royalties: [MetadataViews.Royalty]) {
            pre{
                NFTContract.whiteListedAccounts.contains(self.owner!.address): "you are not authorized for this action"
                self.ownedTemplates[templateId]!= nil: "Minter does not have specific template Id"
                NFTContract.allTemplates[templateId] != nil: "Template Id must be valid"
                }
            let receiptAccount = getAccount(account)
            let recipientCollection = receiptAccount
                .getCapability(NFTContract.CollectionPublicPath)
                .borrow<&{NFTContract.NFTContractCollectionPublic}>()
                ?? panic("Could not get receiver reference to the NFT Collection")
            let metadata : {String: AnyStruct} = {}
            var newNFT: @NFT <- create NFT(
                templateID: templateId,
                mintNumber: NFTContract.allTemplates[templateId]!.incrementIssuedSupply(),
                immutableData:immutableData,
                name: name,
                description: description,
                thumbnail: thumbnail,
                royalties: royalties,
                metadata: metadata)
            recipientCollection.deposit(token: <-newNFT)
        }

          //method to remove template by id
        pub fun removeTemplateById(templateId: UInt64): Bool {
            pre {
                NFTContract.whiteListedAccounts.contains(self.owner!.address): "you are not authorized for this action"
                templateId != nil: "invalid template id"
                NFTContract.allTemplates[templateId]!=nil: "template id does not exist"
                NFTContract.allTemplates[templateId]!.issuedSupply == 0: "could not remove template with given id"   
            }
            let mintsData =  NFTContract.allTemplates.remove(key: templateId)
            emit TemplateRemoved(templateId: templateId)
            return true
        }

        init() {
            self.ownedBrands = {}
            self.ownedTemplates = {}
            self.capability = nil
        }
    }
    
    //method to create empty Collection
    pub fun createEmptyCollection(): @NonFungibleToken.Collection {
        return <- create NFTContract.Collection()
    }

    //method to get all brands
    pub fun getAllBrands(): {UInt64: Brand} {
        return NFTContract.allBrands
    }

    //method to get brand by id
    pub fun getBrandById(brandId: UInt64): Brand {
        pre {
            NFTContract.allBrands[brandId] != nil: "brand Id does not exists"
        }
        return NFTContract.allBrands[brandId]!
    }

    //method to get all templates
    pub fun getAllTemplates(): {UInt64: Template} {
        return NFTContract.allTemplates
    }

    //method to get template by id
    pub fun getTemplateById(templateId: UInt64): Template {
        pre {
            NFTContract.allTemplates[templateId]!=nil: "Template id does not exist"
        }
        return NFTContract.allTemplates[templateId]!
    } 

    //method to get nft-data by id
    pub fun getNFTDataById(nftId: UInt64): NFTData {
        pre {
            NFTContract.allNFTs[nftId]!=nil: "nft id does not exist"
        }
        return NFTContract.allNFTs[nftId]!
    }

    //Initialize all variables with default values
    init() {
        self.lastIssuedBrandId = 1
        self.lastIssuedTemplateId = 1
        self.totalSupply = 0
        self.allBrands = {}
        self.allTemplates = {}
        self.allNFTs = {}
        self.whiteListedAccounts = [self.account.address]

        self.AdminResourceStoragePath = /storage/NFTContractAdminResource
        self.CollectionStoragePath = /storage/NFTContractCollection
        self.CollectionPublicPath = /public/NFTContractCollection
        self.AdminStorageCapability = /storage/NFTContractAdminCapability
        self.AdminCapabilityPrivate = /private/NFTContractAdminCapability
        self.NFTMethodsCapabilityPrivatePath = /private/NFTContractNFTMethodsCapability
        
        self.account.save<@AdminCapability>(<- create AdminCapability(), to: /storage/AdminStorageCapability)
        self.account.link<&AdminCapability>(self.AdminCapabilityPrivate, target: /storage/AdminStorageCapability)
        self.account.save<@AdminResource>(<- create AdminResource(), to: self.AdminResourceStoragePath)
        self.account.link<&{NFTMethodsCapability}>(self.NFTMethodsCapabilityPrivatePath, target: self.AdminResourceStoragePath)

        emit ContractInitialized()
    }
}