package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

func main() {
	// Decode a hex-encoded private key.
	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), big.NewInt(10).Bytes())
	// Sign a message using the private key.
	message := "test message"
	messageHash := chainhash.DoubleHashB([]byte(message))
	signature, err := privKey.Sign(messageHash)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Serialize and display the signature.
	fmt.Printf("Serialized Signature: %x\n", signature.Serialize())
	fmt.Printf("Signature R: %x\n", signature.R.String())
	fmt.Printf("Signature S: %x\n", signature.S.String())

	// Verify the signature for the message using the public key.
	verified := signature.Verify(messageHash, pubKey)
	fmt.Printf("Signature Verified? %v\n", verified)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%+v\n", err)
	os.Exit(1)
}
