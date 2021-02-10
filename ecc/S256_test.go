package ecc

import (
	"encoding/hex"
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

func Test_s256Point_sRMul(t *testing.T) {
	type fields struct {
		point *point
		n     *big.Int
	}

	tests := []struct {
		name        string
		fields      fields
		coefficient *big.Int
		want        *point
		wantErr     bool
	}{
		{
			name: "OK",
			fields: fields{
				point: &point{
					x: &fieldElement{
						number: big.NewInt(192),
						prime:  genPrime(),
					},
					y: &fieldElement{
						number: big.NewInt(105),
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
			coefficient: big.NewInt(3),
			want: &point{
				x: &fieldElement{
					number: mustGetFromHex("0x6cccdbe1d22d7bcc12df177da0d6e6ec4b790f5da805b983d7b1bea1da916b3b"),
					prime:  genPrime(),
				},
				y: &fieldElement{
					number: mustGetFromHex("0x26a9359a5f73ddcad408ff41ce4eb5213564ff9cfecd53877a84b0ce8d209fb9"),
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := s256Point{
				point: tt.fields.point,
				n:     tt.fields.n,
			}
			if err := s.SRMul(tt.coefficient); (err != nil) != tt.wantErr {
				t.Errorf("s256Point.sRMul() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !s.Eq(tt.want) {
				t.Errorf("point.SRMul() = x:%v y:%v want x:%v y:%v", s.x.number, s.y.number, tt.want.x.number, tt.want.y.number)
			}
		})
	}
}

func Test_genG(t *testing.T) {
	tests := []struct {
		name    string
		want    *s256Point
		wantErr bool
	}{
		{
			name: "OK",
			want: &s256Point{
				point: &point{
					x: &fieldElement{
						number: mustGetFromHex("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"),
						prime:  genPrime(),
					},
					y: &fieldElement{
						number: mustGetFromHex("0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"),
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
				n: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := genG()
			if (err != nil) != tt.wantErr {
				t.Errorf("genG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Eq(tt.want.point) {
				t.Errorf("point.SRMul() = x:%v y:%v want x:%v y:%v", got.x.number, got.y.number, tt.want.x.number, tt.want.y.number)
			}
			_ = got.SRMul(got.n)
			if got.x != nil || got.y != nil {
				t.Errorf("got should be infinity but got x:%v, y:%v", got.x, got.y)
			}

		})
	}
}

func Test_Pubpoint(t *testing.T) {
	tests := []struct {
		name   string
		secret *big.Int
		x      *big.Int
		y      *big.Int
	}{
		{
			name:   "OK case1",
			secret: big.NewInt(7),
			x:      mustGetFromHex("0x5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc"),
			y:      mustGetFromHex("0x6aebca40ba255960a3178d6d861a54dba813d0b813fde7b5a5082628087264da"),
		},
		{
			name:   "OK case2",
			secret: big.NewInt(1485),
			x:      mustGetFromHex("0xc982196a7466fbbbb0e27a940b6af926c1a74d5ad07128c82824a11b5398afda"),
			y:      mustGetFromHex("0x7a91f9eae64438afb9ce6448a1c133db2d8fb9254e4546b6f001637d50901f55"),
		},
		{
			name:   "OK case3",
			secret: big.NewInt(0).Exp(big.NewInt(2), big.NewInt(128), nil),
			x:      mustGetFromHex("0x8f68b9d2f63b5f339239c1ad981f162ee88c5678723ea3351b7b444c9ec4c0da"),
			y:      mustGetFromHex("0x662a9f2dba063986de1d90c2b6be215dbbea2cfe95510bfdf23cbf79501fff82"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewS256Point(tt.x, tt.y)
			if err != nil {
				t.Fatal(err)
			}
			g, err := genG()
			if err != nil {
				t.Fatal(err)
			}
			err = g.SRMul(tt.secret)
			if err != nil {
				t.Fatal(err)
			}
			if !s.Eq(g.point) {
				t.Errorf("x:%v y:%v want x:%v y:%v", g.x.number, g.y.number, s.x.number, s.y.number)
			}
		})
	}
}

func Test_s256Point_Verify(t *testing.T) {
	type args struct {
		x   *big.Int
		y   *big.Int
		z   *big.Int
		sig signature
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Ok case1",
			args: args{
				x: mustGetFromHex("0x887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c"),
				y: mustGetFromHex("0x61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34"),
				z: mustGetFromHex("0xec208baa0fc1c19f708a9ca96fdeff3ac3f230bb4a7ba4aede4942ad003c0f60"),
				sig: signature{
					r: mustGetFromHex("0xac8d1c87e51d0d441be8b3dd5b05c8795b48875dffe00b7ffcfac23010d3a395"),
					s: mustGetFromHex("0x68342ceff8935ededd102dd876ffd6ba72d6a427a3edb13d26eb0781cb423c4"),
				},
			},
			want: true,
		},
		{
			name: "Ok case2",
			args: args{
				x: mustGetFromHex("0x887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c"),
				y: mustGetFromHex("0x61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34"),
				z: mustGetFromHex("0x7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d"),
				sig: signature{
					r: mustGetFromHex("0xeff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c"),
					s: mustGetFromHex("0xc7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6"),
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewS256Point(tt.args.x, tt.args.y)
			if err != nil {
				t.Fatal(err)
			}
			got, err := s.Verify(tt.args.z, tt.args.sig)
			if (err != nil) != tt.wantErr {
				t.Errorf("s256Point.Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("s256Point.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mustDecodeString(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func Test_s256Point_Sec(t *testing.T) {

	tests := []struct {
		name        string
		coefficient *big.Int
		compressed  bool
		wantB       []byte
	}{
		{
			name:        "OK uncompressed case1",
			coefficient: big.NewInt(123),
			compressed:  false,
			wantB:       mustDecodeString("04a598a8030da6d86c6bc7f2f5144ea549d28211ea58faa70ebf4c1e665c1fe9b5204b5d6f84822c307e4b4a7140737aec23fc63b65b35f86a10026dbd2d864e6b"),
		},
		{
			name:        "OK uncompressed case2",
			coefficient: big.NewInt(42424242),
			compressed:  false,
			wantB:       mustDecodeString("04aee2e7d843f7430097859e2bc603abcc3274ff8169c1a469fee0f20614066f8e21ec53f40efac47ac1c5211b2123527e0e9b57ede790c4da1e72c91fb7da54a3"),
		},
		{
			name:        "OK compressed case1",
			coefficient: big.NewInt(123),
			compressed:  true,
			wantB:       mustDecodeString("03a598a8030da6d86c6bc7f2f5144ea549d28211ea58faa70ebf4c1e665c1fe9b5"),
		},
		{
			name:        "OK compressed case2",
			coefficient: big.NewInt(42424242),
			compressed:  true,
			wantB:       mustDecodeString("03aee2e7d843f7430097859e2bc603abcc3274ff8169c1a469fee0f20614066f8e"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := genG()
			if err != nil {
				t.Fatal(err)
			}
			if err := g.SRMul(tt.coefficient); err != nil {
				t.Fatal(err)
			}
			if gotB := g.Sec(tt.compressed); !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("s256Point.Sec() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
