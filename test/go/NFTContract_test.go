package test

import (
	"log"
	"regexp"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/assert"
)

var BrandName string = "Nowwhere"
var BrandMetadataKey string = "brandName"
var BrandMetadataValue string = "Mercedes-Benz"
var UpdateBrandMetadataValue string = "Mercedes-Benz_GLA"
var SchemaName string = "artist"
var TemplateName string = "twin-turbo"
var TemplateMetdataKey string = "engine"
var shouldFail = true
var shouldNotFail = false

const (
	zero  = 0
	one   = 1
	two   = 2
	three = 3
)

// Description: Check Contract Deployment and field initialization total NFT supply to be zero for Verification
// Input: Empty
// Expected Output: NFT Supply initialized with zero
// Test-case-type: Positive
func Test_NFTContractTestDeployment(test *testing.T) {
	emulator := newEmulator()

	nonfungibleAddr, ownerAddr, _, _ := NFTContractDeployContracts(emulator, test)

	test.Run("Should have initialized Supply field zero correctly", func(test *testing.T) {
		supply := executeScriptAndCheck(test, emulator, NowwhereGenerateGetSupplyScript(nonfungibleAddr, ownerAddr), nil)
		var supplyOnitial uint64 = uint64(zero)
		assert.EqualValues(test, CadenceUInt64(supplyOnitial), supply)
	})
}

// Description: Deploy Contract and Check for Brand initialization to be zero for verification
// Input: Empty
// Expected Output: Fields are created successfully but empty
// Test-case-type: Positive

func Test_GetBrand_EmptyOnDeployment(test *testing.T) {
	emulator := newEmulator()

	nonfungibleAddr, ownerAddr, _, _ := NFTContractDeployContracts(emulator, test)

	test.Run("Should have Brand initialized field correctly", func(test *testing.T) {
		supply := executeScriptAndCheck(
			test,
			emulator,
			NFTContractGenerateGetBrandCountScript(nonfungibleAddr, ownerAddr),
			nil)

		assert.EqualValues(test, CadenceInt(zero), supply)
	})
}

// Description: Create Brand with the name as input Without making an Account Admin
// Input: BrandName
// Expected Output: Brand not Created
// Test-case-type: Negative
func Test_CreateBrand_WithAnyAddress(test *testing.T) {
	emulator := newEmulator()

	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	// Pre Condition: Check Brand not initialized
	test.Run("Should have Zero Brand count", func(test *testing.T) {
		brandCount := executeScriptAndCheck(test, emulator, NFTContractGenerateGetBrandCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), brandCount)
	})

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brandMetadata := cadence.NewDictionary(metadata)
	NFTContractCreateBrandTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldFail, // This transaction should  fail
		BrandName,
		brandMetadata,
	)

	// Post Condition: Check Brand initialized properly
	test.Run("Should have initialized Brands correctly", func(test *testing.T) {
		brandCount := executeScriptAndCheck(test, emulator, NFTContractGenerateGetBrandCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), brandCount)
	})
}

// Description: Create Brand with the name as input with Admin Account
// Input: BrandName
// Expected Output: Brand Created
// Test-case-type: Positive
func Test_CreateBrand_Success(test *testing.T) {
	emulator := newEmulator()

	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	// Pre Condition: Check brand not initialized
	test.Run("Should have Zero Brand count", func(test *testing.T) {
		brandCount := executeScriptAndCheck(test, emulator, NFTContractGenerateGetBrandCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), brandCount)
	})

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brand := cadence.NewDictionary(metadata)
	NFTContractCreateBrandTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail,
		BrandName,
		brand,
	)

	// Post Condition: Check Brand initialized properly
	test.Run("Should have initialized Brands correctly", func(test *testing.T) {
		brandCount := executeScriptAndCheck(test, emulator, NFTContractGenerateGetBrandCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(one), brandCount)
	})
}

