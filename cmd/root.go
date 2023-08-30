package cmd

import (
	"fmt"

	"github.com/pattonjp/localcluster/pkg/cluster"
	"github.com/pattonjp/localcluster/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	currBuild BuildInfo
	rootCmd   = &cobra.Command{
		Use:          "localcluster",
		Short:        "this cli helps create local development environment in kubernetes",
		SilenceUsage: true,
	}
	config cluster.Config
)

func Run(build BuildInfo) error {
	currBuild = build
	if build.Version != "dev" {
		if updated := updateCheck(); updated {
			fmt.Println("application updated please run again with the latest version")
			return nil
		}
	}
	var err error
	config, err = cluster.GetConfig()
	if err != nil {
		return err
	}
	utils.MustCheckAllDeps()
	rootCmd.AddCommand(k3dCommandRoot())
	rootCmd.AddCommand(deployCommandRoot())
	rootCmd.AddCommand(versionCmd())
	return rootCmd.Execute()
}
