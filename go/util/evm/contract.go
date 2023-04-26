package evm

import (
	"context"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"deepwaters/go-examples/util"
)

type contract struct {
	abi           *abi.ABI
	address       common.Address
	name          string
	deployedBlock *uint64
	evm           Connector
}

func NewContract(abi *abi.ABI, addr, name string, deployedBlock *uint64, evm Connector) Contract {
	return &contract{
		name:          name,
		abi:           abi,
		address:       common.HexToAddress(addr),
		evm:           evm,
		deployedBlock: deployedBlock,
	}
}

func NewContractFromConfig(contractConfig ContractConfig, evm Connector) (Contract, error) {
	abiJSONBytes, err := ioutil.ReadFile(contractConfig.ABIFilePath)
	if err != nil {
		return nil, err
	}
	abi, err := abi.JSON(strings.NewReader(string(abiJSONBytes)))
	if err != nil {
		return nil, err
	}
	return NewContract(&abi, contractConfig.AddressHexStr, contractConfig.Name, contractConfig.DeployedBlock, evm), nil
}

func (c *contract) GetName() string {
	return c.name
}

func (c *contract) GetABI() *abi.ABI {
	return c.abi
}

func (c *contract) GetAddressBytes() common.Address {
	return c.address
}

func (c *contract) GetAddressStr(checkSum bool) string {
	str := c.address.Hex()
	if checkSum {
		return str
	}
	return strings.ToLower(str)
}

func (c *contract) GetConnector() Connector {
	return c.evm
}

func (c *contract) GetDeployedBlock() *uint64 {
	return c.deployedBlock
}

func (c *contract) Call(ctx context.Context, block *uint64, from *string, f string, arg ...interface{}) ([]byte, error) {
	if from == nil {
		return c.evm.Call(ctx, c.abi, c.address, block, nil, f, arg...)
	}
	sender := common.HexToAddress(*from)
	return c.evm.Call(ctx, c.abi, c.address, block, &sender, f, arg...)
}

func (c *contract) Send(mode SendMode, from Wallet, to *string, value *big.Int, nonce *uint64, f *string, args ...interface{}) (*string, *types.Receipt, *uint64, error) {
	hash, receipt, nonce, err := c.evm.Send(SendArgs{
		Mode:    mode,
		Wallet:  from,
		Address: common.HexToAddress(*to),
		ABI:     c.abi,
		Value:   value,
		Nonce:   nonce,
		F:       f,
		Args:    args,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	if hash == nil {
		// dry mode
		return nil, receipt, nil, nil
	}
	return util.StrP(hash.Hex()), receipt, nonce, nil
}
