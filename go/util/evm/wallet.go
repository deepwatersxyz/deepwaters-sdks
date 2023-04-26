package evm

import (
	"crypto/ecdsa"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"deepwaters/go-examples/util"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
}

func NewWalletFromHexStr(privateKeyHexStr string) (Wallet, error) {

	key, err := crypto.HexToECDSA(privateKeyHexStr[2:])
	if err != nil {
		return nil, err
	}
	return &wallet{
		privateKey: key,
	}, nil
}

func NewWalletFromFilePath(filePath string) (Wallet, error) {
	privateKeyBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	privateKeyHexStr := strings.Trim(string(privateKeyBytes), "\n")

	key, err := crypto.HexToECDSA(privateKeyHexStr[2:])
	if err != nil {
		return nil, err
	}
	return &wallet{
		privateKey: key,
	}, nil
}

func (w *wallet) GetPrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *wallet) GetAddressBytes() common.Address {
	publicKey := w.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("publicKey is not of type *ecdsa.PublicKey")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA)
}

func (w *wallet) GetAddressStr(checksum bool) string {
	addr := w.GetAddressBytes().Hex()
	if !checksum {
		return strings.ToLower(addr)
	}
	return addr
}

func (w *wallet) GetECDSASignature(hash common.Hash) (*string, error) {
	return util.GetECDSASignature(util.GetPrefixedHash(hash), w.privateKey)
}
