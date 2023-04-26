package evm

import "deepwaters/go-examples/util"

const (
	positionManagerABIPath = "../abi/PositionManager.json"
	erc20ABIPath = "../abi/ERC20.json"
)

func GetFujiAVAX() NativeCurrency {
	return NativeCurrency{
		Decimals: 18,
		Name:     "Fuji AVAX",
		Symbol:   "AVAX",
	}
}

func GetTestnetProdDomainName() string {
	return "testnet.api.deepwaters.xyz"
}

func GetFujiTestnetProdPositionManagerConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   positionManagerABIPath,
		ChainID:       43113,
		ChainName:     "fuji",
		AddressHexStr: "0x5F2C6309389CA1524402E928dC9c377757B0F947",
		Name:          "PositionManager",
		Description:   util.StrP("Position Manager"),
		DeployedBlock: util.Uint64P(16125650),
	}
}

func GetFujiTestnetProdWAVAXConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       43113,
		ChainName:     "fuji",
		AddressHexStr: "0xd4743A7B6cCeAa5d6EBAb19013012F5D1Fc779CB",
		Name:          "WAVAX",
		Description:   util.StrP("Wrapped AVAX"),
		DeployedBlock: util.Uint64P(16125690),
	}
}

func GetFujiTestnetProdUSDCConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       43113,
		ChainName:     "fuji",
		AddressHexStr: "0xA4C930EbD593197226CEc2Cbdc6927bcF405338C",
		Name:          "USDC",
		Description:   util.StrP("USD Coin"),
		DeployedBlock: util.Uint64P(16125701),
	}
}

func GetGoerliTestnetProdPositionManagerConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   positionManagerABIPath,
		ChainID:       5,
		ChainName:     "goerli",
		AddressHexStr: "0xbaf1fb2CBDfaE48C0169A837EC6E39294cE9Eb15",
		Name:          "PositionManager",
		Description:   util.StrP("Position Manager"),
		DeployedBlock: util.Uint64P(8005362),
	}
}

func GetGoerliTestnetProdWETHConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       5,
		ChainName:     "goerli",
		AddressHexStr: "0x952269331420120dB1583F3628175a0EBbd56113",
		Name:          "WETH",
		Description:   util.StrP("Wrapped Ether"),
		DeployedBlock: util.Uint64P(8005381),
	}
}
func GetGoerliTestnetProdWBTCConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       5,
		ChainName:     "goerli",
		AddressHexStr: "0x4696397b41F0d6449b20320b5eE320DfeF3Bd2B5",
		Name:          "WBTC",
		Description:   util.StrP("Wrapped Bitcoin"),
		DeployedBlock: util.Uint64P(8005386),
	}
}

func GetGoerliTestnetProdUSDCConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       5,
		ChainName:     "goerli",
		AddressHexStr: "0x08235f416aA7Af10eb02b04F167DeFf0B0c84Ccb",
		Name:          "USDC",
		Description:   util.StrP("USD Coin"),
		DeployedBlock: util.Uint64P(8005392),
	}
}

func GetMumbaiTestnetProdPositionManagerConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   positionManagerABIPath,
		ChainID:       80001,
		ChainName:     "mumbai",
		AddressHexStr: "0xa77E7a2e3cAAE4de7D08B4aae7e4BAF93CbF29F0",
		Name:          "PositionManager",
		Description:   util.StrP("Position Manager"),
		DeployedBlock: util.Uint64P(29297109),
	}
}

func GetMumbaiTestnetProdWMATICConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       80001,
		ChainName:     "mumbai",
		AddressHexStr: "0x9F324abb8a8744a18030deB6b2888dB6B7B6841D",
		Name:          "WBNB",
		Description:   util.StrP("Wrapped MATIC"),
		DeployedBlock: util.Uint64P(29297236),
	}
}

func GetMumbaiTestnetProdUSDCConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       80001,
		ChainName:     "mumbai",
		AddressHexStr: "0x08Bc45D8b4b5d9844dB5B1686aB36d0472328C6d",
		Name:          "USDC",
		Description:   util.StrP("USD Coin"),
		DeployedBlock: util.Uint64P(29297242),
	}
}

func GetRialtoTestnetProdPositionManagerConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   positionManagerABIPath,
		ChainID:       97,
		ChainName:     "rialto",
		AddressHexStr: "0xE8AbdC59A0DB16468A874DD10420E3c564c47104",
		Name:          "PositionManager",
		Description:   util.StrP("Position Manager"),
		DeployedBlock: util.Uint64P(24844966),
	}
}

func GetRialtoTestnetProdWBNBConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       97,
		ChainName:     "rialto",
		AddressHexStr: "0x4Ca5898e82260586f6CBa2A25848816AA9d8eE03",
		Name:          "WBNB",
		Description:   util.StrP("Wrapped BNB"),
		DeployedBlock: util.Uint64P(24845210),
	}
}

func GetRialtoTestnetProdUSDConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       97,
		ChainName:     "rialto",
		AddressHexStr: "0x3Fc54ADd69955724169E9aB22D59152320811327",
		Name:          "USDC",
		Description:   util.StrP("USD Coin"),
		DeployedBlock: util.Uint64P(24845232),
	}
}
