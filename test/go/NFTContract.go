package test

import (
	"regexp"
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ft_contracts "github.com/onflow/flow-ft/lib/go/contracts"
	nft_contracts "github.com/onflow/flow-nft/lib/go/contracts"
)

const (
	nowwhereRootPath                  = "../.."
	NFTContractPath                   = nowwhereRootPath + "/contracts/NFTContract.cdc"
	NowwhereContractPath              = nowwhereRootPath + "/contracts/NowwhereContract.cdc"
	NFTContractTransferTokensPath     = nowwhereRootPath + "/transactions/transferNFT.cdc"
	NFTContractDestroyTokensPath      = nowwhereRootPath + "/transactions/destroyNFT.cdc"
	NFTContractMintTokensPath         = nowwhereRootPath + "/transactions/mintNFT.cdc"
	NFTContractGetSupplyPath          = nowwhereRootPath + "/scripts/getTotalSupply.cdc"
	NFTContractGetCollectionPath      = nowwhereRootPath + "/scripts/getBrand.cdc"
	NFTContractGetCollectionCountPath = nowwhereRootPath + "/scripts/getBrandCount.cdc"
	NFTContractGetBrandNamePath       = nowwhereRootPath + "/scripts/getBrandName.cdc"
	NFTContractGetBrandIDPath         = nowwhereRootPath + "/scripts/getBrandIDs.cdc"
	NFTContractGetSchemaCountPath     = nowwhereRootPath + "/scripts/getSchemaCount.cdc"
	NFTContractGetTemplateCountPath   = nowwhereRootPath + "/scripts/getTemplateCount.cdc"
	NFTContractGetNFTAddressPath      = nowwhereRootPath + "/scripts/getNFTAddress.cdc"
	NFTContractGetNFTAddressCountPath = nowwhereRootPath + "/scripts/getAddressOwnedNFTCount.cdc"
	NFTContractCreateCollectionPath   = nowwhereRootPath + "/transactions/createBrand.cdc"
	NFTContractUpdateBrandPath        = nowwhereRootPath + "/transactions/UpdateBrand.cdc"
	NFTContractCreateSchemaPath       = nowwhereRootPath + "/transactions/createSchema.cdc"
	NFTContractCreateTemplatePath     = nowwhereRootPath + "/transactions/createTemplate.cdc"
	NFTContractSetupAccountPath       = nowwhereRootPath + "/transactions/setupAccount.cdc"
	NFTContractSetupAdminAccountPath  = nowwhereRootPath + "/transactions/setupAdminAccount.cdc"
	NFTContractAddAdminCapabilityPath = nowwhereRootPath + "/transactions/addAdminAccount.cdc"
	NFTContractCreateDropPath         = nowwhereRootPath + "/transactions/createDrop.cdc"
	NowwherePurchaseDropPath          = nowwhereRootPath + "/transactions/purchaseDrop.cdc"
	NowwhereRemoveDropPath            = nowwhereRootPath + "/transactions/RemoveDrop.cdc"
	CapabilityAdminCheck              = nowwhereRootPath + "/transactions/CheckAdminCapability.cdc"
	NowwhereContractgetDropCountPath  = nowwhereRootPath + "/scripts/getDropCount.cdc"
	NowwhereContractgetDropIdsPath    = nowwhereRootPath + "/scripts/getDropIds.cdc"
	getDate                           = nowwhereRootPath + "/scripts/getDate.cdc"
)

func NFTContractDeployContracts(emulator *emulator.Blockchain, testing *testing.T) (flow.Address, flow.Address, crypto.Signer, sdk.Address) {
	accountKeys := test.AccountKeyGenerator()
	adminAccountKey, adminSigner := accountKeys.NewWithSigner()

	nftCode := loadNonFungibleToken()
	nftAddr, err := emulator.CreateAccount(
		[]*flow.AccountKey{adminAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "NonFungibleToken",
				Source: string(nftCode),
			},
		},
	)
	require.NoError(testing, err)

	_, err = emulator.CommitBlock()
	assert.NoError(testing, err)

	address, err := emulator.CreateAccount([]*sdk.AccountKey{adminAccountKey}, nil)

	NFTContractCode := loadNFTContract(nftAddr.String())

	adminAddr, err := emulator.CreateAccount(
		[]*flow.AccountKey{adminAccountKey},
		[]templates.Contract{templates.Contract{
			Name:   "NFTContract",
			Source: string(NFTContractCode),
		}},
	)
	assert.NoError(testing, err)

	_, err = emulator.CommitBlock()
	assert.NoError(testing, err)

	return nftAddr, adminAddr, adminSigner, address
}

