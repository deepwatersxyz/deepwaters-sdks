package util

import (
	"fmt"
	"math/big"
)

func BigUint64(i *uint64) *big.Int {
	if i == nil {
		return nil
	}
	b := new(big.Int)
	b.SetUint64(*i)
	return b
}

func BigIntFromStr(s string) *big.Int {
	b, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil
	}
	return b
}

func Str(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}

func StrP(s string) *string {
	return &s
}

func Uint64P(i uint64) *uint64 {
	return &i
}

func Uint64Str(i *uint64, pretty bool) string {
	if i == nil {
		return "<nil>"
	}
	if pretty {
		return PrettyNum(*i)
	}
	return fmt.Sprintf("%+v", *i)
}
