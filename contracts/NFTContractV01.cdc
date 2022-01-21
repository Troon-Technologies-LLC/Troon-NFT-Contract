import NonFungibleToken from "./NonFungibleToken.cdc"

pub contract NFTContractV01 : NonFungibleToken {

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

    // Paths
    pub let AdminResourceStoragePath: StoragePath
    pub let NFTMethodsCapabilityPrivatePath: PrivatePath
    pub let CollectionStoragePath: StoragePath
    pub let CollectionPublicPath: PublicPath

    // Latest brand-id
    pub var lastIssuedBrandId: UInt64

    // Latest schema-id
    pub var lastIssuedSchemaId: UInt64

    // Latest brand-id
    pub var lastIssuedTemplateId: UInt64

    // Total supply of all NFTs that are minted using this contract
    pub var totalSupply:UInt64
    
    // A dictionary that stores all Brands against it's brand-id.
    access(self) var allBrands: {UInt64: Brand}
    access(self) var allSchemas: {UInt64:Schema}
    access(self) var allTemplates: {UInt64:Template}
    access(self) var allNFTs: {UInt64:NFTData}

    // Create Schema Support all the mentioned Types
    pub enum SchemaType: UInt8{
        pub case String
        pub case Int
        pub case Fix64
        pub case Bool
        pub case Address
        pub case Array
        pub case Any        
    }

    // A strcuture that contain all the data related to a Brand
    pub struct Brand {
        pub let brandId: UInt64
        pub let brandName: String
        pub let author: Address
        access(contract) var data: {String: String}
        
        init(brandName: String, author: Address, data: {String:String}) {
            pre {
                brandName.length > 0: "Brand name is required";
            }

            let newBrandId = NFTContractV01.lastIssuedBrandId
            self.brandId = newBrandId
            self.brandName = brandName
            self.author = author
            self.data = data
        }
        pub fun update(data: {String:String}){
            self.data = data
        }
    }

    // A strcuture that contain all the data related to a Schema
    pub struct Schema {
        pub let schemaId: UInt64
        pub let schemaName: String
        pub let author:Address
        access(contract) let format: {String:SchemaType}

        init(schemaName: String, author: Address, format: {String:SchemaType}){
            pre {
                schemaName.length>0: "Could not create schema: name is required"
            }

            let newSchemaId = NFTContractV01.lastIssuedSchemaId
            self.schemaId = newSchemaId
            self.schemaName = schemaName
            self.author = author      
            self.format = format
        }
    }

    // A strcuture that contain all the data and methods related to Template
    pub struct Template {
        pub let templateId: UInt64
        pub let brandId: UInt64
        pub let schemaId: UInt64
        pub var maxSupply: UInt64
        pub var issuedSupply: UInt64
        pub var immutableData:  {String:AnyStruct}

        init(brandId: UInt64, schemaId: UInt64, maxSupply: UInt64, immutableData: {String:AnyStruct}) {
            pre {
                NFTContractV01.allBrands[brandId] != nil:"Brand Id must be valid"
                NFTContractV01.allSchemas[schemaId] !=nil:"Schema Id must be valid"
                maxSupply > 0 : "MaxSupply must be greater than zero"
                immutableData != nil: "ImmutableData must not be nil"   
            }

            // Before creating template, we need to check template data, if it is valid against given schema or not
            let schema = NFTContractV01.allSchemas[schemaId]!
            var invalidKey : String = ""
            var isValidTemplate = true

            for key in immutableData.keys {
                let value = immutableData[key]!
                if(schema.format[key] ==nil){ 
                    isValidTemplate = false
                    invalidKey = "key $".concat(key.concat(" not found"))
                    break
                }
                if schema.format[key] == NFTContractV01.SchemaType.String{
                    if(value as? String ==nil){ 
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    }
                }
                else if schema.format[key] == NFTContractV01.SchemaType.Int{
                    if(value as? Int ==nil){  
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    }
                } 
                else if schema.format[key] == NFTContractV01.SchemaType.Fix64 {
                    if(value as? Fix64 ==nil){  
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    } 
                }else if schema.format[key] == NFTContractV01.SchemaType.Bool{
                    if(value as? Bool ==nil){
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    }
                }else if schema.format[key] == NFTContractV01.SchemaType.Address{
                    if(value as? Address ==nil){
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    } 
                }
                else if schema.format[key] == NFTContractV01.SchemaType.Array{
                    if(value as? [AnyStruct] ==nil){
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    } 
                }
                else if schema.format[key] == NFTContractV01.SchemaType.Any{
                    if(value as? {String:AnyStruct} ==nil){
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    } 
                }   
            }
            assert(isValidTemplate, message: "invalid template data. Error: ".concat(invalidKey))

            self.templateId = NFTContractV01.lastIssuedTemplateId
            self.brandId = brandId
            self.schemaId = schemaId
            self.maxSupply = maxSupply
            self.immutableData = immutableData
            self.issuedSupply = 0
        }

        // a method to increment issued supply for template
        access(contract) fun incrementIssuedSupply() : UInt64 {
            pre {
                self.issuedSupply  < self.maxSupply: "Template reached max supply"
            }   

            self.issuedSupply = self.issuedSupply + 1            
            return self.issuedSupply        
        }
        
    }

    // A structure that link template and mint-no of NFT
    pub struct NFTData {
        pub let templateID:UInt64
        pub let mintNumber:UInt64

        init(templateID:UInt64, mintNumber:UInt64){
            self.templateID = templateID
            self.mintNumber = mintNumber
        }
    }

    // The resource that represents the Troon NFTs
    // 
    pub resource NFT: NonFungibleToken.INFT {
        pub let id: UInt64
        access(contract) let data: NFTData

        init(templateID:UInt64, mintNumber:UInt64){
            NFTContractV01.totalSupply = NFTContractV01.totalSupply + 1
            self.id = NFTContractV01.totalSupply
            NFTContractV01.allNFTs[self.id] = NFTData(templateID:templateID,mintNumber:mintNumber)
            self.data =  NFTContractV01.allNFTs[self.id]!           
            emit NFTMinted(nftId:self.id, templateId:templateID, mintNumber: mintNumber)  
        }
        destroy (){
            emit NFTDestroyed(id:self.id)
        }
    }

    // Collection is a resource that every user who owns NFTs 
    // will store in their account to manage their NFTS
    //
    pub resource Collection: NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.CollectionPublic{    
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
            let token <- token as! @NFTContractV01.NFT
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
        pub fun createNewBrand(brandName:String, data:{String:String})
        pub fun updateBrandData(brandId:UInt64, data:{String:String})
        pub fun createSchema(schemaName : String , format:{String:SchemaType})
        pub fun createTemplate(brandId:UInt64, schemaId:UInt64, maxSupply:UInt64,  immutableData: {String:AnyStruct})
        pub fun mintNFT(templateId:UInt64, account:Address)
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
            }
            // add the SpecialCapability
            self.capability = cap
        }

        //method to create new Brand, only access by the verified user
        pub fun createNewBrand(brandName: String, data: {String: String}){
            pre {
                // the transaction will instantly revert if 
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
            }

            let newBrand = Brand(brandName: brandName, author: self.owner?.address!, data: data)
            NFTContractV01.allBrands[NFTContractV01.lastIssuedBrandId] = newBrand
            emit BrandCreated(brandId:NFTContractV01.lastIssuedBrandId ,brandName:brandName, author: self.owner?.address!, data:data)
            self.ownedBrands[NFTContractV01.lastIssuedBrandId] = newBrand   
            NFTContractV01.lastIssuedBrandId = NFTContractV01.lastIssuedBrandId + 1
        }

        //method to update the existing Brand, only author of brand can update this brand
        pub fun updateBrandData(brandId: UInt64, data: {String: String}){
            pre{
                // the transaction will instantly revert if 
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
                NFTContractV01.allBrands[brandId] != nil:"brand Id does not exists"   
            }
            let oldBrand = NFTContractV01.allBrands[brandId]
            if self.owner?.address! != oldBrand!.author {
                panic("No permission to update others brand")
            }
            NFTContractV01.allBrands[brandId]!.update(data:data)  
            emit BrandUpdated(brandId:brandId, brandName:oldBrand!.brandName, author:oldBrand!.author, data:data)
        }

        //method to create new Schema, only access by the verified user
        pub fun createSchema(schemaName: String, format: {String: SchemaType}) {
            pre {
                // the transaction will instantly revert if 
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
            }

            let newSchema = Schema(schemaName: schemaName, author: self.owner?.address!, format: format)
            NFTContractV01.allSchemas[NFTContractV01.lastIssuedSchemaId] = newSchema
            emit SchemaCreated(schemaId: NFTContractV01.lastIssuedSchemaId, schemaName: schemaName, author: self.owner?.address!)
            self.ownedSchemas[NFTContractV01.lastIssuedSchemaId] = newSchema
            NFTContractV01.lastIssuedSchemaId = NFTContractV01.lastIssuedSchemaId + 1
            
        }

        //method to create new Template, only access by the verified user
        pub fun createTemplate(brandId: UInt64, schemaId: UInt64, maxSupply: UInt64, immutableData: {String:AnyStruct}){
            pre {   
                // the transaction will instantly revert if 
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
                self.ownedBrands[brandId]!= nil :"Collection Id Must be valid"
                self.ownedSchemas[schemaId]!=nil :"Schema Id Must be valid"     
            }
            let newTemplate = Template(brandId: brandId, schemaId: schemaId, maxSupply: maxSupply,  immutableData: immutableData)
            NFTContractV01.allTemplates[NFTContractV01.lastIssuedTemplateId] = newTemplate
            emit TemplateCreated(templateId: NFTContractV01.lastIssuedTemplateId, brandId: brandId, schemaId: schemaId, maxSupply: maxSupply)
            self.ownedTemplates[NFTContractV01.lastIssuedTemplateId] = newTemplate
            NFTContractV01.lastIssuedTemplateId = NFTContractV01.lastIssuedTemplateId + 1
        }

        //method to mint NFT, only access by the verified user
        pub fun mintNFT(templateId: UInt64, account: Address){
            pre{
                // the transaction will instantly revert if 
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
                self.ownedTemplates[templateId]!= nil :"Minter does not have specific template Id"     
                NFTContractV01.allTemplates[templateId] != nil: "Template Id must be valid"
                }
            let receiptAccount = getAccount(account)
            let recipientCollection = receiptAccount
                .getCapability(NFTContractV01.CollectionPublicPath)
                .borrow<&{NonFungibleToken.CollectionPublic}>()
                ?? panic("Could not get receiver reference to the NFT Collection")
            var newNFT: @NFT <- create NFT(templateID:templateId,mintNumber:NFTContractV01.allTemplates[templateId]!.incrementIssuedSupply())  
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
        return <- create NFTContractV01.Collection()
    }

    //method to create Admin Resources
    pub fun createAdminResource(): @AdminResource{
        return <- create AdminResource()    
    }

    //method to get all brands
    pub fun getAllBrands(): {UInt64: Brand}{
        return NFTContractV01.allBrands
    }

    //method to get brand by id
    pub fun getBrandById(brandId: UInt64): Brand{
        pre {
            NFTContractV01.allBrands[brandId] != nil:"brand Id does not exists"  
        }
        return NFTContractV01.allBrands[brandId]!
    }

    //method to get all schema
    pub fun getAllSchemas(): {UInt64: Schema}{
        return NFTContractV01.allSchemas
    }

    //method to get schema by id
    pub fun getSchemaById(schemaId: UInt64): Schema {
        pre {
            NFTContractV01.allSchemas[schemaId]!=nil:"schema id does not exist"
        }
        return NFTContractV01.allSchemas[schemaId]!
    }

    //method to get all templates
    pub fun getAllTemplates(): {UInt64: Template}{
        return NFTContractV01.allTemplates
    }

    //method to get template by id
    pub fun getTemplateById(templateId:UInt64): Template{
        pre {
            NFTContractV01.allTemplates[templateId]!=nil:"Template id does not exist"
        }
        return NFTContractV01.allTemplates[templateId]!
    } 

    //method to get nft-data by id
    pub fun getNFTDataById(nftId: UInt64): NFTData{
        pre {
            NFTContractV01.allNFTs[nftId]!=nil:"nft id does not exist"
        }
        return NFTContractV01.allNFTs[nftId]!
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

        self.AdminResourceStoragePath = /storage/TroonAdminResourcev01  
        self.CollectionStoragePath = /storage/TroonCollectionv01
        self.CollectionPublicPath = /public/TroonCollectionv01

        self.NFTMethodsCapabilityPrivatePath = /private/NFTMethodsCapabilityv01
        
        self.account.save<@AdminResource>(<- create AdminResource(), to: self.AdminResourceStoragePath)
        self.account.link<&{NFTMethodsCapability}>(self.NFTMethodsCapabilityPrivatePath, target: self.AdminResourceStoragePath)

        emit ContractInitialized()
    }
}