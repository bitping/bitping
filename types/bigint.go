package types

import (
	"fmt"
	"math/big"
)

type BigInt big.Int

func NewBigInt(i *big.Int) *BigInt {
	if i == nil {
		i = new(big.Int).SetInt64(0)
	}
	bi := BigInt(*i)
	return &bi
}

func BigIntFromInt(i int64) *BigInt {
	bi := new(big.Int).SetInt64(i)
	return NewBigInt(bi)
}

func BigIntFromString(s string) (*BigInt, bool) {
	i, worked := new(big.Int).SetString(s, 10)

	if i != nil {
		bi := BigInt(*i)
		return &bi, worked
	}

	return nil, worked
}

func (i BigInt) MarshalJSON() ([]byte, error) {
	i2 := big.Int(i)
	return []byte(fmt.Sprintf(`%s`, i2.String())), nil
}