// Description: Update Brand Data with ID as input and updated metdata
// Input: brandID
// Expected Output: Brand Updated
// Test-case-type: Positive
func Test_UpdateBrand_Success(test *testing.T) {
	emulator := newEmulator()

	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	// Pre Condition: Check brand not initialized
	test.Run("Should have Zero brand count", func(test *testing.T) {
		brandCount := executeScriptAndCheck(test, emulator, NFTContractGenerateGetBrandCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), brandCount)
	})

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brand := cadence.NewDictionary(metadata)
	NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, BrandName, brand)

	// Post Condition: Check brand initialized properly
	test.Run("Should have initialized Brands correctly", func(test *testing.T) {
		brandCount := executeScriptAndCheck(test, emulator, NFTContractGenerateGetBrandCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(one), brandCount)
	})
	//brandID := executeScriptAndCheck(test, emulator, NFTContractGenerateGetBrandIDsScript(nonfungibleAddr, ownerAddr), nil)
	//assert.EqualValues(test, CadenceInt(one), brandID)
	//	assert.EqualValues(test, CadenceInt(one), brandID)
	NFTContractUpdateBrandTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail,
		one,                      // brand ID
		UpdateBrandMetadataValue, //UpdateModel
	)
	// Make a Regex to say we only want letters and numbers because Script returns JSON String having some characters other than alphabets
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	// Post Condition: Check Brand metadata updated initialized properly
	test.Run("Should update Brands correctly", func(test *testing.T) {
		ReturnedbrandModel := executeScriptAndCheck(test, emulator, NFTContractGenerateGetBrandNameScript(nonfungibleAddr, ownerAddr), [][]byte{jsoncdc.MustEncode(cadence.NewUInt64(one))})
		assert.EqualValues(test, reg.ReplaceAllString(UpdateBrandMetadataValue, ""), reg.ReplaceAllString(ReturnedbrandModel.String(), ""))
	})

}

// Description: Update Brand with ID as input and update metdata
// Input: brandID
// Expected Output: Brand Not Update
// Test-case-type: Negative
func Test_UpdateBrand_WithInvalidId(test *testing.T) {
	emulator := newEmulator()

	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	// Pre Condition: Check brand not initialized
	test.Run("Should have Zero brand count", func(test *testing.T) {
		brandCount := executeScriptAndCheck(test, emulator, NFTContractGenerateGetBrandCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), brandCount)
	})

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	brandNameField := CadenceString(BrandMetadataKey)
	brandName := CadenceString(BrandMetadataValue)

	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brand := cadence.NewDictionary(metadata)
	NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, BrandName, brand)

	// Post Condition: Check brand initialized properly
	test.Run("Should have initialized Brands correctly", func(test *testing.T) {
		brandCount := executeScriptAndCheck(test, emulator, NFTContractGenerateGetBrandCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(one), brandCount)
	})

	//brandID := executeScriptAndCheck(test, emulator, NFTContractGenerateGetBrandIDsScript(nonfungibleAddr, ownerAddr), nil)
	//assert.EqualValues(test, CadenceInt(one), brandID)
	//	assert.EqualValues(test, CadenceInt(one), brandID)
	NFTContractUpdateBrandTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldFail,
		two,                      // wrong brand ID
		UpdateBrandMetadataValue, //UpdateModel
	)
	// Make a Regex to say we only want letters and numbers because Script returns JSON String having some characters other than alphabets
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	// Post Condition: Check Brand name not updated properly
	test.Run("Should contains metadata Brands correctly", func(test *testing.T) {
		ReturnedbrandModel := executeScriptAndCheck(test, emulator, NFTContractGenerateGetBrandNameScript(nonfungibleAddr, ownerAddr), [][]byte{jsoncdc.MustEncode(cadence.NewUInt64(one))})
		/// remove all characters other than alphabets
		assert.EqualValues(test, reg.ReplaceAllString(BrandMetadataValue, ""), reg.ReplaceAllString(ReturnedbrandModel.String(), ""))
	})

}

// Update Brand Not successful id ghalat

// Description: Create Schema with the name as input
// Input: SchemaName
// Expected Output: Schema Empty
// Test-case-type: Positive

func Test_CreateSchema_InstanceOnDeployment(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, _, _ := NFTContractDeployContracts(emulator, test)

	// Pre Condition: Check brand not initialized
	test.Run("Should have Zero brand count", func(test *testing.T) {
		schemaCount := executeScriptAndCheck(test, emulator, GetSchema_CountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), schemaCount)
	})
}

// Description: Create Schema with the name as input
// Input: schemaName
// Expected Output: Schema Created
// Test-case-type: Positive

