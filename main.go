package main

import (
	"fmt"
	"mutterio/localdev/cmd"
	"os"
)

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
