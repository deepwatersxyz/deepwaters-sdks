package evm

import "deepwaters/go-examples/util"

const (
	positionManagerABIPath = "../abi/PositionManager.json"
	erc20ABIPath = "../abi/ERC20.json"
)

// testnet

func GetFujiAVAX() NativeCurrency {
	return NativeCurrency{
		Decimals: 18,
		Name:     "Fuji AVAX",
		Symbol:   "AVAX",
	}
}

func GetGoerliETH() NativeCurrency {
	return NativeCurrency{
		Decimals: 18,
		Name:     "Goerli Ether",
		Symbol:   "ETH",
	}
}

func GetMumbaiMATIC() NativeCurrency {
	return NativeCurrency{
		Decimals: 18,
		Name:     "Mumbai MATIC",
		Symbol:   "MATIC",
	}
}

func GetRialtoBNB() NativeCurrency {
	return NativeCurrency{
		Decimals: 18,
		Name:     "Rialto BNB",
		Symbol:   "BNB",
	}
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
		Name:          "WMATIC",
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

// mainnet

func GetAVAX() NativeCurrency {
	return NativeCurrency{
		Decimals: 18,
		Name:     "AVAX",
		Symbol:   "AVAX",
	}
}

func GetETH() NativeCurrency {
	return NativeCurrency{
		Decimals: 18,
		Name:     "Ether",
		Symbol:   "ETH",
	}
}

func GetMATIC() NativeCurrency {
	return NativeCurrency{
		Decimals: 18,
		Name:     "MATIC",
		Symbol:   "MATIC",
	}
}

func GetBNB() NativeCurrency {
	return NativeCurrency{
		Decimals: 18,
		Name:     "BNB",
		Symbol:   "BNB",
	}
}

func GetAvalancheMainnetPositionManagerConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   positionManagerABIPath,
		ChainID:       43114,
		ChainName:     "avalance c-chain",
		AddressHexStr: "0x32bb0a6CeEfcE9cC222b54f8159f56aF035D2aBA",
		Name:          "PositionManager",
		Description:   util.StrP("Position Manager"),
		DeployedBlock: util.Uint64P(26134721),
	}
}

func GetAvalancheMainnetWAVAXConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       43114,
		ChainName:     "avalanche c-chain",
		AddressHexStr: "0xb31f66aa3c1e785363f0875a1b74e27b85fd66c7",
		Name:          "WAVAX",
		Description:   util.StrP("Wrapped AVAX"),
		DeployedBlock: util.Uint64P(820),
	}
}

func GetEthereumMainnetPositionManagerConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   positionManagerABIPath,
		ChainID:       1,
		ChainName:     "ethereum mainnet",
		AddressHexStr: "0xC86289E5889eF21bA60dC6D5F1c487EC84FbED61",
		Name:          "PositionManager",
		Description:   util.StrP("Position Manager"),
		DeployedBlock: util.Uint64P(16607968),
	}
}

func GetEthereumMainnetWETHConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       1,
		ChainName:     "ethereum mainnet",
		AddressHexStr: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
		Name:          "WETH",
		Description:   util.StrP("Wrapped Ether"),
		DeployedBlock: util.Uint64P(4719568),
	}
}
func GetEthereumMainnetWBTCConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       1,
		ChainName:     "ethereum mainnet",
		AddressHexStr: "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
		Name:          "WBTC",
		Description:   util.StrP("Wrapped Bitcoin"),
		DeployedBlock: util.Uint64P(6766284),
	}
}

func GetEthereumMainnetUSDCConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       1,
		ChainName:     "ethereum mainnet",
		AddressHexStr: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		Name:          "USDC",
		Description:   util.StrP("USD Coin"),
		DeployedBlock: util.Uint64P(6082465),
	}
}

func GetPolygonMainnetPositionManagerConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   positionManagerABIPath,
		ChainID:       137,
		ChainName:     "polygon mainnet",
		AddressHexStr: "0x32A717699C68ca555D7E1212F8822c980A8d48A3",
		Name:          "PositionManager",
		Description:   util.StrP("Position Manager"),
		DeployedBlock: util.Uint64P(39185355),
	}
}

func GetPolygonMainnetWMATICConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       137,
		ChainName:     "polygon mainnet",
		AddressHexStr: "0x0d500b1d8e8ef31e21c99d1db9a6444d3adf1270",
		Name:          "WMATIC",
		Description:   util.StrP("Wrapped MATIC"),
		DeployedBlock: util.Uint64P(4931456),
	}
}

func GetBSCMainnetPositionManagerConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   positionManagerABIPath,
		ChainID:       56,
		ChainName:     "binance smart chain mainnet",
		AddressHexStr: "0x5612d989868fe9f3A2b2923D66C577d1d8bE2A8D",
		Name:          "PositionManager",
		Description:   util.StrP("Position Manager"),
		DeployedBlock: util.Uint64P(25584258),
	}
}

func GetBSCWBNBConfig() ContractConfig {
	return ContractConfig{
		ABIFilePath:   erc20ABIPath,
		ChainID:       56,
		ChainName:     "binance smart chain mainnet",
		AddressHexStr: "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c",
		Name:          "WBNB",
		Description:   util.StrP("Wrapped BNB"),
		DeployedBlock: util.Uint64P(149268),
	}
}