func Test_CreateSchema_Success(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	// Pre Condition: Check Schema not initialized
	test.Run("Should have Zero brand count", func(test *testing.T) {
		schemaCount := executeScriptAndCheck(test, emulator, GetSchema_CountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), schemaCount)
	})

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	// Create Schema Transaction
	CreateSchema_Transaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, SchemaName)

	// Post Condition: Check Schema initialized properly
	test.Run("Should have initialized Schema correctly", func(test *testing.T) {
		schemaCount := executeScriptAndCheck(test, emulator, GetSchema_CountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(one), schemaCount)
	})
}

// Description: Create Schema with Empty Format
// Input: schemaName
// Expected Output: Schema should not Create
// Test-case-type: Negative

func Test_CreateSchema_InvalidStructure(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	// Pre Condition: Check Schema not initialized
	test.Run("Should have Zero Schema count", func(test *testing.T) {
		schemaCount := executeScriptAndCheck(test, emulator, GetSchema_CountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), schemaCount)
	})

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	// Create Schema Transaction
	CreateSchema_Transaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, SchemaName)

	// Post Condition: Check Schema initialized properly
	test.Run("Should have initialized Schema correctly", func(test *testing.T) {
		schemaCount := executeScriptAndCheck(test, emulator, GetSchema_CountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(one), schemaCount)
	})
}

// Schema with no format

// Schema with type Not supported in our Contract

// Description: Create template with brandId and schemaId as input
// Input: brandId,SchemaId
// Expected Output: Template Not Created
// Test-case-type: Negative
func Test_CreateTemplate_WithoutBrandID(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	TemplateNameField, _ := cadence.NewString(TemplateMetdataKey)

	//NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, Name, play)
	CreateSchema_Transaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, TemplateName)

	// metadata
	metadatatemplate := []cadence.KeyValuePair{{Key: TemplateNameField, Value: cadence.NewUInt64(two)}}

	// Create Template Transaction with brand ID:1 and schema ID:1 and 2 max Supply
	NowwhereCreateTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldFail,
		one, // brand ID
		one, // Schema ID
		two, // Max Supply
		metadatatemplate)

	// Post Condition: Check brand initialized properly
	test.Run("Should have initialized Template correctly", func(test *testing.T) {
		schemaCount := executeScriptAndCheck(test, emulator, NowwhereGenerateGetTemplateCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), schemaCount)
	})
}

// Description: Create template with brandId and schemaId as input
// Input: BrandId, schemaId(Not exist), MaxSupply
// Expected Output: Template Not Created
// Test-case-type: Negative
func Test_CreateTemplate_WithoutSchemaID(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brandField := cadence.NewDictionary(metadata)
	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, BrandName, brandField)
	// metadata
	metadatatemplate := []cadence.KeyValuePair{{Key: brandNameField, Value: cadence.NewUInt64(2)}}
	// Create Template Transaction with brand ID:1 and schema ID:1 and 2 max Supply
	NowwhereCreateTemplateTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldFail, one, one, two, metadatatemplate)

	// Post Condition: Check brand initialized properly
	test.Run("Should not have initialized Schema correctly", func(test *testing.T) {
		schemaCount := executeScriptAndCheck(test, emulator, NowwhereGenerateGetTemplateCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), schemaCount)
	})
}

// Description: Create template with brandId and schemaId as input
// Input: BrandId, schemaId, MaxSupply
// Expected Output: Template Created
// Test-case-type: Positive
func Test_CreateTemplate_withOwnBrandIDandSchemaID(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brandMetadata := cadence.NewDictionary(metadata)
	TemplateField, _ := cadence.NewString(SchemaName) // Metdata Template Field

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, BrandName, brandMetadata)
	CreateSchema_Transaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, SchemaName)
	// metadata
	metadatatemplate := []cadence.KeyValuePair{{Key: TemplateField, Value: TemplateField}}

	// Create Template Transaction  with brand ID:1 and schema ID:1 and 2 max Supply
	NowwhereCreateTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail,
		one, // brand ID
		one, // Schema ID
		two, // Max Supply
		metadatatemplate)

	// Post Condition: Check brand initialized properly
	test.Run("Should have initialized template correctly", func(test *testing.T) {
		templateCount := executeScriptAndCheck(test, emulator, NowwhereGenerateGetTemplateCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(one), templateCount)
	})
}

