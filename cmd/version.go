package cmd

import (
	"fmt"

	"github.com/pattonjp/localcluster/pkg/updater"
	"github.com/spf13/cobra"
)

var appRepo = "pattonjp/localcluster"

func versionCmd(vMgr updater.VersionManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "displays the current version",
		Run: func(cmd *cobra.Command, args []string) {
			vMgr.Print()
		},
	}
	cmd.AddCommand(updateCmd(vMgr))
	return cmd
}

func updateCmd(vMgr updater.VersionManager) *cobra.Command {
	cmd := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			err := vMgr.UpdateToLatest()
			if err != nil {
				fmt.Printf("Could not update version \n %v+", err)
			}
		},
	}
	return cmd
}
