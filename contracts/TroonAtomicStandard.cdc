import NonFungibleToken from "./NonFungibleToken.cdc"

pub contract TroonAtomicStandard: NonFungibleToken {

    // Events
    pub event ContractInitialized()
    pub event Withdraw(id: UInt64, from: Address?)
    pub event Deposit(id: UInt64, to: Address?)
    pub event NFTBorrowed(id: UInt64)
    pub event NFTDestroyed(id: UInt64)
    pub event NFTMinted(nftId: UInt64, templateId: UInt64, mintNumber: UInt64)
    pub event BrandCreated(brandId: UInt64, brandName: String, author: Address, data:{String: String})
    pub event BrandUpdated(brandId: UInt64, brandName: String, author: Address, data:{String: String})
    pub event SchemaCreated(schemaId: UInt64, schemaName: String, author: Address)
    pub event TemplateCreated(templateId: UInt64, brandId: UInt64, schemaId: UInt64, maxSupply: UInt64)
    pub event TemplateRemoved(templateId: UInt64)

    // Paths
    pub let AdminResourceStoragePath: StoragePath
    pub let NFTMethodsCapabilityPrivatePath: PrivatePath
    pub let CollectionStoragePath: StoragePath
    pub let CollectionPublicPath: PublicPath
    pub let AdminStorageCapability: StoragePath
    pub let AdminCapabilityPrivate: PrivatePath

    // Latest brand-id
    pub var lastIssuedBrandId: UInt64

    // Latest schema-id
    pub var lastIssuedSchemaId: UInt64

    // Latest brand-id
    pub var lastIssuedTemplateId: UInt64

    // Total supply of all NFTs that are minted using this contract
    pub var totalSupply: UInt64
    
    // A dictionary that stores all Brands against it's brand-id.
    access(self) var allBrands: {UInt64: Brand}
    access(self) var allSchemas: {UInt64: Schema}
    access(self) var allTemplates: {UInt64: Template}
    access(self) var allNFTs: {UInt64: TroonAtomicNFTData}

    // Accounts ability to add capability
    access(self) var whiteListedAccounts: [Address]

    // Create Schema Support all the mentioned Types
    pub enum SchemaType: UInt8 {
        pub case String
        pub case Int
        pub case Fix64
        pub case Bool
        pub case Address
        pub case Array
        pub case Any
    }

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

            let newBrandId = TroonAtomicStandard.lastIssuedBrandId
            self.brandId = newBrandId
            self.brandName = brandName
            self.author = author
            self.data = data
        }
        pub fun update(data: {String: String}) {
            self.data = data
        }
    }

    // A structure that contain all the data related to a Schema
    pub struct Schema {
        pub let schemaId: UInt64
        pub let schemaName: String
        pub let author: Address
        access(contract) let format: {String: SchemaType}

        init(schemaName: String, author: Address, format: {String: SchemaType}){
            pre {
                schemaName.length > 0: "Could not create schema: name is required"
            }

            let newSchemaId = TroonAtomicStandard.lastIssuedSchemaId
            self.schemaId = newSchemaId
            self.schemaName = schemaName
            self.author = author
            self.format = format
        }
    }

    // A structure that contain all the data and methods related to Template
    pub struct Template {
        pub let templateId: UInt64
        pub let brandId: UInt64
        pub let schemaId: UInt64
        pub var maxSupply: UInt64
        pub var issuedSupply: UInt64
        pub var immutableData: {String: AnyStruct}

        init(brandId: UInt64, schemaId: UInt64, maxSupply: UInt64, immutableData: {String: AnyStruct}) {
            pre {
                TroonAtomicStandard.allBrands[brandId] != nil:"Brand Id must be valid"
                TroonAtomicStandard.allSchemas[schemaId] != nil:"Schema Id must be valid"
                maxSupply > 0 : "MaxSupply must be greater than zero"
                immutableData != nil: "ImmutableData must not be nil"
            }

            self.templateId = TroonAtomicStandard.lastIssuedTemplateId
            self.brandId = brandId
            self.schemaId = schemaId
            self.maxSupply = maxSupply
            self.immutableData = immutableData
            self.issuedSupply = 0
            // Before creating template, we need to check template data, if it is valid against given schema or not
            let schema = TroonAtomicStandard.allSchemas[schemaId]!
            var invalidKey: String = ""
            var isValidTemplate = true

            for key in immutableData.keys {
                let value = immutableData[key]!
                if(schema.format[key] == nil) {
                    isValidTemplate = false
                    invalidKey = "key $".concat(key.concat(" not found"))
                    break
                }
                if schema.format[key] == TroonAtomicStandard.SchemaType.String {
                    if(value as? String == nil) {
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    }
                }
                else if schema.format[key] == TroonAtomicStandard.SchemaType.Int {
                    if(value as? Int == nil) {
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    }
                } 
                else if schema.format[key] == TroonAtomicStandard.SchemaType.Fix64 {
                    if(value as? Fix64 == nil) {
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    }
                }else if schema.format[key] == TroonAtomicStandard.SchemaType.Bool {
                    if(value as? Bool == nil) {
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    }
                }else if schema.format[key] == TroonAtomicStandard.SchemaType.Address {
                    if(value as? Address == nil) {
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    }
                }
                else if schema.format[key] == TroonAtomicStandard.SchemaType.Array {
                    if(value as? [AnyStruct] == nil) {
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    }
                }
                else if schema.format[key] == TroonAtomicStandard.SchemaType.Any {
                    if(value as? {String:AnyStruct} ==nil) {
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    }
                }
            }
            assert(isValidTemplate, message: "invalid template data. Error: ".concat(invalidKey))

            
        }

        // a method to increment issued supply for template
        access(contract) fun incrementIssuedSupply(): UInt64 {
            pre {
                self.issuedSupply < self.maxSupply: "Template reached max supply"
            }   

            self.issuedSupply = self.issuedSupply + 1
            return self.issuedSupply
        }
    }

    // A structure that link template and mint-no of NFT
    pub struct TroonAtomicNFTData {
        pub let templateID: UInt64
        pub let mintNumber: UInt64

        init(templateID: UInt64, mintNumber: UInt64) {
            self.templateID = templateID
            self.mintNumber = mintNumber
        }
    }

    // The resource that represents the Troon NFTs
    // 
    pub resource NFT: NonFungibleToken.INFT {
        pub let id: UInt64
        access(contract) let data: TroonAtomicNFTData

        init(templateID: UInt64, mintNumber: UInt64) {
            TroonAtomicStandard.totalSupply = TroonAtomicStandard.totalSupply + 1
            self.id = TroonAtomicStandard.totalSupply
            TroonAtomicStandard.allNFTs[self.id] = TroonAtomicNFTData(templateID: templateID, mintNumber: mintNumber)
            self.data = TroonAtomicStandard.allNFTs[self.id]!
            emit NFTMinted(nftId: self.id, templateId: templateID, mintNumber: mintNumber)
        }
        destroy(){
            emit NFTDestroyed(id: self.id)
        }
    }

    pub resource interface TroonAtomicStandardCollectionPublic {
        pub fun deposit(token: @NonFungibleToken.NFT)
        pub fun getIDs(): [UInt64]
        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT
        pub fun borrowNFTTroonAtomicStandard(id: UInt64): &TroonAtomicStandard.NFT? {
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
    pub resource Collection: TroonAtomicStandardCollectionPublic, NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.CollectionPublic {
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
            let token <- token as! @TroonAtomicStandard.NFT
            let id = token.id
            let oldToken <- self.ownedNFTs[id] <- token
            if self.owner?.address != nil {
                emit Deposit(id: id, to: self.owner?.address)
            }
            destroy oldToken
        }

        pub fun borrowNFT(id: UInt64): &NonFungibleToken.NFT {
            emit NFTBorrowed(id:id)
            return &self.ownedNFTs[id] as &NonFungibleToken.NFT
        }

        // borrowNFTTriQuetaContract returns a borrowed reference to a TriQuetaNFT
        // so that the caller can read data and call methods from it.
        //
        // Parameters: id: The ID of the NFT to get the reference for
        //
        // Returns: A reference to the NFT
        pub fun borrowNFTTroonAtomicStandard(id: UInt64): &TroonAtomicStandard.NFT? {
            if self.ownedNFTs[id] != nil {
                let ref = &self.ownedNFTs[id] as auth &NonFungibleToken.NFT
                return ref as! &TroonAtomicStandard.NFT
            } else {
                return nil
            }
        }

        init() {
            self.ownedNFTs <- {}
        }
        
        destroy () {
            destroy self.ownedNFTs
        }
    }

    // Special Capability, that is needed by user to utilize our contract. Only verified user can get this capability so it will add a KYC layer in our white-lable-solution
    pub resource interface UserSpecialCapability {
        pub fun addCapability(cap: Capability<&{NFTMethodsCapability}>)
    }

    // Interface, which contains all the methods that are called by any user to mint NFT and manage brand, schema and template funtionality
    pub resource interface NFTMethodsCapability {
        pub fun createNewBrand(brandName: String, data: {String: String})
        pub fun updateBrandData(brandId: UInt64, data: {String: String})
        pub fun createSchema(schemaName: String, format: {String: SchemaType})
        pub fun createTemplate(brandId: UInt64, schemaId: UInt64, maxSupply: UInt64, immutableData: {String: AnyStruct})
        pub fun mintNFT(templateId: UInt64, account: Address)
    }
    
    //AdminCapability to add whiteListedAccounts
     pub resource AdminCapability{
        
        pub fun addwhiteListedAccount(_user: Address) {
            pre{
                TroonAtomicStandard.whiteListedAccounts.contains(_user) == false: "user already exist"
            }
            TroonAtomicStandard.whiteListedAccounts.append(_user)
        }

        pub fun isWhiteListedAccount(_user: Address): Bool {
            return TroonAtomicStandard.whiteListedAccounts.contains(_user)
        }

        init(){}
    }

    // AdminResource, where are defining all the methods related to Brands, Schema, Template and NFTs
    pub resource AdminResource: UserSpecialCapability, NFTMethodsCapability {
        // a variable which stores all Brands owned by a user
        priv var ownedBrands: {UInt64: Brand}
        // a variable which stores all Schema owned by a user
        priv var ownedSchemas: {UInt64: Schema}
        // a variable which stores all Templates owned by a user
        priv var ownedTemplates: {UInt64: Template}
        // a variable that store user capability to utilize methods 
        access(contract) var capability: Capability<&{NFTMethodsCapability}>?
        // method which provide capability to user to utilize methods
        pub fun addCapability(cap: Capability<&{NFTMethodsCapability}>) {
            pre {
                // we make sure the SpecialCapability is
                // valid before executing the method
                cap.borrow() != nil: "could not borrow a reference to the SpecialCapability"
                self.capability == nil: "resource already has the SpecialCapability"
                TroonAtomicStandard.whiteListedAccounts.contains(self.owner!.address): "you are not authorized for this action"
            }
            // add the SpecialCapability
            self.capability = cap
        }

        //method to create new Brand, only access by the verified user
        pub fun createNewBrand(brandName: String, data: {String: String}) {
            pre {
                // the transaction will instantly revert if
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
                TroonAtomicStandard.whiteListedAccounts.contains(self.owner!.address): "you are not authorized for this action"
            }

            let newBrand = Brand(brandName: brandName, author: self.owner?.address!, data: data)
            TroonAtomicStandard.allBrands[TroonAtomicStandard.lastIssuedBrandId] = newBrand
            emit BrandCreated(brandId: TroonAtomicStandard.lastIssuedBrandId ,brandName: brandName, author: self.owner?.address!, data: data)
            self.ownedBrands[TroonAtomicStandard.lastIssuedBrandId] = newBrand 
            TroonAtomicStandard.lastIssuedBrandId = TroonAtomicStandard.lastIssuedBrandId + 1
        }

        //method to update the existing Brand, only author of brand can update this brand
        pub fun updateBrandData(brandId: UInt64, data: {String: String}) {
            pre{
                // the transaction will instantly revert if
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
                TroonAtomicStandard.whiteListedAccounts.contains(self.owner!.address): "you are not authorized for this action"
                TroonAtomicStandard.allBrands[brandId] != nil: "brand Id does not exists"
            }

            let oldBrand = TroonAtomicStandard.allBrands[brandId]
            if self.owner?.address! != oldBrand!.author {
                panic("No permission to update others brand")
            }

            TroonAtomicStandard.allBrands[brandId]!.update(data: data)
            emit BrandUpdated(brandId: brandId, brandName: oldBrand!.brandName, author: oldBrand!.author, data: data)
        }

        //method to create new Schema, only access by the verified user
        pub fun createSchema(schemaName: String, format: {String: SchemaType}) {
            pre {
                // the transaction will instantly revert if 
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
                TroonAtomicStandard.whiteListedAccounts.contains(self.owner!.address): "you are not authorized for this action"
            }

            let newSchema = Schema(schemaName: schemaName, author: self.owner?.address!, format: format)
            TroonAtomicStandard.allSchemas[TroonAtomicStandard.lastIssuedSchemaId] = newSchema
            emit SchemaCreated(schemaId: TroonAtomicStandard.lastIssuedSchemaId, schemaName: schemaName, author: self.owner?.address!)
            self.ownedSchemas[TroonAtomicStandard.lastIssuedSchemaId] = newSchema
            TroonAtomicStandard.lastIssuedSchemaId = TroonAtomicStandard.lastIssuedSchemaId + 1
            
        }

        //method to create new Template, only access by the verified user
        pub fun createTemplate(brandId: UInt64, schemaId: UInt64, maxSupply: UInt64, immutableData: {String: AnyStruct}) {
            pre { 
                // the transaction will instantly revert if 
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
                TroonAtomicStandard.whiteListedAccounts.contains(self.owner!.address): "you are not authorized for this action"
                self.ownedBrands[brandId] != nil: "Collection Id Must be valid"
                self.ownedSchemas[schemaId] != nil: "Schema Id Must be valid"
            }

            let newTemplate = Template(brandId: brandId, schemaId: schemaId, maxSupply: maxSupply, immutableData: immutableData)
            TroonAtomicStandard.allTemplates[TroonAtomicStandard.lastIssuedTemplateId] = newTemplate
            emit TemplateCreated(templateId: TroonAtomicStandard.lastIssuedTemplateId, brandId: brandId, schemaId: schemaId, maxSupply: maxSupply)
            self.ownedTemplates[TroonAtomicStandard.lastIssuedTemplateId] = newTemplate
            TroonAtomicStandard.lastIssuedTemplateId = TroonAtomicStandard.lastIssuedTemplateId + 1
        }

         //method to remove template by id
        pub fun removeTemplateById(templateId: UInt64): Bool {
            pre {
                templateId != nil: "invalid template id"
                TroonAtomicStandard.allTemplates[templateId]!=nil: "template id does not exist"
                TroonAtomicStandard.allTemplates[templateId]!.issuedSupply == 0: "could not remove template with given id"   
            }
            let mintsData =  TroonAtomicStandard.allTemplates.remove(key: templateId)
            emit TemplateRemoved(templateId: templateId)
            return true
        }


        //method to mint NFT, only access by the verified user
        pub fun mintNFT(templateId: UInt64, account: Address) {
            pre{
                // the transaction will instantly revert if 
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
                TroonAtomicStandard.whiteListedAccounts.contains(self.owner!.address): "you are not authorized for this action"
                self.ownedTemplates[templateId]!= nil: "Minter does not have specific template Id"
                TroonAtomicStandard.allTemplates[templateId] != nil: "Template Id must be valid"
                }
            let receiptAccount = getAccount(account)
            let recipientCollection = receiptAccount
                .getCapability(TroonAtomicStandard.CollectionPublicPath)
                .borrow<&{NonFungibleToken.CollectionPublic}>()
                ?? panic("Could not get receiver reference to the NFT Collection")
            var newNFT: @NFT <- create NFT(templateID: templateId, mintNumber: TroonAtomicStandard.allTemplates[templateId]!.incrementIssuedSupply())
            recipientCollection.deposit(token: <-newNFT)
        }

        init() {
            self.ownedBrands = {}
            self.ownedSchemas = {}
            self.ownedTemplates = {}
            self.capability = nil
        }
    }
    
    //method to create empty Collection
    pub fun createEmptyCollection(): @NonFungibleToken.Collection {
        return <- create TroonAtomicStandard.Collection()
    }

    //method to create Admin Resources
    pub fun createAdminResource(): @AdminResource {
        return <- create AdminResource()
    }

    //method to get all brands
    pub fun getAllBrands(): {UInt64: Brand} {
        return TroonAtomicStandard.allBrands
    }

    //method to get brand by id
    pub fun getBrandById(brandId: UInt64): Brand {
        pre {
            TroonAtomicStandard.allBrands[brandId] != nil: "brand Id does not exists"
        }
        return TroonAtomicStandard.allBrands[brandId]!
    }

    //method to get all schema
    pub fun getAllSchemas(): {UInt64: Schema} {
        return TroonAtomicStandard.allSchemas
    }

    //method to get schema by id
    pub fun getSchemaById(schemaId: UInt64): Schema {
        pre {
            TroonAtomicStandard.allSchemas[schemaId] != nil: "schema id does not exist"
        }
        return TroonAtomicStandard.allSchemas[schemaId]!
    }

    //method to get all templates
    pub fun getAllTemplates(): {UInt64: Template} {
        return TroonAtomicStandard.allTemplates
    }

    //method to get template by id
    pub fun getTemplateById(templateId: UInt64): Template {
        pre {
            TroonAtomicStandard.allTemplates[templateId]!=nil: "Template id does not exist"
        }
        return TroonAtomicStandard.allTemplates[templateId]!
    } 

    //method to get nft-data by id
    pub fun getTroonAtomicNFTDataById(nftId: UInt64): TroonAtomicNFTData {
        pre {
            TroonAtomicStandard.allNFTs[nftId]!=nil:"nft id does not exist"
        }
        return TroonAtomicStandard.allNFTs[nftId]!
    }

    //Initialize all variables with default values
    init(){
        self.lastIssuedBrandId = 1
        self.lastIssuedSchemaId = 1
        self.lastIssuedTemplateId = 1
        self.totalSupply = 0
        self.allBrands = {}
        self.allSchemas = {}
        self.allTemplates = {}
        self.allNFTs = {}
        self.whiteListedAccounts = [self.account.address]

        self.AdminResourceStoragePath = /storage/TroonAdminResourcev01
        self.CollectionStoragePath = /storage/TroonCollectionv01
        self.CollectionPublicPath = /public/TroonCollectionv01
        self.AdminStorageCapability = /storage/AdminCapability
        self.AdminCapabilityPrivate = /private/AdminCapability
        self.NFTMethodsCapabilityPrivatePath = /private/NFTMethodsCapabilityv01
        
        self.account.save<@AdminCapability>(<- create AdminCapability(), to: /storage/AdminStorageCapability)
        self.account.link<&AdminCapability>(self.AdminCapabilityPrivate, target: /storage/AdminStorageCapability)
        self.account.save<@AdminResource>(<- create AdminResource(), to: self.AdminResourceStoragePath)
        self.account.link<&{NFTMethodsCapability}>(self.NFTMethodsCapabilityPrivatePath, target: self.AdminResourceStoragePath)

        emit ContractInitialized()
    }
}
