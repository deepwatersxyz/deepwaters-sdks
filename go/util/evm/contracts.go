package evm

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"deepwaters/go-examples/util"
)

type contracts struct {
	contracts map[string]Contract
}

func NewContracts() Contracts {
	return &contracts{
		contracts: make(map[string]Contract),
	}
}

func (c *contracts) GetKey(chainName string, address string) string {
	fullName := chainName + "-" + address
	fullName = strings.ToLower(fullName)
	return fullName
}

func (c *contracts) LoadFromConfigs(contractConfigs []ContractConfig, evm EVMConnector, lg *log.Logger) error {
	lge := lg.WithField("package", "micro.web3.evm")

	for _, cfg := range contractConfigs {

		if cfg.ChainID != evm.GetConfig().ChainID {
			return fmt.Errorf("inconsistent chainIDs")
		}

		contract, err := NewContractFromConfig(cfg, evm)
		if err != nil {
			return err
		}

		logFields := log.Fields{
			"contract": contract.GetName(),
			"address":  contract.GetAddressStr(true),
			"network":  evm.GetConfig().ChainName,
			"block":    util.Uint64Str(contract.GetDeployedBlock(), true),
		}

		lge.WithFields(logFields).Debug("successfully loaded contract from abi")

		if c.HasContract(evm.GetConfig().ChainName, contract.GetAddressStr(true)) {
			return fmt.Errorf("contract already configured: %+v %+v", contract.GetEVM().GetConfig().ChainName, contract.GetAddressStr(true))
		}

		if err := c.AddContract(contract); err != nil {
			return err
		}

		lge.WithFields(logFields).Debug("added contract")
	}

	return nil
}

func (c *contracts) AddContract(contract Contract) error {
	if c.HasContract(contract.GetChainName(), contract.GetAddressStr(true)) {
		return fmt.Errorf("contract already exists: %+v", c.GetKey(contract.GetChainName(), contract.GetAddressStr(true)))
	}

	fullName := c.GetKey(contract.GetChainName(), contract.GetAddressStr(true))

	c.contracts[fullName] = contract

	return nil
}

func (c *contracts) HasContract(chainName string, address string) bool {
	fullName := c.GetKey(chainName, address)
	_, ok := c.contracts[fullName]
	return ok
}

func (c *contracts) GetContract(chainName string, address string) Contract {
	fullName := c.GetKey(chainName, address)
	if contract, ok := c.contracts[fullName]; ok {
		return contract
	}
	return nil
}

func (c *contracts) ForEach(f func(c Contract) error) error {
	for _, c := range c.contracts {
		if err := f(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *contracts) Clone() Contracts {
	clone := NewContracts().(*contracts)

	for k, v := range c.contracts {
		clone.contracts[k] = v
	}

	return clone
}
