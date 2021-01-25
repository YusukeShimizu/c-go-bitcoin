package ecc

import (
	"math/big"

	"golang.org/x/xerrors"
)

type point struct {
	x *fieldElement
	y *fieldElement
	a *fieldElement
	b *fieldElement
}

func NewPoint(x, y, a, b *fieldElement) (*point, error) {
	// nil means infinity.
	if x == nil && y == nil {
		return &point{x, y, a, b}, nil
	}
	left, err := NewFieldElement(big.NewInt(0), x.prime)
	if err != nil {
		return nil, err
	}
	left, err = left.Mul(y, y)
	if err != nil {
		return nil, err
	}
	right, err := NewFieldElement(big.NewInt(0), x.prime)
	if err != nil {
		return nil, err
	}
	right, err = right.Pow(x, big.NewInt(3))
	if err != nil {
		return nil, err
	}
	ax, err := NewFieldElement(big.NewInt(0), x.prime)
	if err != nil {
		return nil, err
	}
	axplusb, err := ax.Mul(a, x)
	if err != nil {
		return nil, err
	}
	axplusb, err = axplusb.Add(axplusb, b)
	if err != nil {
		return nil, err
	}
	right, err = right.Add(right, axplusb)
	if err != nil {
		return nil, err
	}
	// y2 = x3 + ax + b
	if !left.Equal(right) {
		return nil, xerrors.Errorf("(%v, %v) is not on the curve", left.number, right.number)
	}
	return &point{x, y, a, b}, nil
}

func (p *point) Eq(other *point) bool {
	// TODO: nil(infinity)の場合でも正常に比較できるようにする
	return p.x.Equal(other.x) &&
		p.y.Equal(other.y) &&
		p.a.Equal(other.a) &&
		p.b.Equal(other.b)
}