// Description: Create template with brandId and schemaId as input and max supply zero to check it still creates template
// Input: BrandId, schemaId, MaxSupply=0
// Expected Output: Template not Created
// Test-case-type: Negative
func Test_CreateTemplate_withzeroSupply(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brandMetadata := cadence.NewDictionary(metadata)
	TemplateField, _ := cadence.NewString(SchemaName) // Metdata Template Field

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, BrandName, brandMetadata)
	CreateSchema_Transaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, SchemaName)
	// metadata
	metadatatemplate := []cadence.KeyValuePair{{Key: TemplateField, Value: TemplateField}}
	// Create Template Transaction  with brand ID:1 and schema ID:1 and 2 max Supply
	NowwhereCreateTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldFail,
		one,  // brand ID
		one,  // Schema ID
		zero, // Max Supply
		metadatatemplate)

	// Post Condition: Check brand initialized properly
	test.Run("Should have initialized template correctly", func(test *testing.T) {
		templateCount := executeScriptAndCheck(test, emulator, NowwhereGenerateGetTemplateCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), templateCount)
	})
}

// Description: Create template with brandId and schemaId as input and give template not according to schema required
// Input: BrandId, schemaId, MaxSupply=0
// Expected Output: Template not Created
// Test-case-type: Negative
func Test_CreateTemplate_SchemaCheck(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brandMetadata := cadence.NewDictionary(metadata)
	WrongTemplateField, _ := cadence.NewString(TemplateMetdataKey) // wrong Metdata Template Field should be SchemaName

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, BrandName, brandMetadata)
	CreateSchema_Transaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, SchemaName)
	// metadata
	metadatatemplate := []cadence.KeyValuePair{{Key: WrongTemplateField, Value: cadence.NewInt(two)}}
	// Create Template Transaction  with brand ID:1 and schema ID:1 and 2 max Supply
	NowwhereCreateTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldFail,
		one, // brand ID
		one, // Schema ID
		two, // Max Supply
		metadatatemplate)

	// Post Condition: Check brand initialized properly
	test.Run("Should have initialized template correctly", func(test *testing.T) {
		templateCount := executeScriptAndCheck(test, emulator, NowwhereGenerateGetTemplateCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), templateCount)
	})
}

// Description: Create template with Normal users(That is not admin) brandId and schemaId as input(already created by Admin)
// Input: brandId, schemaId, MaxSupply
// Expected Output: Template Not Created
// Test-case-type: Negative
func Test_CreateTemplate_withOthersBrandIDandSchemaID(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, adminSigner, _ := NFTContractDeployContracts(emulator, test)

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brandMetadata := cadence.NewDictionary(metadata)
	TemplateField, _ := cadence.NewString(SchemaName) // Metdata Template Field

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, adminSigner)

	NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, adminSigner, shouldNotFail, BrandName, brandMetadata)
	CreateSchema_Transaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, adminSigner, shouldNotFail, SchemaName)
	// metadata
	metadatatemplate := []cadence.KeyValuePair{{Key: TemplateField, Value: TemplateField}}
	// Normal User
	userAddress, usersigner := NFTContractSetupAccount(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail)

	// Create Template Transaction  with brand ID:1 and schema ID:1 and 2 max Supply
	NowwhereCreateTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		userAddress,
		usersigner,
		shouldFail,
		one, // brand ID
		one, // Schema ID
		two, // Max Supply
		metadatatemplate)

	// Post Condition: Check brand initialized properly
	test.Run("Should not have initialized template with others brand correctly", func(test *testing.T) {
		templateCount := executeScriptAndCheck(test, emulator, NowwhereGenerateGetTemplateCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), templateCount)
	})
}

// Description: Create template with own brandId and other users schemaId as input
// Input: brandId, schemaId, MaxSupply
// Expected Output: Template Not Created
// Test-case-type: Negative
func Test_CreateTemplate_withOwnBrandIDandOthersSchemaID(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, adminSigner, _ := NFTContractDeployContracts(emulator, test)

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brandMetadata := cadence.NewDictionary(metadata)

	// 1st account(admin) configure
	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, adminSigner)

	// 2nd account(admin) configure
	userAddress, usersigner := NFTContractSetupNewAdminAccount(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		shouldNotFail, // if error arise throw error
	)

	NFTContractAddAdminCapability(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		usersigner,
		shouldNotFail, // if error arise throw error
		userAddress,   // setup Admin account to that address
	)

	NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, userAddress, usersigner, shouldNotFail, BrandName, brandMetadata)
	// Post Condition: Check brand initialized properly
	test.Run("Should have initialized brand correctly", func(test *testing.T) {
		brandCount := executeScriptAndCheck(test, emulator, NFTContractGenerateGetBrandCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(one), brandCount)
	})

	CreateSchema_Transaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		adminSigner,
		shouldNotFail,
		SchemaName,
	)

	// Post Condition: Check brand initialized properly
	test.Run("Should have initialized Schema correctly", func(test *testing.T) {
		schemaCount := executeScriptAndCheck(test, emulator, GetSchema_CountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(one), schemaCount)
	})

	// metadata
	metadatatemplate := []cadence.KeyValuePair{{Key: brandNameField, Value: cadence.NewUInt64(two)}}

	// Create Template Transaction  with brand ID:1 and schema ID:1 and 2 max Supply
	NowwhereCreateTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		userAddress,
		usersigner,
		shouldFail,
		one, // brand ID
		one, // Schema ID
		two, // Max Supply
		metadatatemplate)

	// Post Condition: Check brand initialized properly
	test.Run("Should not have initialized template with others brand correctly", func(test *testing.T) {
		templateCount := executeScriptAndCheck(test, emulator, NowwhereGenerateGetTemplateCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), templateCount)
	})
}