func nowwhereReplaceAddressPlaceholders(code string, nonfungibleAddress, NFTContractAddress string) []byte {
	return []byte(replaceImports(
		code,
		map[string]*regexp.Regexp{
			nonfungibleAddress: nftAddressPlaceholder,
			NFTContractAddress: NFTContractAddressPlaceHolder,
		},
	))
}

func nowwhereContractReplaceAddressPlaceholders(code string, nonfungibleAddress, NFTContractAddress, nowwhereAddress string) []byte {
	return []byte(replaceImports(
		code,
		map[string]*regexp.Regexp{
			nonfungibleAddress: nftAddressPlaceholder,
			NFTContractAddress: NFTContractAddressPlaceHolder,
			nowwhereAddress:    NowwherePlaceholder,
		},
	))
}

func loadFungibleToken() []byte {
	return ft_contracts.FungibleToken()
}

func loadNFTContract(nftAddr string) []byte {
	return []byte(replaceImports(
		string(readFile(NFTContractPath)),
		map[string]*regexp.Regexp{
			nftAddr: nftAddressPlaceholder,
		},
	))
}
func loadNowwhereContract(nftAddr string, NFTContractAddr string) []byte {
	return []byte(replaceImports(
		string(readFile(NowwhereContractPath)),
		map[string]*regexp.Regexp{
			nftAddr:         nftAddressPlaceholder,
			NFTContractAddr: NFTContractAddressPlaceHolder,
		},
	))
}

func loadNFT(fungibleAddr flow.Address) []byte {
	return []byte(replaceImports(
		string(readFile(NFTContractPath)),
		map[string]*regexp.Regexp{
			fungibleAddr.String(): ftAddressPlaceholder,
		},
	))
}

func NowwhereGenerateGetSupplyScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractGetSupplyPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereGenerateGetCollectionScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractGetCollectionPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractGenerateGetBrandCountScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractGetCollectionCountPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractGenerateGetBrandNameScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractGetBrandNamePath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractGenerateGetBrandIDsScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractGetBrandIDPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func GetSchema_CountScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractGetSchemaCountPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

// Template Script
func NowwhereGenerateGetTemplateCountScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractGetTemplateCountPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

// Drops Script
func NowwhereGenerateGetDropCountScript(fungibleAddr, NFTContract, nowwhereAddr flow.Address) []byte {
	return nowwhereContractReplaceAddressPlaceholders(
		string(readFile(NowwhereContractgetDropCountPath)),
		fungibleAddr.String(),
		NFTContract.String(),
		nowwhereAddr.String(),
	)
}

// Drops Script
func NowwhereGenerateGetDropIdsScript(fungibleAddr, NFTContract, nowwhereAddr flow.Address) []byte {
	return nowwhereContractReplaceAddressPlaceholders(
		string(readFile(NowwhereContractgetDropIdsPath)),
		fungibleAddr.String(),
		NFTContract.String(),
		nowwhereAddr.String(),
	)
}

func getCurrentTime(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(getDate)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereGenerateGetNFTAddressScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractGetNFTAddressPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereGenerateGetNFTAddressCountScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractGetNFTAddressCountPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func loadNonFungibleToken() []byte {
	return nft_contracts.NonFungibleToken()
}

func NowwhereCreateGenerateCollectionScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractCreateCollectionPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func CapabilityAccessScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(CapabilityAdminCheck)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereUpdateBrandScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractUpdateBrandPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereCreateGenerateSchemaScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractCreateSchemaPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereCreateGenerateTemplateScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractCreateTemplatePath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractSetupAccountScript(nonfungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractSetupAccountPath)),
		nonfungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereSetupAdminAccountScript(nonfungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractSetupAdminAccountPath)),
		nonfungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractAddAdminCapabilityScript(nonfungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractAddAdminCapabilityPath)),
		nonfungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractTransferNFTScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractTransferTokensPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractDestroyNFTScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractDestroyTokensPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereMintTokensScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractMintTokensPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func CheckCapabilityTransaction(
	testing *testing.T,
	emulator *emulator.Blockchain,
	fungibleAddr,
	nowwhereAddr flow.Address,
	userAddress sdk.Address,
	userSigner crypto.Signer,
	shouldFail bool,
) {
	tx := flow.NewTransaction().
		SetScript(CapabilityAccessScript(fungibleAddr, nowwhereAddr)).
		SetGasLimit(100).
		SetProposalKey(emulator.ServiceKey().Address, emulator.ServiceKey().Index, emulator.ServiceKey().SequenceNumber).
		SetPayer(emulator.ServiceKey().Address).
		AddAuthorizer(userAddress)

	signAndSubmit(
		testing, emulator, tx,
		[]flow.Address{emulator.ServiceKey().Address, userAddress},
		[]crypto.Signer{emulator.ServiceKey().Signer(), userSigner},
		shouldFail,
	)
}

