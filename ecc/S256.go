package ecc

import (
	"math/big"

	"golang.org/x/xerrors"
)

// 2**256 - 2**32 - 977
func genPrime() *big.Int {
	prime := big.NewInt(0)
	prime.Exp(big.NewInt(2), big.NewInt(256), nil)
	fraction := big.NewInt(0)
	fraction.Exp(big.NewInt(2), big.NewInt(32), nil)
	fraction.Add(fraction, big.NewInt(977))
	return prime.Sub(prime, fraction)
}

func NewS256Field(number *big.Int) (*fieldElement, error) {
	prime := genPrime()
	if number.Cmp(prime) >= 0 {
		return nil, xerrors.New("number is larger than prime")
	}
	return &fieldElement{number: number, prime: prime}, nil
}
