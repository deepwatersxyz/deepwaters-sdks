package util

import (
	"fmt"
	"math/big"

	"database/sql/driver"
)

type BigInt struct {
	Int *big.Int `json:"int"`
}

func NewNullBigInt() *BigInt {
	return &BigInt{Int: big.NewInt(0)}
}

func (dst *BigInt) Scan(src interface{}) error {
	if src == nil {
		*dst = BigInt{}
		return nil
	}

	switch t := src.(type) {
	case bool:
		if src.(bool) {
			dst.Int = big.NewInt(1)
		} else {
			dst.Int = big.NewInt(0)
		}
		return nil
	case []byte:
		return dst.Scan(string(src.([]byte)))
	case string:
		dst.Int = &big.Int{}
		r, _ := dst.Int.SetString(src.(string), 10)
		if r == nil {
			dst.Int = nil
			return fmt.Errorf("unable to scan BigInt: %T, %+v", t, src)
		}
		return nil
	default:
		return fmt.Errorf("unable to scan BigInt: %T, %+v", t, src)
	}
}

func (src BigInt) Value() (driver.Value, error) {
	if src.Int == nil {
		return nil, nil
	}
	return src.Int.String(), nil
}

func (b *BigInt) AddInt64(i int64) {
	if b.Int == nil {
		b.Int = &big.Int{}
	}
	b.Int.Add(b.Int, big.NewInt(i))
}

func NewBigIntFromStr(s string) (*BigInt, error) {
	i, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("unable to convert to big.Int: %+v", s)
	}
	return &BigInt{Int: i}, nil
}
