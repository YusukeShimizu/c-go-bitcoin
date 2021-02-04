package ecc

import "math/big"

type signature struct {
	r *big.Int
	s *big.Int
}

func newSignature(r *big.Int, s *big.Int) *signature {
	return &signature{r, s}
}
