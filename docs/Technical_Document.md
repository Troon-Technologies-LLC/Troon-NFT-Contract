
## Technical Summary and Code Documentation

## Instructions for creating Brand, Schema, Template and Mint Templates

A common order of creating NFT would be
 - Create Admin Account with `transaction/setupAdminAccount.cdc`.
 - Owner then create that account Admin, and gives that account ability to create own Brand, Schema and Template with `transactions/AddAdminCapability` 
 - Create new Brand with `transactions/createBrand.cdc` using Admin Account.
 - Create new Schema with `transactions/createSchema.cdc` using Admin Account.
 - Create new Template with `transactions/createTemplate.cdc` using Admin Account.
 - Create NFT Receiver with `transaction/SetupAccount` .
 - Create Mint of Templates and transfer to Address(having Setup Account enabled) with `transaction/mintNFT`
to receive NFTs.
You can also see the scripts in `transactions/scripts` to see how information
can be read from the NFTContract. 

### NFTContract Events

 - Contract Initialized ->
` pub event ContractInitialized()` 
This event is emitted when the `NFTContract` will be initialized.

- Event for Brand ->
`pub event BrandCreated(brandName:String)`
Emitted when a new Brand will be created and added to the smart Contract.

- Event for Brand Updation ->
`pub event BrandUpdated(brandId:UInt64, brandName:String, author:Address, data:{String:String})` 
Emitted when a Brand will be update

- Event for Schema ->
`pub event SchemaCreated(schemaId:UInt64, schemaName:String, author:Address)`
Emitted when a new Schema will be created

- Event for Template ->
`pub event TemplateCreated(templateId:UInt64, brandId:UInt64, schemaId:UInt64, maxSupply:UInt64)`
Emitted when a new Template will be created

-  Event for Template Mint ->
`pub event NFTMinted(nftId:UInt64, templateId:UInt64, mintNumber: UInt64`
Emitted when a Template will be Minted and save as NFT


## NFTContract Addresses

`NFTContract.cdc`: This is the main NFTContract smart contract that defines
the core functionality of the NFT.

| Network | Contract Address     |
|---------|----------------------|
| Testnet | `0x8f5c3c561b83eae3` |


## NFTContract Overview Technical

Each NFTContract represent a standard to create an NFT. We inherited NonFungibleToken to conform our standard with the existent NFT standard.
To Create an NFT, you first have to create a Brand structure which contains following fields:
- brandId: UInt64 (Id of Brand)
- brandName: String (Name of Brand)
- author: Address (Owner that created Collection)
- data:{String:String} (Metadata of Brand)
The transaction will create the brand taking input above mentioned fields. We can update metadata later using Update function(only owner can perform this action).

We also have to Create Schema Structure before creating a template using the following fields:
- schemaId: UInt64 (Id of Schema)
- schameName: String (Name of Schema)
- author: Address (Address of Owner)
- format: {String:SchemaType} 
The transaction will create the schema taking input above mentioned fields. This schema is like a database structure which is already given and if you want to create a template using that schema. You have to follow schema Structure.

We will then create Template using brandId and schemaId that we created before. Without brandId and schemaId we can't create template. We can create Template using following fields:
- templateId: UInt64 (Id of template)
- brandId: UInt64 (Foreign Id of Brand)
- schemaId: UInt64 (Foreign Id of Schema)
- maxSupply: UInt64 (maximum NFTs that could be created using that template)
- issuedSupply: UInt64 (Issued NFTs that are minted using that template)
- immutableData: {String:AnyStruct} (Immutable metadata of template)

We then have our Resource type NFT(actual asset) that represents a template owns by a user. It stores its unique Id and NFTData structure contains TemplateId and mintNumber of Template. 

The above transaction can only be performed by an Admin having an Admin resource that will give special capability to any user to create Brands, Schema and Template.

### Deployment Contract on Emulator

-  Run `flow project deploy --network emulator`
    - All contracts are deployed to the emulator.

After the contracts have been deployed, you can run the sample transactions
to interact with the contracts. The sample transactions are meant to be used
in an automated context, so they use transaction arguments and string template
fields. These make it easier for a program to use and interact with them.
If you are running these transactions manually in the Flow Playground or
vscode extension, you will need to remove the transaction arguments and
hard code the values that they are used for. 