func (p *point) Add(other *point) error {
	if !p.a.Equal(other.a) || !p.b.Equal(other.b) {
		return xerrors.Errorf("(%v, %v) is not on the same curve", p, other)
	}
	if p.x == nil {
		p.x = other.x
		p.y = other.y
		return nil
	}
	if other.x == nil {
		return nil
	}
	if !p.x.Equal(other.x) {
		// s = (other.y - self.y) / (other.x - self.x)
		subOtherY, err := NewFieldElement(big.NewInt(0), p.x.prime)
		if err != nil {
			return err
		}
		subOtherY, err = subOtherY.Sub(other.y, p.y)
		if err != nil {
			return err
		}
		subOtherX, err := NewFieldElement(big.NewInt(0), p.x.prime)
		if err != nil {
			return err
		}
		subOtherX, err = subOtherX.Sub(other.x, p.x)
		if err != nil {
			return err
		}
		s, err := NewFieldElement(big.NewInt(0), p.x.prime)
		if err != nil {
			return err
		}
		s, err = s.Div(subOtherY, subOtherX)
		if err != nil {
			return err
		}
		x, err := NewFieldElement(big.NewInt(0).Set(p.x.number), p.x.prime)
		if err != nil {
			return err
		}
		// x = s**2 - self.x - other.x
		x, err = x.Pow(s, big.NewInt(2))
		if err != nil {
			return err
		}
		x, err = x.Sub(x, p.x)
		if err != nil {
			return err
		}
		x, err = x.Sub(x, other.x)
		if err != nil {
			return err
		}
		// y = s * (self.x - x) - self.y
		mulS, err := NewFieldElement(big.NewInt(0), p.x.prime)
		if err != nil {
			return err
		}
		mulS, err = mulS.Sub(p.x, x)
		if err != nil {
			return err
		}
		mulS, err = mulS.Mul(s, mulS)
		if err != nil {
			return err
		}
		y, err := NewFieldElement(big.NewInt(0), p.x.prime)
		if err != nil {
			return err
		}
		y, err = y.Sub(mulS, p.y)
		if err != nil {
			return err
		}
		p.x = x
		p.y = y
		return nil
	}
	if p.Eq(other) {
		// s=(3*x1**2+a)/(2*y1)
		powX, err := NewFieldElement(big.NewInt(0).Set(p.x.number), p.x.prime)
		if err != nil {
			return err
		}
		powX, err = powX.Pow(powX, big.NewInt(2))
		if err != nil {
			return err
		}
		tripleX, err := NewFieldElement(big.NewInt(0), p.x.prime)
		if err != nil {
			return err
		}
		triple, err := NewFieldElement(big.NewInt(3), p.x.prime)
		if err != nil {
			return err
		}
		tripleX, err = tripleX.Mul(powX, triple)
		if err != nil {
			return err
		}
		tripleX, err = tripleX.Add(tripleX, p.a)
		if err != nil {
			return err
		}
		doubleY, err := NewFieldElement(big.NewInt(0).Set(p.y.number), p.x.prime)
		if err != nil {
			return err
		}
		double, err := NewFieldElement(big.NewInt(2), p.x.prime)
		if err != nil {
			return err
		}
		doubleY, err = doubleY.Mul(doubleY, double)
		if err != nil {
			return err
		}
		s, err := NewFieldElement(big.NewInt(0), p.x.prime)
		if err != nil {
			return err
		}
		s, err = s.Div(tripleX, doubleY)
		if err != nil {
			return err
		}
		// x = s**2 - 2 * self.x
		x, err := NewFieldElement(big.NewInt(0).Set(p.x.number), p.x.prime)
		if err != nil {
			return err
		}
		x, err = x.Pow(s, big.NewInt(2))
		if err != nil {
			return err
		}
		x, err = x.Sub(x, p.x)
		if err != nil {
			return err
		}
		x, err = x.Sub(x, other.x)
		if err != nil {
			return err
		}
		// y = s * (self.x - x) - self.y
		mulS, err := NewFieldElement(big.NewInt(0), p.x.prime)
		if err != nil {
			return err
		}
		mulS, err = mulS.Sub(p.x, x)
		if err != nil {
			return err
		}
		mulS, err = mulS.Mul(s, mulS)
		if err != nil {
			return err
		}
		y, err := NewFieldElement(big.NewInt(0), p.x.prime)
		if err != nil {
			return err
		}
		y, err = y.Sub(mulS, p.y)
		if err != nil {
			return err
		}
		p.x = x
		p.y = y
		return nil
	}
	if p.x.Equal(other.x) && !p.y.Equal(other.y) {
		p.x = nil
		p.y = nil
		return nil
	}
	return nil
}

func (p *point) RMul(coefficient *big.Int) error {
	product, err := NewPoint(nil, nil, p.a, p.b)
	if err != nil {
		return xerrors.New("failed to new point")
	}
	for i := big.NewInt(0); i.Cmp(coefficient) != 0; i.Add(i, big.NewInt(1)) {
		err = product.Add(p)
		if err != nil {
			return xerrors.Errorf("failed to add point, %v", err)
		}
	}
	p.x = product.x
	p.y = product.y
	return nil
}

func (p *point) FastRMul(coefficient *big.Int) error {
	coef := coefficient
	current := p
	result, err := NewPoint(nil, nil, p.a, p.b)
	if err != nil {
		return xerrors.New("failed to new point")
	}
	for {
		// coef & 1 != 0
		if big.NewInt(0).And(coef, big.NewInt(1)).Cmp(big.NewInt(0)) != 0 {
			err := result.Add(current)
			if err != nil {
				return xerrors.Errorf("failed to add point, %v", err)
			}
		}
		err := current.Add(current)
		if err != nil {
			return xerrors.Errorf("failed to add point, %v", err)
		}
		coef.Rsh(coef, 1)
		if coef.Cmp(big.NewInt(0)) == 0 {
			break
		}
	}
	p.x = result.x
	p.y = result.y
	return nil
}
