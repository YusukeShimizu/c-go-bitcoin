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
