package test

import "regexp"

var (
	ftAddressPlaceholder = regexp.MustCompile(`"[^"\s].*/FungibleToken.cdc"`)

	NFTContractV01AddressPlaceHolder = regexp.MustCompile(`"[^"\s].*/NFTContractV01.cdc"`)
	nftAddressPlaceholder            = regexp.MustCompile(`"[^"\s].*/NonFungibleToken.cdc"`)
	NowwherePlaceholder              = regexp.MustCompile(`"[^"\s].*/NowwhereContract.cdc"`)
)
