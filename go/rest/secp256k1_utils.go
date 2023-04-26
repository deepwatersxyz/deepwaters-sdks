package rest

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"encoding/hex"

	"deepwaters/go-examples/util"
)

func CheckPrivateKeyHexStr(privateKeyHexStr string) error {
	if privateKeyHexStr[:2] != "0x" {
		return fmt.Errorf("private key must be a hex string beginning with '0x'")
	}
	if len(privateKeyHexStr) != 66 {
		return fmt.Errorf("private key must have length: 66")
	}
	return nil
}

func GetKeccak256Hash(str string) *common.Hash {

	rawData := []byte(str)
	sigHashBytes := crypto.Keccak256(rawData)
	hash := common.BytesToHash(sigHashBytes) // these stringify with leading '0x'

	return &hash
}

func GetECDSASignatureHexStr(hash common.Hash, privateKey *ecdsa.PrivateKey) (*string, error) {

	sig, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}
	result := util.StrP("0x" + hex.EncodeToString(sig))
	return result, nil
}

func HashAndSignFromPrivateKey(str string, privateKey *ecdsa.PrivateKey) (*string, error) {

	hash := GetKeccak256Hash(str)
	return GetECDSASignatureHexStr(*hash, privateKey)
}

func HashAndSignFromPrivateKeyHexStr(str string, privateKeyHexStr string) (*string, error) {

	if err := CheckPrivateKeyHexStr(privateKeyHexStr); err != nil {
		return nil, err
	}

	prv, err := crypto.HexToECDSA(privateKeyHexStr[2:]) // this does further checks on the private key hex string
	if err != nil {
		return nil, err
	}

	signature, err := HashAndSignFromPrivateKey(str, prv)
	return signature, err
}