func NFTContractCreateBrandTransaction(
	testing *testing.T,
	emulator *emulator.Blockchain,
	fungibleAddr,
	nowwhereAddr flow.Address,
	userAddress sdk.Address,
	userSigner crypto.Signer,
	shouldFail bool,
	brandName string,
	metaData cadence.Dictionary,
) {
	tx := flow.NewTransaction().
		SetScript(NowwhereCreateGenerateCollectionScript(fungibleAddr, nowwhereAddr)).
		SetGasLimit(100).
		SetProposalKey(emulator.ServiceKey().Address, emulator.ServiceKey().Index, emulator.ServiceKey().SequenceNumber).
		SetPayer(emulator.ServiceKey().Address).
		AddAuthorizer(userAddress)

	brand, _ := cadence.NewString(brandName)

	_ = tx.AddArgument(brand)    // brandName
	_ = tx.AddArgument(metaData) // Metadata

	signAndSubmit(
		testing, emulator, tx,
		[]flow.Address{emulator.ServiceKey().Address, userAddress},
		[]crypto.Signer{emulator.ServiceKey().Signer(), userSigner},
		shouldFail,
	)
}

func NFTContractUpdateBrandTransaction(
	testing *testing.T,
	emulator *emulator.Blockchain,
	fungibleAddr,
	nowwhereAddr flow.Address,
	userAddress sdk.Address,
	userSigner crypto.Signer,
	shouldFail bool,
	brandId int,
	brandNameToUpdate string,
) {
	tx := flow.NewTransaction().
		SetScript(NowwhereUpdateBrandScript(fungibleAddr, nowwhereAddr)).
		SetGasLimit(100).
		SetProposalKey(emulator.ServiceKey().Address, emulator.ServiceKey().Index, emulator.ServiceKey().SequenceNumber).
		SetPayer(emulator.ServiceKey().Address).
		AddAuthorizer(userAddress)

	brand := cadence.UInt64(uint64(brandId))

	_ = tx.AddArgument(brand) // brandName
	_ = tx.AddArgument(CadenceString(brandNameToUpdate))

	signAndSubmit(
		testing, emulator, tx,
		[]flow.Address{emulator.ServiceKey().Address, userAddress},
		[]crypto.Signer{emulator.ServiceKey().Signer(), userSigner},
		shouldFail,
	)
}

func CreateSchema_Transaction(
	testing *testing.T,
	emulator *emulator.Blockchain,
	fungibleAddr,
	nowwhereAddr flow.Address,
	userAddress sdk.Address,
	userSigner crypto.Signer,
	shouldFail bool,
	schemaName string,
) {
	tx := flow.NewTransaction().
		SetScript(NowwhereCreateGenerateSchemaScript(fungibleAddr, nowwhereAddr)).
		SetGasLimit(100).
		SetProposalKey(emulator.ServiceKey().Address, emulator.ServiceKey().Index, emulator.ServiceKey().SequenceNumber).
		SetPayer(emulator.ServiceKey().Address).
		AddAuthorizer(userAddress)
	schema, _ := cadence.NewString(schemaName)

	_ = tx.AddArgument(schema)

	signAndSubmit(
		testing, emulator, tx,
		[]flow.Address{emulator.ServiceKey().Address, userAddress},
		[]crypto.Signer{emulator.ServiceKey().Signer(), userSigner},
		shouldFail,
	)
}

