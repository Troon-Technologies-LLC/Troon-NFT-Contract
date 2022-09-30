package test

import "regexp"

var (
	ftAddressPlaceholder = regexp.MustCompile(`"[^"\s].*/FungibleToken.cdc"`)

	TroonAtomicStandardAddressPlaceHolder = regexp.MustCompile(`"[^"\s].*/TroonAtomicStandard.cdc"`)
	nftAddressPlaceholder         = regexp.MustCompile(`"[^"\s].*/NonFungibleToken.cdc"`)
	NowwherePlaceholder           = regexp.MustCompile(`"[^"\s].*/NowwhereContract.cdc"`)
)
