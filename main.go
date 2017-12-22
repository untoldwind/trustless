package main

import (
	"fmt"
	"os"

	"github.com/untoldwind/trustless/cmds"
)

func showError(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())

	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "%+v\n", err)
}

func main() {
	cmds.Execute()
}