func NowwhereCreateTemplateTransaction(
	testing *testing.T,
	emulator *emulator.Blockchain,
	fungibleAddr,
	nowwhereAddr flow.Address,
	userAddress sdk.Address,
	userSigner crypto.Signer,
	shouldFail bool,
	collectionId uint64,
	schemaId uint64,
	maxSupply uint64,
	metadata []cadence.KeyValuePair,
) {
	tx := flow.NewTransaction().
		SetScript(NowwhereCreateGenerateTemplateScript(fungibleAddr, nowwhereAddr)).
		SetGasLimit(100).
		SetProposalKey(emulator.ServiceKey().Address, emulator.ServiceKey().Index, emulator.ServiceKey().SequenceNumber).
		SetPayer(emulator.ServiceKey().Address).
		AddAuthorizer(userAddress)

	_ = tx.AddArgument(cadence.NewUInt64(collectionId))
	_ = tx.AddArgument(cadence.NewUInt64(schemaId))
	_ = tx.AddArgument(cadence.NewUInt64(maxSupply))
	_ = tx.AddArgument(cadence.NewDictionary(metadata))

	signAndSubmit(
		testing, emulator, tx,
		[]flow.Address{emulator.ServiceKey().Address, userAddress},
		[]crypto.Signer{emulator.ServiceKey().Signer(), userSigner},
		shouldFail,
	)
}

func NowwhereMintTemplateTransaction(
	testing *testing.T,
	emulator *emulator.Blockchain,
	fungibleAddr,
	nowwhereAddr flow.Address,
	userAddress sdk.Address,
	userSigner crypto.Signer,
	shouldFail bool,
	templateId uint64,
	receiverAccount sdk.Address,
) {
	tx := flow.NewTransaction().
		SetScript(NowwhereMintTokensScript(fungibleAddr, nowwhereAddr)).
		SetGasLimit(100).
		SetProposalKey(emulator.ServiceKey().Address, emulator.ServiceKey().Index, emulator.ServiceKey().SequenceNumber).
		SetPayer(emulator.ServiceKey().Address).
		AddAuthorizer(userAddress)

	_ = tx.AddArgument(cadence.NewUInt64(templateId))
	_ = tx.AddArgument(cadence.NewAddress(receiverAccount))

	signAndSubmit(
		testing, emulator, tx,
		[]flow.Address{emulator.ServiceKey().Address, userAddress},
		[]crypto.Signer{emulator.ServiceKey().Signer(), userSigner},
		shouldFail,
	)
}

func NFTContractSetupAccount(
	testing *testing.T,
	emulator *emulator.Blockchain,
	nonfungibleAddr,
	nowwhereAddr sdk.Address,
	shouldFail bool,
) (sdk.Address, crypto.Signer) {
	accountKeys := test.AccountKeyGenerator()
	AccountKey, Signer := accountKeys.NewWithSigner()
	address, _ := emulator.CreateAccount([]*sdk.AccountKey{AccountKey}, nil)

	tx := flow.NewTransaction().
		SetScript(NFTContractSetupAccountScript(nonfungibleAddr, nowwhereAddr)).
		SetGasLimit(100).
		SetProposalKey(emulator.ServiceKey().Address, emulator.ServiceKey().Index, emulator.ServiceKey().SequenceNumber).
		SetPayer(emulator.ServiceKey().Address).
		AddAuthorizer(address)

	signAndSubmit(
		testing, emulator, tx,
		[]flow.Address{emulator.ServiceKey().Address, address},
		[]crypto.Signer{emulator.ServiceKey().Signer(), Signer},
		false,
	)

	return address, Signer
}

func GenerateAddress(
	testing *testing.T,
	emulator *emulator.Blockchain,
	nonfungibleAddr,
	nowwhereAddr sdk.Address,
	shouldFail bool,
) sdk.Address {
	accountKeys := test.AccountKeyGenerator()
	AccountKey, _ := accountKeys.NewWithSigner()
	address, _ := emulator.CreateAccount([]*sdk.AccountKey{AccountKey}, nil)

	return address
}

func NFTContractTransferNFT(
	testing *testing.T,
	emulator *emulator.Blockchain,
	nonfungibleAddr,
	nowwhereAddr sdk.Address,
	userAddress sdk.Address,
	userSigner crypto.Signer,
	shouldFail bool,
	recieverAddress sdk.Address,
	NFTId uint64,
) {

	tx := flow.NewTransaction().
		SetScript(NFTContractTransferNFTScript(nonfungibleAddr, nowwhereAddr)).
		SetGasLimit(100).
		SetProposalKey(emulator.ServiceKey().Address, emulator.ServiceKey().Index, emulator.ServiceKey().SequenceNumber).
		SetPayer(emulator.ServiceKey().Address).
		AddAuthorizer(userAddress)
	_ = tx.AddArgument(cadence.NewAddress(recieverAddress))
	_ = tx.AddArgument(cadence.NewUInt64(NFTId))
	signAndSubmit(
		testing, emulator, tx,
		[]flow.Address{emulator.ServiceKey().Address, userAddress},
		[]crypto.Signer{emulator.ServiceKey().Signer(), userSigner},
		false,
	)

}

