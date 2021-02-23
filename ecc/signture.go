package ecc

import (
	"math/big"
)

type Signature struct {
	r *big.Int
	s *big.Int
}

func NewSignature(r *big.Int, s *big.Int) *Signature {
	return &Signature{r, s}
}

func (s *Signature) Der() []byte {
	rb := canonicalizeInt(s.r)
	sb := canonicalizeInt(s.s)
	length := 6 + len(rb) + len(sb)
	b := make([]byte, length)
	b[0] = 0x30
	b[1] = byte(length - 2)
	b[2] = 0x02
	b[3] = byte(len(rb))
	offset := copy(b[4:], rb) + 4
	b[offset] = 0x02
	b[offset+1] = byte(len(sb))
	copy(b[offset+2:], sb)
	return b
}

func ParseDer(der []byte) (*Signature, error) {
	return nil, nil
}
func lstrip(bs []byte) []byte {
	lstriped := []byte{}
	for _, b := range bs {
		if b != byte(0x00) {
			lstriped = append(lstriped, b)
		}
	}
	return lstriped
}

func canonicalizeInt(val *big.Int) []byte {
	b := val.Bytes()
	if len(b) == 0 {
		b = []byte{0x00}
	}
	if b[0]&0x80 != 0 {
		paddedBytes := make([]byte, len(b)+1)
		copy(paddedBytes[1:], b)
		b = paddedBytes
	}
	return b
}
