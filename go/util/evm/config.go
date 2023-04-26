package evm

import (
	"fmt"
	"math/big"
)

type NativeCurrency struct {
	Decimals uint8
	Name     string
	Symbol   string
}

type ConnectorConfig struct {
	ChainID   int
	ChainName string
	URL       string
	//PerSecRateLimit int
	//Sender          SenderConfig
	NativeCurrency NativeCurrency
	//ContractConfigs []ContractConfig
	MaxGasPriceWei         *big.Int
	AddEstimatedGasPercent int
}

//type SenderConfig struct {
//	MaxGasPriceWei         *big.Int
//	AddEstimatedGasPercent int
//}

type ContractConfig struct {
	ABIFilePath   string
	ChainID       int
	ChainName     string
	AddressHexStr string
	Name          string
	Description   *string
	DeployedBlock *uint64
}

func NewConnectorConfig(chainID int, chainName string, rpcURL string /*perSecRateLimit int,*/, maxGasPriceWeiStr string, addEstimatedGasPercent int, nativeCurrency NativeCurrency /*, contractConfigs []ContractConfig*/) (*ConnectorConfig, error) {
	maxGasPriceWei, success := big.NewInt(0).SetString(maxGasPriceWeiStr, 10)
	if !success {
		return nil, fmt.Errorf("unable to parse maxGasPriceWeiStr: %s", maxGasPriceWeiStr)
	}

	cfg := &ConnectorConfig{
		ChainID:   chainID,
		ChainName: chainName,
		URL:       rpcURL,
		//PerSecRateLimit: perSecRateLimit,
		NativeCurrency: nativeCurrency,
		//ContractConfigs: contractConfigs,
		MaxGasPriceWei:         maxGasPriceWei,
		AddEstimatedGasPercent: addEstimatedGasPercent,
		//		Sender: SenderConfig{
		//			MaxGasPriceWei:         maxGasPriceWei,
		//			AddEstimatedGasPercent: addEstimatedGasPercent,
		//		},
	}

	return cfg, nil
}
