package ecc

import (
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1"
)

type PrivateKey struct {
	p      *s256Point
	secret *big.Int
}

func NewPrivateKey(secret *big.Int) (*PrivateKey, error) {
	p := &PrivateKey{}
	p.secret = new(big.Int).Set(secret)
	g, err := genG()
	if err != nil {
		return nil, err
	}
	err = g.SRMul(secret)
	if err != nil {
		return nil, err
	}
	p.p = g
	return p, nil
}

// https://github.com/btcsuite/btcd/blob/master/btcec/signature.go#L440
func (p *PrivateKey) Sign(hash []byte) (*Signature, error) {
	k := secp256k1.NonceRFC6979(p.secret, hash, nil, nil)
	inv := new(big.Int).ModInverse(k, p.p.n)
	g, err := genG()
	if err != nil {
		return nil, err
	}
	err = g.SRMul(k)
	if err != nil {
		return nil, err
	}
	r := g.n
	rMulSec := big.NewInt(0).Mul(r, p.secret)
	zRMulSec := big.NewInt(0).Add(rMulSec, big.NewInt(0).SetBytes(hash))
	zRMulSecMulKinv := big.NewInt(0).Mul(zRMulSec, inv)
	zRMulSecMulKinv = zRMulSecMulKinv.Mod(zRMulSecMulKinv, p.p.n)
	if big.NewInt(0).Div(p.p.n, big.NewInt(2)).Cmp(zRMulSecMulKinv) == -1 {
		zRMulSecMulKinv = zRMulSecMulKinv.Sub(p.p.n, zRMulSecMulKinv)
	}
	return NewSignature(r, zRMulSecMulKinv), nil
}

func (p *PrivateKey) Wif(string, error) {

}
