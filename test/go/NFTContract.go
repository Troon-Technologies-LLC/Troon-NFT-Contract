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
	nowwhereRootPath                     = "../.."
	NFTContractV01Path                   = nowwhereRootPath + "/contracts/NFTContractV01.cdc"
	NowwhereContractPath                 = nowwhereRootPath + "/contracts/NowwhereContract.cdc"
	NFTContractV01TransferTokensPath     = nowwhereRootPath + "/transactions/transferNFT.cdc"
	NFTContractV01DestroyTokensPath      = nowwhereRootPath + "/transactions/destroyNFT.cdc"
	NFTContractV01MintTokensPath         = nowwhereRootPath + "/transactions/mintNFT.cdc"
	NFTContractV01GetSupplyPath          = nowwhereRootPath + "/scripts/getTotalSupply.cdc"
	NFTContractV01GetCollectionPath      = nowwhereRootPath + "/scripts/getBrand.cdc"
	NFTContractV01GetCollectionCountPath = nowwhereRootPath + "/scripts/getBrandCount.cdc"
	NFTContractV01GetBrandNamePath       = nowwhereRootPath + "/scripts/getBrandName.cdc"
	NFTContractV01GetBrandIDPath         = nowwhereRootPath + "/scripts/getBrandIDs.cdc"
	NFTContractV01GetSchemaCountPath     = nowwhereRootPath + "/scripts/getSchemaCount.cdc"
	NFTContractV01GetTemplateCountPath   = nowwhereRootPath + "/scripts/getTemplateCount.cdc"
	NFTContractV01GetNFTAddressPath      = nowwhereRootPath + "/scripts/getNFTAddress.cdc"
	NFTContractV01GetNFTAddressCountPath = nowwhereRootPath + "/scripts/getAddressOwnedNFTCount.cdc"
	NFTContractV01CreateCollectionPath   = nowwhereRootPath + "/transactions/createBrand.cdc"
	NFTContractV01UpdateBrandPath        = nowwhereRootPath + "/transactions/UpdateBrand.cdc"
	NFTContractV01CreateSchemaPath       = nowwhereRootPath + "/transactions/createSchema.cdc"
	NFTContractV01CreateTemplatePath     = nowwhereRootPath + "/transactions/createTemplate.cdc"
	NFTContractV01SetupAccountPath       = nowwhereRootPath + "/transactions/setupAccount.cdc"
	NFTContractV01SetupAdminAccountPath  = nowwhereRootPath + "/transactions/setupAdminAccount.cdc"
	NFTContractV01AddAdminCapabilityPath = nowwhereRootPath + "/transactions/addAdminAccount.cdc"
	NFTContractV01CreateDropPath         = nowwhereRootPath + "/transactions/createDrop.cdc"
	NowwherePurchaseDropPath             = nowwhereRootPath + "/transactions/purchaseDrop.cdc"
	NowwhereRemoveDropPath               = nowwhereRootPath + "/transactions/RemoveDrop.cdc"
	CapabilityAdminCheck                 = nowwhereRootPath + "/transactions/CheckAdminCapability.cdc"
	NowwhereContractgetDropCountPath     = nowwhereRootPath + "/scripts/getDropCount.cdc"
	NowwhereContractgetDropIdsPath       = nowwhereRootPath + "/scripts/getDropIds.cdc"
	getDate                              = nowwhereRootPath + "/scripts/getDate.cdc"
)

func NFTContractV01DeployContracts(emulator *emulator.Blockchain, testing *testing.T) (flow.Address, flow.Address, crypto.Signer, sdk.Address) {
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

	NFTContractV01Code := loadNFTContractV01(nftAddr.String())

	adminAddr, err := emulator.CreateAccount(
		[]*flow.AccountKey{adminAccountKey},
		[]templates.Contract{templates.Contract{
			Name:   "NFTContractV01",
			Source: string(NFTContractV01Code),
		}},
	)
	assert.NoError(testing, err)

	_, err = emulator.CommitBlock()
	assert.NoError(testing, err)

	return nftAddr, adminAddr, adminSigner, address
}

func nowwhereReplaceAddressPlaceholders(code string, nonfungibleAddress, NFTContractV01Address string) []byte {
	return []byte(replaceImports(
		code,
		map[string]*regexp.Regexp{
			nonfungibleAddress:    nftAddressPlaceholder,
			NFTContractV01Address: NFTContractV01AddressPlaceHolder,
		},
	))
}

func nowwhereContractReplaceAddressPlaceholders(code string, nonfungibleAddress, NFTContractV01Address, nowwhereAddress string) []byte {
	return []byte(replaceImports(
		code,
		map[string]*regexp.Regexp{
			nonfungibleAddress:    nftAddressPlaceholder,
			NFTContractV01Address: NFTContractV01AddressPlaceHolder,
			nowwhereAddress:       NowwherePlaceholder,
		},
	))
}

