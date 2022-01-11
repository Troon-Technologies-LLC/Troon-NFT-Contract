import NonFungibleToken from "./NonFungibleToken.cdc"

pub contract NFTContract : NonFungibleToken {

    pub event ContractInitialized()
    pub event Withdraw(id: UInt64, from: Address?)
    pub event Deposit(id: UInt64, to: Address?)
    pub event NFTBorrowed(id:UInt64)
    pub event NFTDestroyed(id:UInt64)
    pub event NFTMinted(nftId:UInt64, templateId:UInt64, mintNumber: UInt64)
    pub event BrandCreated(brandId:UInt64, brandName:String, author:Address, data:{String:String})
    pub event BrandUpdated(brandId:UInt64, brandName:String, author:Address, data:{String:String})
    pub event SchemaCreated(schemaId:UInt64, schemaName:String, author:Address)
    pub event TemplateCreated(templateId:UInt64, brandId:UInt64, schemaId:UInt64, maxSupply:UInt64)

    pub let AdminResourceStoragePath: StoragePath
    
    pub let NFTMethodsCapabilityPrivatePath: PrivatePath

    pub let CollectionStoragePath: StoragePath
    pub let CollectionPublicPath: PublicPath

    pub var lastIssuedBrandId: UInt64
    pub var lastIssuedSchemaId: UInt64
    pub var lastIssuedTemplateId: UInt64
    pub var totalSupply:UInt64

    access(self) var allBrands: {UInt64: Brand}
    access(self) var allSchemas:{UInt64:Schema}
    access(self) var allTemplates:{UInt64:Template}
    access(self) var allNFTs : {UInt64:NFTData}

    pub enum SchemaType :UInt8{
        pub case String
        pub case Int
        pub case Fix64
        pub case Bool
        pub case Address
        pub case Array
        pub case Any        
    }

    pub struct Brand {
        pub let brandId : UInt64
        pub let brandName : String
        pub let author : Address
        pub var data : {String: String}
        
        init(brandName:String, author: Address, data:{String:String}){
            pre{
                brandName.length > 0: "Brand name is required";
            }
            let newBrandId = NFTContract.lastIssuedBrandId + 1
            self.brandId = newBrandId
            self.brandName = brandName
            self.author = author
            self.data = data
            NFTContract.lastIssuedBrandId = newBrandId
        }
        pub fun update(data:{String:String}){
            self.data = data
        }
    }

    pub struct Schema {
        pub let schemaId: UInt64
        pub let schemaName: String
        pub let author:Address
        pub let format: {String:SchemaType}

        init(schemaName : String, author:Address, format:{String:SchemaType}){
            pre{
                schemaName.length> 0 : "Could not create schema: name is required"
            }
            let newSchemaId = NFTContract.lastIssuedSchemaId + 1
            self.schemaId = newSchemaId
            self.schemaName = schemaName
            self.author = author      
            self.format = format
            NFTContract.lastIssuedSchemaId = NFTContract.lastIssuedSchemaId + 1
        }
    }

    pub struct Template {
        pub let templateId: UInt64
        pub let brandId: UInt64
        pub let schemaId: UInt64
        pub var maxSupply: UInt64
        pub var issuedSupply: UInt64
        pub var immutableData:  {String:AnyStruct}

        init(brandId:UInt64, schemaId:UInt64, maxSupply:UInt64,  immutableData: {String:AnyStruct}) {
            pre {
                NFTContract.allBrands[brandId] != nil:"Brand Id must be valid"
                NFTContract.allSchemas[schemaId] !=nil:"Schema Id must be valid"
                maxSupply > 0 : "MaxSupply must be greater than zero"
                immutableData != nil: "ImmutableData must not be nil"   
            }
            let schema = NFTContract.allSchemas[schemaId]!
            var invalidKey : String = ""
            var isValidTemplate = true

            for key in immutableData.keys {
                let value = immutableData[key]!

                if(schema.format[key] ==nil){ 
                    isValidTemplate = false
                    invalidKey = "key $".concat(key.concat(" not found"))
                    break
                }
                if schema.format[key] == NFTContract.SchemaType.String{
                    if(value as? String ==nil){ 
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    }
                }
                else if schema.format[key] == NFTContract.SchemaType.Int{
                    if(value as? Int ==nil){  
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    }
                } 
                else if schema.format[key] == NFTContract.SchemaType.Fix64 {
                    if(value as? Fix64 ==nil){  
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    } 
                }else if schema.format[key] == NFTContract.SchemaType.Bool{
                    if(value as? Bool ==nil){
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    }
                }else if schema.format[key] == NFTContract.SchemaType.Address{
                    if(value as? Address ==nil){
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    } 
                }
                else if schema.format[key] == NFTContract.SchemaType.Array{
                    if(value as? [AnyStruct] ==nil){
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    } 
                }
                else if schema.format[key] == NFTContract.SchemaType.Any{
                    if(value as? {String:AnyStruct} ==nil){
                        isValidTemplate = false
                        invalidKey = "key $".concat(key.concat(" has type mismatch"))
                        break
                    } 
                }   
            }
            assert(isValidTemplate, message: "invalid template data. Error: ".concat(invalidKey))

            let newTemplateId = NFTContract.lastIssuedTemplateId + 1
            self.templateId = newTemplateId
            self.brandId = brandId
            self.schemaId = schemaId
            self.maxSupply = maxSupply
            self.immutableData = immutableData
            self.issuedSupply = 0
            NFTContract.lastIssuedTemplateId = NFTContract.lastIssuedTemplateId + 1
        }

        pub fun incrementIssuedSupply() :UInt64{
            pre {
                self.issuedSupply  < self.maxSupply: "Template reached max supply"
            }   
            self.issuedSupply = self.issuedSupply + 1            
            return self.issuedSupply        
        }
    }

    pub struct NFTData {
        pub let templateID:UInt64
        pub let mintNumber:UInt64

        init(templateID:UInt64, mintNumber:UInt64){
            self.templateID = templateID
            self.mintNumber = mintNumber
        }

    }

    pub resource NFT: NonFungibleToken.INFT {
        pub let id: UInt64

        pub let data: NFTData

        init(templateID:UInt64, mintNumber:UInt64){
            NFTContract.totalSupply = NFTContract.totalSupply + 1
            self.id = NFTContract.totalSupply
            NFTContract.allNFTs[self.id] = NFTData(templateID:templateID,mintNumber:mintNumber)
            self.data =  NFTContract.allNFTs[self.id]!           
            emit NFTMinted(nftId:self.id, templateId:templateID, mintNumber: mintNumber)  
        }
        destroy (){
            emit NFTDestroyed(id:self.id)
        }
    }

    pub resource Collection : NonFungibleToken.Provider, NonFungibleToken.Receiver, NonFungibleToken.CollectionPublic{    
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
            emit NFTBorrowed(id:id)
            return &self.ownedNFTs[id] as &NonFungibleToken.NFT
        }

        init(){
            self.ownedNFTs <- {}    
        }
        
        destroy (){
            destroy self.ownedNFTs
        }
    }

    pub resource interface UserSpecialCapability {
        pub fun addCapability(cap: Capability<&{NFTMethodsCapability}>)
    }

    pub resource interface NFTMethodsCapability {
        pub fun createNewBrand(brandName:String, data:{String:String})
        pub fun updateBrandData(brandId:UInt64, data:{String:String})
        pub fun createSchema(schemaName : String , format:{String:SchemaType})
        pub fun createTemplate(brandId:UInt64, schemaId:UInt64, maxSupply:UInt64,  immutableData: {String:AnyStruct})
        pub fun mintNFT(templateId:UInt64, account:Address)
    }
    
    pub resource AdminResource : UserSpecialCapability, NFTMethodsCapability {

        priv var ownedBrands : {UInt64:Brand}
        priv var ownedSchemas : {UInt64:Schema}
        priv var ownedTemplates : {UInt64:Template}
        access(contract) var capability: Capability<&{NFTMethodsCapability}>?
        // this is the addCapability method that the Admin owner calls
        // to add the SpecialCapability to the AdminResource
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

        pub fun createNewBrand(brandName:String, data:{String:String}){
            pre {
                // the transaction will instantly revert if 
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
            }
            let newBrand = Brand(brandName: brandName, author: self.owner?.address!, data: data)
            NFTContract.allBrands[NFTContract.lastIssuedBrandId] = newBrand
            self.ownedBrands[NFTContract.lastIssuedBrandId] = newBrand
            emit BrandCreated(brandId:NFTContract.lastIssuedBrandId ,brandName:brandName, author: self.owner?.address!, data:data)
        }

        pub fun updateBrandData(brandId:UInt64, data:{String:String}){
            pre{
                // the transaction will instantly revert if 
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
                NFTContract.allBrands[brandId] != nil:"brand Id does not exists"   
            }
            let oldBrand = NFTContract.allBrands[brandId]
            if self.owner?.address! != oldBrand!.author {
                panic("No permission to update others brand")
            }
            NFTContract.allBrands[brandId]!.update(data:data)  
            emit BrandUpdated(brandId:brandId, brandName:oldBrand!.brandName, author:oldBrand!.author, data:data)
        }

        pub fun createSchema(schemaName: String, format: {String: SchemaType}){
            pre {
                // the transaction will instantly revert if 
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
            }
            let newSchema = Schema(schemaName: schemaName, author: self.owner?.address!, format: format)
            NFTContract.allSchemas[NFTContract.lastIssuedSchemaId] = newSchema
            self.ownedSchemas[NFTContract.lastIssuedSchemaId] = newSchema
            emit SchemaCreated(schemaId: NFTContract.lastIssuedSchemaId, schemaName: schemaName, author: self.owner?.address!)
        }

        pub fun createTemplate(brandId: UInt64, schemaId: UInt64,maxSupply: UInt64, immutableData: {String:AnyStruct}){
            pre {   
                // the transaction will instantly revert if 
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
                self.ownedBrands[brandId]!= nil :"Collection Id Must be valid"
                self.ownedSchemas[schemaId]!=nil :"Schema Id Must be valid"     
            }
            let newTemplate = Template(brandId:brandId, schemaId:schemaId, maxSupply:maxSupply,  immutableData: immutableData)
            NFTContract.allTemplates[NFTContract.lastIssuedTemplateId] = newTemplate
            self.ownedTemplates[NFTContract.lastIssuedTemplateId] = newTemplate
            emit TemplateCreated(templateId:NFTContract.lastIssuedTemplateId, brandId:brandId, schemaId:schemaId, maxSupply:maxSupply)
        }

        pub fun mintNFT(templateId:UInt64, account:Address){
            pre{
                // the transaction will instantly revert if 
                // the capability has not been added
                self.capability != nil: "I don't have the special capability :("
                self.ownedTemplates[templateId]!= nil :"Minter does not have specific template Id"     
                NFTContract.allTemplates[templateId] != nil: "Template Id must be valid"
                }
            let receiptAccount = getAccount(account)
            let recipientCollection = receiptAccount
                .getCapability(NFTContract.CollectionPublicPath)
                .borrow<&{NonFungibleToken.CollectionPublic}>()
                ?? panic("Could not get receiver reference to the NFT Collection")
            var newNFT: @NFT <- create NFT(templateID:templateId,mintNumber:NFTContract.allTemplates[templateId]!.incrementIssuedSupply())  
            recipientCollection.deposit(token: <-newNFT)
        }
        init(){
            self.ownedBrands = {}
            self.ownedSchemas = {}
            self.ownedTemplates = {}    
            self.capability = nil
        }
    }
    
    pub fun createEmptyCollection(): @NonFungibleToken.Collection {
        return <- create NFTContract.Collection()
    }

    pub fun createAdminResource(): @AdminResource{
        return <- create AdminResource()    
    }

    pub fun getAllBrands():{UInt64:Brand}{
        return NFTContract.allBrands
    }

    pub fun getBrandById(brandId:UInt64):Brand{
        pre {
            NFTContract.allBrands[brandId] != nil:"brand Id does not exists"  
        }
        return NFTContract.allBrands[brandId]!
    }

    pub fun getAllSchemas():{UInt64:Schema}{
        return NFTContract.allSchemas
    }

    pub fun getSchemaById(schemaId:UInt64): Schema{
        pre {
            NFTContract.allSchemas[schemaId]!=nil:"schema id does not exist"
        }
        return NFTContract.allSchemas[schemaId]!
    }
    
    pub fun getAllTemplates():{UInt64:Template}{
        return NFTContract.allTemplates
    }

    pub fun getTemplateById(templateId:UInt64): Template{
        pre {
            NFTContract.allTemplates[templateId]!=nil:"Template id does not exist"
        }
        return NFTContract.allTemplates[templateId]!
    } 
    pub fun getNFTDataById(nftId:UInt64):NFTData{
        return NFTContract.allNFTs[nftId]!
    }
    init(){
        self.lastIssuedBrandId = 0
        self.lastIssuedSchemaId = 0
        self.lastIssuedTemplateId = 0
        self.totalSupply = 0
        self.allBrands = {}
        self.allSchemas = {}
        self.allTemplates = {}
        self.allNFTs = {}

        self.AdminResourceStoragePath = /storage/TroonAdminResource  
        self.CollectionStoragePath = /storage/TroonCollection
        self.CollectionPublicPath = /public/TroonCollection

        self.NFTMethodsCapabilityPrivatePath = /private/NFTMethodsCapability

        emit ContractInitialized()
    }
}