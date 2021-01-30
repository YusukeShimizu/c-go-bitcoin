package ecc

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func Test_genPrime(t *testing.T) {
	got := fmt.Sprintf("%x", genPrime())
	wantHex := "fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f"
	if got != wantHex {
		t.Errorf("genPrime() = %v, want %v", got, wantHex)
	}

}

func TestNewS256Field(t *testing.T) {

	tests := []struct {
		name    string
		number  *big.Int
		want    *fieldElement
		wantErr bool
	}{
		{
			name:    "Ok if correct",
			number:  big.NewInt(1),
			want:    &fieldElement{big.NewInt(1), genPrime()},
			wantErr: false,
		},
		{
			name:    "Error if exceed prime",
			number:  big.NewInt(0).Add(genPrime(), big.NewInt(1)),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewS256Field(tt.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewS256Field() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewS256Field() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mustGetFromHex(hex string) *big.Int {
	hexb, ok := new(big.Int).SetString(hex, 0)
	if !ok {
		return nil
	}
	return hexb
}

func TestNewS256Point(t *testing.T) {
	type args struct {
		x *big.Int
		y *big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    *s256Point
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				x: mustGetFromHex("0x5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc"),
				y: mustGetFromHex("0x6aebca40ba255960a3178d6d861a54dba813d0b813fde7b5a5082628087264da"),
			},
			want: &s256Point{
				point: &point{
					x: &fieldElement{
						number: mustGetFromHex("0x5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc"),
						prime:  genPrime(),
					},
					y: &fieldElement{
						number: mustGetFromHex("0x6aebca40ba255960a3178d6d861a54dba813d0b813fde7b5a5082628087264da"),
						prime:  genPrime(),
					},
					a: &fieldElement{
						number: big.NewInt(0),
						prime:  genPrime(),
					},
					b: &fieldElement{
						number: big.NewInt(7),
						prime:  genPrime(),
					},
				},
				n: mustGetFromHex("0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"),
			},
			wantErr: false,
		},
		{
			name: "Error if not on the vurve",
			args: args{
				x: mustGetFromHex("0x5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc"),
				y: mustGetFromHex("0x1"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewS256Point(tt.args.x, tt.args.y)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewS256Point() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !got.point.Eq(tt.want.point) {
				t.Errorf("NewS256Point() = %v, want %v", got.point, tt.want.point)
			}
			if got.n.Cmp(tt.want.n) != 0 {
				t.Errorf("NewS256Point() generated N is wrong. %v, want %v", got.n, tt.want.n)
			}
		})
	}
}
