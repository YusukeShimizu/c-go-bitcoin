package ecc

import (
	"math/big"
)

type PrivateKey struct {
	p      *s256Point
	secret *big.Int
}

func NewPrivateKey(secret *big.Int) (*PrivateKey, error) {
	p := &PrivateKey{}
	p.secret = secret
	g, err := genG()
	if err != nil {
		return nil, err
	}
	err = g.SRMul(secret)
	if err != nil {
		return nil, err
	}
	p = &PrivateKey{
		p:      g,
		secret: secret,
	}
	return p, nil
}

func (p *PrivateKey) Wif(string, error) {

}
