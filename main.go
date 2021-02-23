package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/btcsuite/btcd/btcec"
)

func main() {
	s := btcec.Signature{
		R: big.NewInt(100),
		S: big.NewInt(50),
	}
	fmt.Println(s.Serialize())
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%+v\n", err)
	os.Exit(1)
}
