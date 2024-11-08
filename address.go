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
	QuorumFactory                common.Address
	SafeERC20Transfer            common.Address
	SelfHostedApplicationFactory common.Address
}

// NewAddressBook returns the contract addresses for mainnet and devnet.
func NewAddressBook() AddressBook {
	return AddressBook{
		ApplicationFactory:           common.HexToAddress("0xA1DA32BF664109D62208a1cb0d69aACc6a484873"),
		AuthorityFactory:             common.HexToAddress("0xbDC5D42771A4Ae55eC7670AAdD2458D1d9C7C8A8"),
		ERC1155BatchPortal:           common.HexToAddress("0x4a218D331C0933d7E3EB496ac901669f28D94981"),
		ERC1155SinglePortal:          common.HexToAddress("0x2f0D587DD6EcF67d25C558f2e9c3839c579e5e38"),
		ERC20Portal:                  common.HexToAddress("0xB0e28881FF7ee9CD5B1229d570540d74bce23D39"),
		ERC721Portal:                 common.HexToAddress("0x874b3245ead7474Cb9f3b83cD1446dC522f6bd36"),
		EtherPortal:                  common.HexToAddress("0xfa2292f6D85ea4e629B156A4f99219e30D12EE17"),
		InputBox:                     common.HexToAddress("0x593E5BCf894D6829Dd26D0810DA7F064406aebB6"),
		QuorumFactory:                common.HexToAddress("0x68C3d53a095f66A215a8bEe096Cd3Ba4fFB7bAb3"),
		SafeERC20Transfer:            common.HexToAddress("0x817b126F242B5F184Fa685b4f2F91DC99D8115F9"),
		SelfHostedApplicationFactory: common.HexToAddress("0x0678FAA399F0193Fb9212BE41590316D275b1392"),
	}
}
