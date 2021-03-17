package ecc

import (
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func TestPrivateKey_Sign(t *testing.T) {
	tests := []struct {
		name    string
		secret  *big.Int
		z       []byte
		wantErr bool
	}{
		{
			name:   "OK",
			secret: big.NewInt(10),
			z:      chainhash.DoubleHashB([]byte("test message")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := NewPrivateKey(tt.secret)
			got, err := p.Sign(tt.z)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKey.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Signature R: %x\n", got.r.String())
			t.Logf("Signature S: %x\n", got.s.String())
			// /Users/bruwbird/go/1.15.0/pkg/mod/github.com/btcsuite/btcd@v0.20.1-beta/btcec/signature.goをみる
			ok, _ := p.p.Verify(big.NewInt(0).SetBytes(tt.z), *got)
			if !ok {
				t.Errorf("PrivateKey.Sign() failed")
			}
		})
	}
}
