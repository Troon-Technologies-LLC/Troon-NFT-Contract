package test

import "regexp"

var (
	ftAddressPlaceholder = regexp.MustCompile(`"[^"\s].*/FungibleToken.cdc"`)

	NFTContractAddressPlaceHolder = regexp.MustCompile(`"[^"\s].*/NFTContract.cdc"`)
	nftAddressPlaceholder         = regexp.MustCompile(`"[^"\s].*/NonFungibleToken.cdc"`)
	NowwherePlaceholder           = regexp.MustCompile(`"[^"\s].*/NowwhereContract.cdc"`)
)
