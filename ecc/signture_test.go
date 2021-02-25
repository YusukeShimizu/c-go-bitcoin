package ecc

import (
	"encoding/hex"
	"math/big"
	"reflect"
	"testing"
)

func Test_lstrip(t *testing.T) {

	tests := []struct {
		name string
		bs   []byte
		want []byte
	}{
		{
			name: "OK",
			bs:   []byte{0x11, 0x00, 0x00, 0x10},
			want: []byte{0x11, 0x10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lstrip(tt.bs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lstrip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSignature_Der_Parse(t *testing.T) {
	type fields struct {
		r *big.Int
		s *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   *big.Int
	}{
		{
			name: "OK",
			fields: fields{
				r: big.NewInt(100),
				s: big.NewInt(50),
			},
			want: mustGetFromHex("3006020164020132"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Signature{
				r: tt.fields.r,
				s: tt.fields.s,
			}
			got := s.Der()
			if hex.EncodeToString(got) != tt.want.Text(10) {
				t.Errorf("Signature.Der() = %v, want %v", hex.EncodeToString(got), tt.want.Text(10))
			}
			gotSig, err := ParseDer(got)
			if err != nil {
				t.Error(err)
			}
			if gotSig.r.Cmp(s.r) != 0 || gotSig.s.Cmp(s.s) != 0 {
				t.Errorf("ParseDer() = r:%v s:%v,  want r:%v s:%v", gotSig.r, gotSig.s, s.r, s.s)
			}
		})
	}
}
