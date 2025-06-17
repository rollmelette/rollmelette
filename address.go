// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import "github.com/ethereum/go-ethereum/common"

// AddressBook contains the addresses of the rollups contracts.
type AddressBook struct {
	ApplicationFactory           common.Address
	AuthorityFactory             common.Address
	ERC1155BatchPortal           common.Address
	ERC1155SinglePortal          common.Address
	ERC20Portal                  common.Address
	ERC721Portal                 common.Address
	EtherPortal                  common.Address
	InputBox                     common.Address
	SelfHostedApplicationFactory common.Address
	TestToken                    common.Address
	TestNFT                      common.Address
	TestMultiToken               common.Address
}

// NewAddressBook returns the contract addresses for mainnet and devnet.
func NewAddressBook() AddressBook {
	return AddressBook{
		ApplicationFactory:           common.HexToAddress("0xc7006f70875BaDe89032001262A846D3Ee160051"),
		AuthorityFactory:             common.HexToAddress("0xC7003566dD09Aa0fC0Ce201aC2769aFAe3BF0051"),
		ERC1155BatchPortal:           common.HexToAddress("0xc700A2e5531E720a2434433b6ccf4c0eA2400051"),
		ERC1155SinglePortal:          common.HexToAddress("0xc700A261279aFC6F755A3a67D86ae43E2eBD0051"),
		ERC20Portal:                  common.HexToAddress("0xc700D6aDd016eECd59d989C028214Eaa0fCC0051"),
		ERC721Portal:                 common.HexToAddress("0xc700d52F5290e978e9CAe7D1E092935263b60051"),
		EtherPortal:                  common.HexToAddress("0xc70076a466789B595b50959cdc261227F0D70051"),
		InputBox:                     common.HexToAddress("0xc70074BDD26d8cF983Ca6A5b89b8db52D5850051"),
		SelfHostedApplicationFactory: common.HexToAddress("0xc700285Ab555eeB5201BC00CFD4b2CC8DED90051"),
		TestToken:                    common.HexToAddress("0xFBdB734EF6a23aD76863CbA6f10d0C5CBBD8342C"),
		TestNFT:                      common.HexToAddress("0xBa46623aD94AB45850c4ecbA9555D26328917c3B"),
		TestMultiToken:               common.HexToAddress("0xDC6d64971B77a47fB3E3c6c409D4A05468C398D2"),
	}
}