func NFTContractDestroyNFT(
	testing *testing.T,
	emulator *emulator.Blockchain,
	nonfungibleAddr,
	nowwhereAddr sdk.Address,
	userAddress sdk.Address,
	userSigner crypto.Signer,
	shouldFail bool,
	NFTId uint64,
) {

	tx := flow.NewTransaction().
		SetScript(NFTContractDestroyNFTScript(nonfungibleAddr, nowwhereAddr)).
		SetGasLimit(100).
		SetProposalKey(emulator.ServiceKey().Address, emulator.ServiceKey().Index, emulator.ServiceKey().SequenceNumber).
		SetPayer(emulator.ServiceKey().Address).
		AddAuthorizer(userAddress)
	_ = tx.AddArgument(cadence.NewUInt64(NFTId))
	signAndSubmit(
		testing, emulator, tx,
		[]flow.Address{emulator.ServiceKey().Address, userAddress},
		[]crypto.Signer{emulator.ServiceKey().Signer(), userSigner},
		false,
	)

}

func NFTContractSetupNewAdminAccount(
	testing *testing.T,
	emulator *emulator.Blockchain,
	nonfungibleAddr,
	nowwhereAddr sdk.Address,
	shouldFail bool,
) (sdk.Address, crypto.Signer) {
	accountKeys := test.AccountKeyGenerator()
	AccountKey, Signer := accountKeys.NewWithSigner()
	address, _ := emulator.CreateAccount([]*sdk.AccountKey{AccountKey}, nil)

	tx := flow.NewTransaction().
		SetScript(NowwhereSetupAdminAccountScript(nonfungibleAddr, nowwhereAddr)).
		SetGasLimit(100).
		SetProposalKey(emulator.ServiceKey().Address, emulator.ServiceKey().Index, emulator.ServiceKey().SequenceNumber).
		SetPayer(emulator.ServiceKey().Address).
		AddAuthorizer(address)

	signAndSubmit(
		testing, emulator, tx,
		[]flow.Address{emulator.ServiceKey().Address, address},
		[]crypto.Signer{emulator.ServiceKey().Signer(), Signer},
		false,
	)

	return address, Signer
}

func NFTContractSetupAdminAccount(
	testing *testing.T,
	emulator *emulator.Blockchain,
	nonfungibleAddr,
	nowwhereAddr sdk.Address,
	shouldFail bool,
	adminAddress sdk.Address,
	Signer crypto.Signer,
) {

	tx := flow.NewTransaction().
		SetScript(NowwhereSetupAdminAccountScript(nonfungibleAddr, nowwhereAddr)).
		SetGasLimit(100).
		SetProposalKey(emulator.ServiceKey().Address, emulator.ServiceKey().Index, emulator.ServiceKey().SequenceNumber).
		SetPayer(emulator.ServiceKey().Address).
		AddAuthorizer(adminAddress)

	signAndSubmit(
		testing, emulator, tx,
		[]flow.Address{emulator.ServiceKey().Address, adminAddress},
		[]crypto.Signer{emulator.ServiceKey().Signer(), Signer},
		false,
	)

	return
}

func NFTContractAddAdminCapability(
	testing *testing.T,
	emulator *emulator.Blockchain,
	nonfungibleAddr,
	nowwhereAddr sdk.Address,
	userSigner crypto.Signer,
	shouldFail bool,
	adminAddress sdk.Address,
) {

	tx := flow.NewTransaction().
		SetScript(NFTContractAddAdminCapabilityScript(nonfungibleAddr, nowwhereAddr)).
		SetGasLimit(100).
		SetProposalKey(emulator.ServiceKey().Address, emulator.ServiceKey().Index, emulator.ServiceKey().SequenceNumber).
		SetPayer(emulator.ServiceKey().Address).
		AddAuthorizer(nowwhereAddr)

	_ = tx.AddArgument(cadence.NewAddress(adminAddress))

	signAndSubmit(
		testing, emulator, tx,
		[]flow.Address{emulator.ServiceKey().Address, nowwhereAddr},
		[]crypto.Signer{emulator.ServiceKey().Signer(), userSigner},
		false,
	)

	return
}