func loadFungibleToken() []byte {
	return ft_contracts.FungibleToken()
}

func loadNFTContractV01(nftAddr string) []byte {
	return []byte(replaceImports(
		string(readFile(NFTContractV01Path)),
		map[string]*regexp.Regexp{
			nftAddr: nftAddressPlaceholder,
		},
	))
}
func loadNowwhereContract(nftAddr string, NFTContractV01Addr string) []byte {
	return []byte(replaceImports(
		string(readFile(NowwhereContractPath)),
		map[string]*regexp.Regexp{
			nftAddr:            nftAddressPlaceholder,
			NFTContractV01Addr: NFTContractV01AddressPlaceHolder,
		},
	))
}

func loadNFT(fungibleAddr flow.Address) []byte {
	return []byte(replaceImports(
		string(readFile(NFTContractV01Path)),
		map[string]*regexp.Regexp{
			fungibleAddr.String(): ftAddressPlaceholder,
		},
	))
}

func NowwhereGenerateGetSupplyScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01GetSupplyPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereGenerateGetCollectionScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01GetCollectionPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractV01GenerateGetBrandCountScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01GetCollectionCountPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractV01GenerateGetBrandNameScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01GetBrandNamePath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractV01GenerateGetBrandIDsScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01GetBrandIDPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func GetSchema_CountScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01GetSchemaCountPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

// Template Script
func NowwhereGenerateGetTemplateCountScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01GetTemplateCountPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

// Drops Script
func NowwhereGenerateGetDropCountScript(fungibleAddr, NFTContractV01, nowwhereAddr flow.Address) []byte {
	return nowwhereContractReplaceAddressPlaceholders(
		string(readFile(NowwhereContractgetDropCountPath)),
		fungibleAddr.String(),
		NFTContractV01.String(),
		nowwhereAddr.String(),
	)
}

// Drops Script
func NowwhereGenerateGetDropIdsScript(fungibleAddr, NFTContractV01, nowwhereAddr flow.Address) []byte {
	return nowwhereContractReplaceAddressPlaceholders(
		string(readFile(NowwhereContractgetDropIdsPath)),
		fungibleAddr.String(),
		NFTContractV01.String(),
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
		string(readFile(NFTContractV01GetNFTAddressPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereGenerateGetNFTAddressCountScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01GetNFTAddressCountPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func loadNonFungibleToken() []byte {
	return nft_contracts.NonFungibleToken()
}

func NowwhereCreateGenerateCollectionScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01CreateCollectionPath)),
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
		string(readFile(NFTContractV01UpdateBrandPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereCreateGenerateSchemaScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01CreateSchemaPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereCreateGenerateTemplateScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01CreateTemplatePath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractV01SetupAccountScript(nonfungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01SetupAccountPath)),
		nonfungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereSetupAdminAccountScript(nonfungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01SetupAdminAccountPath)),
		nonfungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractV01AddAdminCapabilityScript(nonfungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01AddAdminCapabilityPath)),
		nonfungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractV01TransferNFTScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01TransferTokensPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NFTContractV01DestroyNFTScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01DestroyTokensPath)),
		fungibleAddr.String(),
		nowwhereAddr.String(),
	)
}

func NowwhereMintTokensScript(fungibleAddr, nowwhereAddr flow.Address) []byte {
	return nowwhereReplaceAddressPlaceholders(
		string(readFile(NFTContractV01MintTokensPath)),
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

func NFTContractV01CreateBrandTransaction(
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

func NFTContractV01UpdateBrandTransaction(
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

func NFTContractV01SetupAccount(
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
		SetScript(NFTContractV01SetupAccountScript(nonfungibleAddr, nowwhereAddr)).
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

func NFTContractV01TransferNFT(
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
		SetScript(NFTContractV01TransferNFTScript(nonfungibleAddr, nowwhereAddr)).
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

func NFTContractV01DestroyNFT(
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
		SetScript(NFTContractV01DestroyNFTScript(nonfungibleAddr, nowwhereAddr)).
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

func NFTContractV01SetupNewAdminAccount(
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

func NFTContractV01SetupAdminAccount(
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

func NFTContractV01AddAdminCapability(
	testing *testing.T,
	emulator *emulator.Blockchain,
	nonfungibleAddr,
	nowwhereAddr sdk.Address,
	userSigner crypto.Signer,
	shouldFail bool,
	adminAddress sdk.Address,
) {

	tx := flow.NewTransaction().
		SetScript(NFTContractV01AddAdminCapabilityScript(nonfungibleAddr, nowwhereAddr)).
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
