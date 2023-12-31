// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import "github.com/ethereum/go-ethereum/common"

// AddressBook contains the addresses of the rollups contracts.
type AddressBook struct {
	CartesiAppFactory   common.Address
	AppAddressRelay     common.Address
	ERC1155BatchPortal  common.Address
	ERC1155SinglePortal common.Address
	ERC20Portal         common.Address
	ERC721Portal        common.Address
	EtherPortal         common.Address
	InputBox            common.Address
}

// NewAddressBook returns the contract addresses for mainnet and devnet.
func NewAddressBook() AddressBook {
	return AddressBook{
		CartesiAppFactory:   common.HexToAddress("0x7122cd1221C20892234186facfE8615e6743Ab02"),
		AppAddressRelay:     common.HexToAddress("0xF5DE34d6BbC0446E2a45719E718efEbaaE179daE"),
		ERC1155BatchPortal:  common.HexToAddress("0xedB53860A6B52bbb7561Ad596416ee9965B055Aa"),
		ERC1155SinglePortal: common.HexToAddress("0x7CFB0193Ca87eB6e48056885E026552c3A941FC4"),
		ERC20Portal:         common.HexToAddress("0x9C21AEb2093C32DDbC53eEF24B873BDCd1aDa1DB"),
		ERC721Portal:        common.HexToAddress("0x237F8DD094C0e47f4236f12b4Fa01d6Dae89fb87"),
		EtherPortal:         common.HexToAddress("0xFfdbe43d4c855BF7e0f105c400A50857f53AB044"),
		InputBox:            common.HexToAddress("0x59b22D57D4f067708AB0c00552767405926dc768"),
	}
}