// Description: Mint Without TemplateId
// Input: templateId and receiver account
// Expected Output: NFT not minted
// Test-case-type: Negateive
func Test_MintTemplate_WithoutTemplate(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	// Create Setup New Account
	receiverAccount, _ := NFTContractSetupAccount(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail)
	var templateId uint64 = one
	NowwhereMintTemplateTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldFail, templateId, receiverAccount)

	// Post Condition: Check brand initialized properly
	test.Run("Should not have initialized Template correctly", func(test *testing.T) {
		templateCount := executeScriptAndCheck(test, emulator, NowwhereGenerateGetTemplateCountScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceInt(zero), templateCount)
	})
}

// Description: Mint template and send NFT to own Account
// Input: templateId and receiver account
// Expected Output: NFT minted
// Test-case-type: Positive

func Test_MintTemplate_Success(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brandMetadata := cadence.NewDictionary(metadata)
	TemplateField, _ := cadence.NewString(SchemaName) // Metdata Template Field

	NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, BrandName, brandMetadata)

	CreateSchema_Transaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, SchemaName)
	// metadata
	metadatatemplate := []cadence.KeyValuePair{{Key: TemplateField, Value: TemplateField}}
	// Create Template Transaction  with brand ID:1 and schema ID:1 and 2 max Supply
	NowwhereCreateTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail, // Throw error if occurs
		one,           // brand ID
		one,           // Schema ID
		two,           // Max Supply
		metadatatemplate)

	// Directly Minted to reciever address
	receiverAddress, _ := NFTContractSetupAccount(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail)
	var templateId uint64 = one
	NowwhereMintTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail,
		templateId,      // template ID wants to Mint
		receiverAddress, // Send NFT to that address
	)

	// NFT Supply
	test.Run("Should have initialized Supply field correctly", func(test *testing.T) {
		supply := executeScriptAndCheck(test, emulator, NowwhereGenerateGetSupplyScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceUInt64(one), supply)
	})

	// Post Condition: Check brand initialized properly
	test.Run("Should Mint Template correctly", func(test *testing.T) {
		MintCount := executeScriptAndCheck(test, emulator, NowwhereGenerateGetNFTAddressScript(nonfungibleAddr, ownerAddr),
			[][]byte{jsoncdc.MustEncode(cadence.Address(receiverAddress))})
		assert.EqualValues(test, CadenceArray(one), MintCount)
	})
}

// Description: Mint template and send NFT to Account Without setup user
// Input: templateId and receiver account
// Expected Output: NFT Not minted
// Test-case-type: Positive

func Test_MintTemplate_WithoutSetupUser(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brandMetadata := cadence.NewDictionary(metadata)
	TemplateField, _ := cadence.NewString(SchemaName) // Metdata Template Field

	NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, BrandName, brandMetadata)

	CreateSchema_Transaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, SchemaName)
	// metadata
	metadatatemplate := []cadence.KeyValuePair{{Key: TemplateField, Value: TemplateField}}
	// Create Template Transaction  with brand ID:1 and schema ID:1 and 2 max Supply
	NowwhereCreateTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail, // This transaction should not fail
		one,           // brand ID
		one,           // Schema ID
		two,           // Max Supply
		metadatatemplate)

	receiverAddress := GenerateAddress(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail)

	// Directly Minted to reciever address
	//receiverAddress, _ := NFTContractSetupAccount(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail)
	var templateId uint64 = one
	NowwhereMintTemplateTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldFail, templateId, receiverAddress)

	// NFT Supply
	test.Run("Should not have initialized Supply field correctly", func(test *testing.T) {
		supply := executeScriptAndCheck(test, emulator, NowwhereGenerateGetSupplyScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceUInt64(zero), supply)
	})

}

