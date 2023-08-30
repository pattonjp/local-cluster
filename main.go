package main

import (
	"fmt"
	"os"

	"github.com/pattonjp/localcluster/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	info := cmd.BuildInfo{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
	if err := cmd.Run(info); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
