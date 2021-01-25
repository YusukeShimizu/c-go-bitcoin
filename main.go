package main

import (
	"fmt"
	"os"
)

func main() {

}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%+v\n", err)
	os.Exit(1)
}