// Description: Mint template and send NFT to own Account, after that send that NFT to other Account
// Input:  receiver account and NFTID
// Expected Output: NFT minted
// Test-case-type: Positive
func Test_TransferNFT_Success(test *testing.T) {
	emulator := newEmulator()

	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brandMetadata := cadence.NewDictionary(metadata)
	TemplateField, _ := cadence.NewString(SchemaName) // Metdata Template Field

	NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, BrandName, brandMetadata)
	CreateSchema_Transaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, SchemaName)
	// metadata
	metadatatemplate := []cadence.KeyValuePair{{Key: TemplateField, Value: TemplateField}}
	// Create Template Transaction  with brand ID:1 and schema ID:1 and 2 max Supply
	NowwhereCreateTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail,
		one, // brand ID
		one, // Schema ID
		two, // Max Supply
		metadatatemplate)

	var templateId uint64 = one
	// Mint NFT to own Address
	NowwhereMintTemplateTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, templateId, ownerAddr)
	// Setup Account
	receiverAddress, _ := NFTContractSetupAccount(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail)
	// Post Condition: Check brand initialized properly
	test.Run("Should have initialized brand correctly", func(test *testing.T) {
		MintCount := executeScriptAndCheck(test, emulator, NowwhereGenerateGetNFTAddressScript(nonfungibleAddr, ownerAddr),
			[][]byte{jsoncdc.MustEncode(cadence.Address(ownerAddr))})
		assert.EqualValues(test, CadenceArray(one), MintCount)
	})

	// NFT transfer with parameter NFT ID = 1
	NFTContractTransferNFT(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail,
		receiverAddress,
		one) // NFT ID

	// NFT Supply
	test.Run("Should have initialized NFT Supply field correctly after Minting", func(test *testing.T) {
		supply := executeScriptAndCheck(test, emulator, NowwhereGenerateGetSupplyScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceUInt64(1), supply)
	})

	// Post Condition: Check brand initialized properly
	test.Run("Should have transferred NFT correctly", func(test *testing.T) {
		MintCount := executeScriptAndCheck(test, emulator, NowwhereGenerateGetNFTAddressScript(nonfungibleAddr, ownerAddr),
			[][]byte{jsoncdc.MustEncode(cadence.Address(receiverAddress))})
		assert.EqualValues(test, CadenceArray(one), MintCount)
	})
}

// Description: Mint template more than supply should give error and supply should be equal to max supply
// Input: templateId and receiver account
// Expected Output: max supply NFT minted and stops minting after that
// Test-case-type: Negative

func Test_MintTemplate_morethanSupply(test *testing.T) {
	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brandMetadata := cadence.NewDictionary(metadata)
	TemplateField, _ := cadence.NewString(SchemaName) // Metdata Template Field

	NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, BrandName, brandMetadata)
	CreateSchema_Transaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, SchemaName)
	// metadata
	metadatatemplate := []cadence.KeyValuePair{{Key: TemplateField, Value: TemplateField}}
	// Create Template Transaction  with brand ID:1 and schema ID:1 and 2 max Supply
	NowwhereCreateTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail,
		one, // brand ID
		one, // Schema ID
		two, // Max Supply
		metadatatemplate)

	// Pre Condition
	test.Run("Should have initialized Supply field Zero correctly", func(test *testing.T) {
		supply := executeScriptAndCheck(test, emulator, NowwhereGenerateGetSupplyScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceUInt64(zero), supply)
	})

	// Directly Minted to reciever address
	var templateId uint64 = one
	// Mint transaction # 1
	NowwhereMintTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail,
		templateId,
		ownerAddr)

	test.Run("Should have initialized Supply field One correctly", func(test *testing.T) {
		supply := executeScriptAndCheck(test, emulator, NowwhereGenerateGetSupplyScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceUInt64(one), supply)
	})

	// Mint transaction # 2
	NowwhereMintTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail,
		templateId,
		ownerAddr)

	test.Run("Should have initialized Supply field Two  correctly", func(test *testing.T) {
		supply := executeScriptAndCheck(test, emulator, NowwhereGenerateGetSupplyScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceUInt64(two), supply)
	})

	// Mint transaction # 3
	NowwhereMintTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldFail,
		templateId,
		ownerAddr)

	// Check NFT Supply
	test.Run("Should have initialized Supply field two correctly", func(test *testing.T) {
		supply := executeScriptAndCheck(test, emulator, NowwhereGenerateGetSupplyScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceUInt64(two), supply)
	})

}

