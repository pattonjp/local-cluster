package main

import (
	"fmt"
	"os"

	"github.com/pattonjp/localcluster/cmd"
	"github.com/pattonjp/localcluster/pkg/updater"
)

var (
	version string
	commit  = "none"
	date    = "unknown"
)

func main() {
	vMgr := updater.New(version, commit, date, "pattonjp/localcluster")

	if err := cmd.Run(vMgr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
