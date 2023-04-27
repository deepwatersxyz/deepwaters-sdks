package util

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func GetPrefixedHash(hash common.Hash) common.Hash {
	return crypto.Keccak256Hash(
		[]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%v", len(hash))),
		hash.Bytes(),
	)
}

func GetECDSASignature(hash common.Hash, privateKey *ecdsa.PrivateKey) (*string, error) {
	sig, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}
	sig[len(sig)-1] += 27
	return StrP("0x" + hex.EncodeToString(sig)), nil
}

func ECRecover(hash common.Hash, signature string, allowNoRecoveryIDOffset bool) (common.Address, error) {
	// signature
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return common.Address{}, err
	}
	if len(sig) != crypto.SignatureLength {
		return common.Address{}, fmt.Errorf("signature must be %d bytes long", crypto.SignatureLength)
	}
	if allowNoRecoveryIDOffset && sig[crypto.RecoveryIDOffset] < 2 {
		sig[crypto.RecoveryIDOffset] += 27
	}
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[crypto.RecoveryIDOffset] -= 27 // Transform yellow paper V from 27/28 to 0/1

	// ecrecover
	pubKeyBytes, err := crypto.Ecrecover(hash.Bytes(), sig)
	if err != nil {
		return common.Address{}, err
	}

	// address
	pubKey, err := crypto.UnmarshalPubkey(pubKeyBytes)
	if err != nil {
		return common.Address{}, err
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}