// Description:Destroy NFT
// Input: templateId and receiver account
// Expected Output: Own NFT Destroyed
// Test-case-type: Positive
func Test_DestroyNFT_Success(test *testing.T) {

	emulator := newEmulator()
	nonfungibleAddr, ownerAddr, signer, _ := NFTContractDeployContracts(emulator, test)

	// 1st account(admin) configure
	SetupAdminAndGiveCapability(test, emulator, nonfungibleAddr, ownerAddr, shouldNotFail, signer)

	brandNameField, _ := cadence.NewString(BrandMetadataKey)
	brandName, _ := cadence.NewString(BrandMetadataValue)
	metadata := []cadence.KeyValuePair{{Key: brandNameField, Value: brandName}}
	brandMetadata := cadence.NewDictionary(metadata)
	TemplateField, _ := cadence.NewString(SchemaName) // Metdata Template Field

	NFTContractCreateBrandTransaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, BrandName, brandMetadata)
	CreateSchema_Transaction(test, emulator, nonfungibleAddr, ownerAddr, ownerAddr, signer, shouldNotFail, SchemaName)
	// metadata
	metadatatemplate := []cadence.KeyValuePair{{Key: TemplateField, Value: TemplateField}}
	// Create Template Transaction  with brand ID:1 and schema ID:1 and 2 max Supply
	NowwhereCreateTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail,
		one, // brand ID
		one, // Schema ID
		two, // Max Supply
		metadatatemplate)

	// Pre Condition
	test.Run("Should have initialized Supply field zero correctly", func(test *testing.T) {
		supply := executeScriptAndCheck(test, emulator, NowwhereGenerateGetSupplyScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceUInt64(zero), supply)
	})

	// Directly Minted to reciever address
	var templateId uint64 = one
	// Mint transaction # 1
	NowwhereMintTemplateTransaction(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail,
		templateId,
		ownerAddr)
	// Mint Should received
	test.Run("Should have initialized Supply field one correctly", func(test *testing.T) {
		supply := executeScriptAndCheck(test, emulator, NowwhereGenerateGetSupplyScript(nonfungibleAddr, ownerAddr), nil)
		assert.EqualValues(test, CadenceUInt64(one), supply)
	})
	test.Run("Should Mint Template correctly", func(test *testing.T) {
		MintCount := executeScriptAndCheck(test, emulator, NowwhereGenerateGetNFTAddressScript(nonfungibleAddr, ownerAddr),
			[][]byte{jsoncdc.MustEncode(cadence.Address(ownerAddr))})
		assert.EqualValues(test, CadenceArray(one), MintCount)
	})
	// NFT transfer with parameter NFT ID = 1
	NFTContractDestroyNFT(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		ownerAddr,
		signer,
		shouldNotFail,
		one) // NFT ID

	// Post Condition: Check brand initialized properly
	test.Run("Should Mint Template correctly", func(test *testing.T) {
		MintCount := executeScriptAndCheck(test, emulator, NowwhereGenerateGetNFTAddressScript(nonfungibleAddr, ownerAddr),
			[][]byte{jsoncdc.MustEncode(cadence.Address(ownerAddr))})
		assert.NotEqualValues(test, CadenceArray(one), MintCount)
	})

}

// Give capability to user after account setup
func SetupAdminAndGiveCapability(test *testing.T, emulator *emulator.Blockchain, nonfungibleAddr flow.Address, ownerAddr flow.Address, shouldNotFail bool, signer crypto.Signer) {
	NFTContractSetupAdminAccount(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		shouldNotFail,
		ownerAddr, // setup Admin account to that address
		signer,    // Signer of Admin Account
	)

	NFTContractAddAdminCapability(
		test,
		emulator,
		nonfungibleAddr,
		ownerAddr,
		signer,
		shouldNotFail,
		ownerAddr, // setup Admin account to that address
	)

}
