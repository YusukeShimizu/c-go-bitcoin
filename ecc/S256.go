package ecc

import (
	"fmt"
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

func genG() (*s256Point, error) {
	gxhex := "0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
	gx, ok := new(big.Int).SetString(gxhex, 0)
	if !ok {
		return nil, xerrors.Errorf("coundn't generate x of the generator point G from hex:,%s", gxhex)
	}
	gyhex := "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"
	gy, ok := new(big.Int).SetString(gyhex, 0)
	if !ok {
		return nil, xerrors.Errorf("coundn't generate y of the generator point G from hex:,%s", gxhex)
	}
	return NewS256Point(gx, gy)
}

func NewS256Field(number *big.Int) (*fieldElement, error) {
	prime := genPrime()
	if number.Cmp(prime) >= 0 {
		return nil, xerrors.New("number is larger than prime")
	}
	return &fieldElement{number: number, prime: prime}, nil
}

type s256Point struct {
	*point
	g *s256Point
	n *big.Int
}

func NewS256Point(bx, by *big.Int) (*s256Point, error) {
	prime := genPrime()
	x, err := NewFieldElement(bx, prime)
	if err != nil {
		return nil, err
	}
	y, err := NewFieldElement(by, prime)
	if err != nil {
		return nil, err
	}
	a, err := NewFieldElement(big.NewInt(0), prime)
	if err != nil {
		return nil, err
	}
	b, err := NewFieldElement(big.NewInt(7), prime)
	if err != nil {
		return nil, err
	}
	sp, err := NewPoint(x, y, a, b)
	if err != nil {
		return nil, err
	}
	hexN := "0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"
	// We specify the order of the group generated by G, n.s
	n, ok := new(big.Int).SetString(hexN, 0)
	if !ok {
		return nil, xerrors.Errorf("coundn't generate the order of the group generated by G from hex:,%s", hexN)
	}
	g, err := genG()
	if err != nil {
		return nil, err
	}
	return &s256Point{sp, g, n}, nil
}

func (s s256Point) SRMul(coefficient *big.Int) error {
	return s.FastRMul(coefficient.Mod(coefficient, s.n))
}

func (s s256Point) Verify(z *big.Int, sig signature) (bool, error) {
	// s_inv = pow(sig.s, N - 2, N)
	s_inv := big.NewInt(0).Exp(sig.s, big.NewInt(0).Sub(s.n, big.NewInt(2)), s.n)
	// u = z * s_inv % N
	u := big.NewInt(0).Mul(z, s_inv)
	u = u.Mod(u, s.n)
	// v = sig.r * s_inv % N
	v := big.NewInt(0).Mul(sig.r, s_inv)
	v = v.Mod(v, s.n)
	// total = u * G + v * self
	g, err := genG()
	if err != nil {
		return false, err
	}
	// u * G
	err = g.SRMul(u)
	if err != nil {
		return false, err
	}
	sCopy, err := NewS256Point(s.x.number, s.y.number)
	if err != nil {
		return false, err
	}
	// v * self
	err = sCopy.SRMul(v)
	if err != nil {
		return false, err
	}
	fmt.Printf("v*self:%s", sCopy.x.number.String())
	err = g.Add(sCopy.point)
	if err != nil {
		return false, err
	}
	return g.x.number.Cmp(sig.r) == 0, nil
}
