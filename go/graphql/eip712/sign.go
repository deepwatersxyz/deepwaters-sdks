package eip712

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

var (
	DomainName         = "EIP712Domain"
	DomainContractType = []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	}
	DomainAMType = []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
	}
)

func StrP(s string) *string {
	return &s
}

type Compatible interface {
	ToTypedData() apitypes.TypedData
}

// GetHash computes the hash of an EIP712 formatted data, according to the EIP712 specification
func GetHash(model Compatible) (*common.Hash, error) {
	typedData := model.ToTypedData()

	domainSeparator, err := typedData.HashStruct(DomainName, typedData.Domain.Map())
	if err != nil {
		return nil, err
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	sigHashBytes := crypto.Keccak256(rawData)
	hash := common.BytesToHash(sigHashBytes)
	return &hash, nil
}

func GetECDSASignature(hash common.Hash, privateKey *ecdsa.PrivateKey) (*string, error) {
	sig, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}
	sig[len(sig)-1] += 27
	return StrP("0x" + hex.EncodeToString(sig)), nil
}

// Sign the hashed EIP712 formatted model with the provided privateKey
func Sign(model Compatible, privateKey *ecdsa.PrivateKey) (*string, error) {
	hash, err := GetHash(model)
	if err != nil {
		return nil, err
	}

	// Do not call evm.wallet.GetECDSASignature() because it adds the "\x19Ethereum Signed Message:\n" prefix
	// EIP712 takes care of encoding when hashing. https://eips.ethereum.org/EIPS/eip-712#specification
	return GetECDSASignature(*hash, privateKey)
}

// SignFromHexVerbose the hashed EIP712 formatted model with the provided privateKey (in hex format), also gives address
func SignFromHexVerbose(model Compatible, privateKey string) (signature *string, address *string, err error) {
	prv, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, nil, err
	}
	addr := crypto.PubkeyToAddress(prv.PublicKey).String()
	address = &addr
	signature, err = Sign(model, prv)
	return
}

// SignFromHex the hashed EIP712 formatted model with the provided privateKey (in hex format)
func SignFromHex(model Compatible, privateKey string) (signature *string, err error) {
	prv, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	signature, err = Sign(model, prv)
	return
}
