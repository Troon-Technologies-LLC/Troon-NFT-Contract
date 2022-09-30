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
	TroonAtomicStandardPath                   = nowwhereRootPath + "/contracts/TroonAtomicStandard.cdc"
	NowwhereContractPath              = nowwhereRootPath + "/contracts/NowwhereContract.cdc"
	TroonAtomicStandardTransferTokensPath     = nowwhereRootPath + "/transactions/transferNFT.cdc"
	TroonAtomicStandardDestroyTokensPath      = nowwhereRootPath + "/transactions/destroyNFT.cdc"
	TroonAtomicStandardMintTokensPath         = nowwhereRootPath + "/transactions/mintNFT.cdc"
	TroonAtomicStandardGetSupplyPath          = nowwhereRootPath + "/scripts/getTotalSupply.cdc"
	TroonAtomicStandardGetCollectionPath      = nowwhereRootPath + "/scripts/getBrand.cdc"
	TroonAtomicStandardGetCollectionCountPath = nowwhereRootPath + "/scripts/getBrandCount.cdc"
	TroonAtomicStandardGetBrandNamePath       = nowwhereRootPath + "/scripts/getBrandName.cdc"
	TroonAtomicStandardGetBrandIDPath         = nowwhereRootPath + "/scripts/getBrandIDs.cdc"
	TroonAtomicStandardGetSchemaCountPath     = nowwhereRootPath + "/scripts/getSchemaCount.cdc"
	TroonAtomicStandardGetTemplateCountPath   = nowwhereRootPath + "/scripts/getTemplateCount.cdc"
	TroonAtomicStandardGetNFTAddressPath      = nowwhereRootPath + "/scripts/getNFTAddress.cdc"
	TroonAtomicStandardGetNFTAddressCountPath = nowwhereRootPath + "/scripts/getAddressOwnedNFTCount.cdc"
	TroonAtomicStandardCreateCollectionPath   = nowwhereRootPath + "/transactions/createBrand.cdc"
	TroonAtomicStandardUpdateBrandPath        = nowwhereRootPath + "/transactions/UpdateBrand.cdc"
	TroonAtomicStandardCreateSchemaPath       = nowwhereRootPath + "/transactions/createSchema.cdc"
	TroonAtomicStandardCreateTemplatePath     = nowwhereRootPath + "/transactions/createTemplate.cdc"
	TroonAtomicStandardSetupAccountPath       = nowwhereRootPath + "/transactions/setupAccount.cdc"
	TroonAtomicStandardSetupAdminAccountPath  = nowwhereRootPath + "/transactions/setupAdminAccount.cdc"
	TroonAtomicStandardAddAdminCapabilityPath = nowwhereRootPath + "/transactions/addAdminAccount.cdc"
	TroonAtomicStandardCreateDropPath         = nowwhereRootPath + "/transactions/createDrop.cdc"
	NowwherePurchaseDropPath          = nowwhereRootPath + "/transactions/purchaseDrop.cdc"
	NowwhereRemoveDropPath            = nowwhereRootPath + "/transactions/RemoveDrop.cdc"
	CapabilityAdminCheck              = nowwhereRootPath + "/transactions/CheckAdminCapability.cdc"
	NowwhereContractgetDropCountPath  = nowwhereRootPath + "/scripts/getDropCount.cdc"
	NowwhereContractgetDropIdsPath    = nowwhereRootPath + "/scripts/getDropIds.cdc"
	getDate                           = nowwhereRootPath + "/scripts/getDate.cdc"
)

func TroonAtomicStandardDeployContracts(emulator *emulator.Blockchain, testing *testing.T) (flow.Address, flow.Address, crypto.Signer, sdk.Address) {
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

	TroonAtomicStandardCode := loadTroonAtomicStandard(nftAddr.String())

	adminAddr, err := emulator.CreateAccount(
		[]*flow.AccountKey{adminAccountKey},
		[]templates.Contract{templates.Contract{
			Name:   "TroonAtomicStandard",
			Source: string(TroonAtomicStandardCode),
		}},
	)
	assert.NoError(testing, err)

	_, err = emulator.CommitBlock()
	assert.NoError(testing, err)

	return nftAddr, adminAddr, adminSigner, address
}

func nowwhereReplaceAddressPlaceholders(code string, nonfungibleAddress, TroonAtomicStandardAddress string) []byte {
	return []byte(replaceImports(
		code,
		map[string]*regexp.Regexp{
			nonfungibleAddress: nftAddressPlaceholder,
			TroonAtomicStandardAddress: TroonAtomicStandardAddressPlaceHolder,
		},
	))
}

func nowwhereContractReplaceAddressPlaceholders(code string, nonfungibleAddress, TroonAtomicStandardAddress, nowwhereAddress string) []byte {
	return []byte(replaceImports(
		code,
		map[string]*regexp.Regexp{
			nonfungibleAddress: nftAddressPlaceholder,
			TroonAtomicStandardAddress: TroonAtomicStandardAddressPlaceHolder,
			nowwhereAddress:    NowwherePlaceholder,
		},
	))
}

