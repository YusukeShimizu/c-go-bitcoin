package main

import (
	"fmt"
	"os"

	"golang.org/x/xerrors"
)

func main() {
	fatal(xerrors.New("go bitcoin"))
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "%+v\n", err)
	os.Exit(1)
}
