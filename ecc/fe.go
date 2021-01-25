package ecc

import (
	"math/big"

	"golang.org/x/xerrors"
)

type fieldElement struct {
	number *big.Int
	prime  *big.Int
}

func NewFieldElement(number, prime *big.Int) (*fieldElement, error) {
	if number.Cmp(prime) >= 0 {
		return nil, xerrors.New("number is larger than prime")
	}
	return &fieldElement{number: number, prime: prime}, nil
}

func (f *fieldElement) Equal(other *fieldElement) bool {
	if other == nil {
		return false
	}
	return f.number.Cmp(other.number) == 0 && f.prime.Cmp(other.prime) == 0
}

func (f *fieldElement) Add(x, y *fieldElement) (*fieldElement, error) {
	if x.prime.Cmp(y.prime) != 0 {
		return nil, xerrors.New("Cannot add two numbers in different Fields")
	}

	f.number.Mod(f.number.Add(x.number, y.number), x.prime)
	return f, nil
}

func (f *fieldElement) Sub(x, y *fieldElement) (*fieldElement, error) {
	if x.prime.Cmp(y.prime) != 0 {
		return nil, xerrors.New("Cannot sub two numbers in different Fields")
	}
	f.number.Mod(f.number.Sub(x.number, y.number), x.prime)
	return f, nil
}

func (f *fieldElement) Mul(x, y *fieldElement) (*fieldElement, error) {
	if x.prime.Cmp(y.prime) != 0 {
		return nil, xerrors.New("Cannot mul two numbers in different Fields")
	}
	f.number.Mod(f.number.Mul(x.number, y.number), x.prime)
	return f, nil
}

func (f *fieldElement) Pow(x *fieldElement, exponent *big.Int) (*fieldElement, error) {
	f.number.Mod(f.number.Exp(x.number, exponent, x.prime), x.prime)
	return f, nil
}

func (f *fieldElement) Div(x, y *fieldElement) (*fieldElement, error) {
	if x.prime.Cmp(y.prime) != 0 {
		return nil, xerrors.New("Cannot div two numbers in different Fields")
	}
	exp := big.NewInt(0)
	exp.Sub(x.prime, big.NewInt(2))

	pow := big.NewInt(0)
	pow.Exp(y.number, exp, x.prime)

	pow.Mul(x.number, pow)
	f.number.Mod(pow, x.prime)
	return f, nil
}
