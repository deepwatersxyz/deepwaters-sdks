package evm

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/params"
)

func WeiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}

func EtherToWei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", eth), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

func ParseToBigFloat(value string) (*big.Float, error) {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	_, err := fmt.Sscan(value, f)
	return f, err
}

func EtherStrToWei(value string) (*big.Int, error) {
	f, err := ParseToBigFloat(value)
	if err != nil {
		return nil, err
	}

	wei := EtherToWei(f)
	return wei, nil
}

var (
	big10      = big.NewInt(10)
	zeroString = "00000000000000000000000000000000000000000000000000000000000000000000000000000000"
)

func ConvertAmountDecimals(amount *big.Int, fromDecimals uint8, toDecimals uint8) *big.Int {
	converted := new(big.Int).Set(amount)

	for ; fromDecimals < toDecimals; fromDecimals += 1 {
		converted = converted.Mul(converted, big10)
	}
	for ; fromDecimals > toDecimals; fromDecimals -= 1 {
		converted = converted.Div(converted, big10)
	}

	return converted
}

func AmountFromParametrizedText(amount string, fromDecimals uint8, toDecimals uint8) (*big.Int, error) {
	// get floating point
	index := strings.LastIndex(amount, ".")
	if index == -1 {
		amount = amount + zeroString[:fromDecimals]
	} else {
		numDecimals := uint8(len(amount) - index - 1)
		if numDecimals > fromDecimals {
			return nil, fmt.Errorf("too many decimals")
		}

		// remove "." and add pending 0 at the end
		amount = amount[:index] + amount[index+1:] + zeroString[:fromDecimals-numDecimals]
	}

	amountBig, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse amount")
	}

	return ConvertAmountDecimals(amountBig, fromDecimals, toDecimals), nil
}

func AmountFromText(amount string, toDecimals uint8) (*big.Int, error) {
	var fromDecimals uint8
	index := strings.LastIndex(amount, ".")
	if index == -1 {
		fromDecimals = 0
	} else {
		fromDecimals = uint8(len(amount) - index - 1)
	}

	return AmountFromParametrizedText(amount, fromDecimals, toDecimals)
}

func AmountToText(amount *big.Int, fromDecimals uint8, toDecimals uint8) string {
	converted := ConvertAmountDecimals(amount, fromDecimals, toDecimals)

	stringified := converted.String()

	// Pad with 0 if length is too small. Ensure at least uiDecimals + 1 length
	if uint8(len(stringified)) < toDecimals+1 {
		stringified = zeroString[:(toDecimals+1)-uint8(len(stringified))] + stringified
	}

	length := uint8(len(stringified))
	return stringified[:length-toDecimals] + "." + stringified[length-toDecimals:]
}
