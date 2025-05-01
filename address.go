// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package rollmelette

import "github.com/ethereum/go-ethereum/common"

// AddressBook contains the addresses of the rollups contracts.
type AddressBook struct {
	ApplicationFactory           common.Address
	AuthorityFactory             common.Address
	EntryPointV06                common.Address
	EntryPointV07                common.Address
	ERC1155BatchPortal           common.Address
	ERC1155SinglePortal          common.Address
	ERC20Portal                  common.Address
	ERC721Portal                 common.Address
	EtherPortal                  common.Address
	InputBox                     common.Address
	LightAccountFactory          common.Address
	SelfHostedApplicationFactory common.Address
	SimpleAccountFactory         common.Address
	SmartAccountFactory          common.Address
	KernelFactoryV2              common.Address
	KernelFactoryV3              common.Address
	KernelFactoryV3_1            common.Address
	TestToken                    common.Address
	TestNFT                      common.Address
	TestMultiToken               common.Address
	VerifyingPaymasterV06        common.Address
	VerifyingPaymasterV07        common.Address
}

// NewAddressBook returns the contract addresses for mainnet and devnet.
func NewAddressBook() AddressBook {
	return AddressBook{
		ApplicationFactory:           common.HexToAddress("0x2210ad1d9B0bD2D470c2bfA4814ab6253BC421A0"),
		AuthorityFactory:             common.HexToAddress("0x451f57Ca716046D114Ab9ff23269a2F9F4a1bdaF"),
		EntryPointV06:                common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789"),
		EntryPointV07:                common.HexToAddress("0x0000000071727De22E5E9d8BAf0edAc6f37da032"),
		ERC1155BatchPortal:           common.HexToAddress("0xBc70d79F916A6d48aB0b8F03AC58f89742dEDA34"),
		ERC1155SinglePortal:          common.HexToAddress("0xB778147D50219544F113A55DE1d8de626f0cC1bB"),
		ERC20Portal:                  common.HexToAddress("0x05355c2F9bA566c06199DEb17212c3B78C1A3C31"),
		ERC721Portal:                 common.HexToAddress("0x0F5A20d3729c44FedabBb560b3D633dc1c246DDe"),
		EtherPortal:                  common.HexToAddress("0xd31aD6613bDaA139E7D12B2428C0Dd00fdBF8aDa"),
		InputBox:                     common.HexToAddress("0xB6b39Fb3dD926A9e3FBc7A129540eEbeA3016a6c"),
		LightAccountFactory:          common.HexToAddress("0x00004EC70002a32400f8ae005A26081065620D20"),
		SelfHostedApplicationFactory: common.HexToAddress("0x4a409e1CaB9229711C4e1f68625DdbC75809e721"),
		SimpleAccountFactory:         common.HexToAddress("0x9406Cc6185a346906296840746125a0E44976454"),
		SmartAccountFactory:          common.HexToAddress("0x000000a56Aaca3e9a4C479ea6b6CD0DbcB6634F5"),
		KernelFactoryV2:              common.HexToAddress("0x5de4839a76cf55d0c90e2061ef4386d962E15ae3"),
		KernelFactoryV3:              common.HexToAddress("0x6723b44Abeec4E71eBE3232BD5B455805baDD22f"),
		KernelFactoryV3_1:            common.HexToAddress("0xaac5D4240AF87249B3f71BC8E4A2cae074A3E419"),
		TestToken:                    common.HexToAddress("0xFBdB734EF6a23aD76863CbA6f10d0C5CBBD8342C"),
		TestNFT:                      common.HexToAddress("0xBa46623aD94AB45850c4ecbA9555D26328917c3B"),
		TestMultiToken:               common.HexToAddress("0xDC6d64971B77a47fB3E3c6c409D4A05468C398D2"),
		VerifyingPaymasterV06:        common.HexToAddress("0x28ec0633192d0cBd9E1156CE05D5FdACAcB93947"),
		VerifyingPaymasterV07:        common.HexToAddress("0xc5c97885C67F7361aBAfD2B95067a5bBDa603608"),
	}
}