func loadFungibleToken() []byte {
	return ft_contracts.FungibleToken()
}

func loadTroonAtomicStandard(nftAddr string) []byte {
	return []byte(replaceImports(
		string(readFile(TroonAtomicStandardPath)),
		map[string]*regexp.Regexp{
			nftAddr: nftAddressPlaceholder,
		},
	))
}
func loadNowwhereContract(nftAddr string, TroonAtomicStandardAddr string) []byte {
	return []byte(replaceImports(
		string(readFile(NowwhereContractPath)),
		map[string]*regexp.Regexp{
			nftAddr:         nftAddressPlaceholder,
			TroonAtomicStandardAddr: TroonAtomicStandardAddressPlaceHolder,
		},
	))
}

func loadNFT(fungibleAddr flow.Address) []byte {
	return []byte(replaceImports(
		string(readFile(TroonAtomicStandardPath)),
		map[string]*regexp.Regexp{
			fungibleAddr.String(): ftAddressPlaceholder,
		},
	))
}

func NowwhereGenerateGetSupplyScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardGetSupplyPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereGenerateGetCollectionScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardGetCollectionPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func TroonAtomicStandardGenerateGetBrandCountScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardGetCollectionCountPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func TroonAtomicStandardGenerateGetBrandNameScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardGetBrandNamePath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func TroonAtomicStandardGenerateGetBrandIDsScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardGetBrandIDPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func GetSchema_CountScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardGetSchemaCountPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

// Template Script
func NowwhereGenerateGetTemplateCountScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardGetTemplateCountPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

// Drops Script
func NowwhereGenerateGetDropCountScript(fungibleAddr, TroonAtomicStandard, nowwhereAddr flow.Address) []byte {
	return nowwhereContractReplaceAddressPlaceholders(
		string(readFile(NowwhereContractgetDropCountPath)),
		fungibleAddr.String(),
		TroonAtomicStandard.String(),
		nowwhereAddr.String(),
	)
}

// Drops Script
func NowwhereGenerateGetDropIdsScript(fungibleAddr, TroonAtomicStandard, nowwhereAddr flow.Address) []byte {
	return nowwhereContractReplaceAddressPlaceholders(
		string(readFile(NowwhereContractgetDropIdsPath)),
		fungibleAddr.String(),
		TroonAtomicStandard.String(),
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
		string(readFile(TroonAtomicStandardGetNFTAddressPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereGenerateGetNFTAddressCountScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardGetNFTAddressCountPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func loadNonFungibleToken() []byte {
	return nft_contracts.NonFungibleToken()
}

func NowwhereCreateGenerateCollectionScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardCreateCollectionPath)),
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
		string(readFile(TroonAtomicStandardUpdateBrandPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereCreateGenerateSchemaScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardCreateSchemaPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereCreateGenerateTemplateScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardCreateTemplatePath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func TroonAtomicStandardSetupAccountScript(nonfungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardSetupAccountPath)),
		nonfungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereSetupAdminAccountScript(nonfungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardSetupAdminAccountPath)),
		nonfungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func TroonAtomicStandardAddAdminCapabilityScript(nonfungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardAddAdminCapabilityPath)),
		nonfungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func TroonAtomicStandardTransferNFTScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardTransferTokensPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func TroonAtomicStandardDestroyNFTScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardDestroyTokensPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereMintTokensScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(TroonAtomicStandardMintTokensPath)),
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

func TroonAtomicStandardCreateBrandTransaction(
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

func TroonAtomicStandardUpdateBrandTransaction(
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

func TroonAtomicStandardSetupAccount(
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
		SetScript(TroonAtomicStandardSetupAccountScript(nonfungibleAddr, nowwhereAddr)).
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

func TroonAtomicStandardTransferNFT(
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
		SetScript(TroonAtomicStandardTransferNFTScript(nonfungibleAddr, nowwhereAddr)).
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

func TroonAtomicStandardDestroyNFT(
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
		SetScript(TroonAtomicStandardDestroyNFTScript(nonfungibleAddr, nowwhereAddr)).
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

func TroonAtomicStandardSetupNewAdminAccount(
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

func TroonAtomicStandardSetupAdminAccount(
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

func TroonAtomicStandardAddAdminCapability(
	testing *testing.T,
	emulator *emulator.Blockchain,
	nonfungibleAddr,
	nowwhereAddr sdk.Address,
	userSigner crypto.Signer,
	shouldFail bool,
	adminAddress sdk.Address,
) {

	tx := flow.NewTransaction().
		SetScript(TroonAtomicStandardAddAdminCapabilityScript(nonfungibleAddr, nowwhereAddr)).
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
