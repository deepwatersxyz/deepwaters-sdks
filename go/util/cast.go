package util

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

func BigUint64(i *uint64) *big.Int {
	if i == nil {
		return nil
	}
	b := new(big.Int)
	b.SetUint64(*i)
	return b
}

func BigIntFromBytes(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}

func BigIntFromStr(s string) *big.Int {
	b, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil
	}
	return b
}

func BigIntFromStrP(s *string) *big.Int {
	if s == nil {
		return nil
	}
	return BigIntFromStr(*s)
}

func BigIntStr(i *big.Int, pretty bool) string {
	if i == nil {
		return "<nil>"
	}
	if pretty {
		return PrettyNum(i)
	}
	return i.String()
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

func Int64P(i int64) *int64 {
	return &i
}

func IntP(i int) *int {
	return &i
}

func UIntP(i uint) *uint {
	return &i
}

func Float32P(f float32) *float32 {
	return &f
}

func Float64P(f float64) *float64 {
	return &f
}

func BoolStr(b *bool) string {
	if b == nil {
		return "<nil>"
	}
	if *b {
		return "true"
	}
	return "false"
}

func BoolP(b bool) *bool {
	return &b
}

func IntStr(i *int, pretty bool) string {
	if i == nil {
		return "<nil>"
	}
	if pretty {
		return PrettyNum(*i)
	}
	return fmt.Sprintf("%+v", *i)
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

func UintStr(i *uint, pretty bool) string {
	if i == nil {
		return "<nil>"
	}
	if pretty {
		return PrettyNum(*i)
	}
	return fmt.Sprintf("%+v", *i)
}

func TimeStr(t *time.Time) string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%+v", *t)
}

func DurationStr(d *time.Duration) string {
	if d == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%+v", *d)
}

func DurationP(d time.Duration) *time.Duration {
	return &d
}

func TimeP(t time.Time) *time.Time {
	return &t
}

func ConvertStruct(from, to interface{}) error {
	b, err := json.Marshal(from)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, to)
	if err != nil {
		return err
	}
	return nil
}

func UUIDP(id uuid.UUID) *uuid.UUID {
	return &id
}

func UUIDStr(id *uuid.UUID) string {
	if id == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%+v", *id)
}
