package cmd

import (
	"github.com/pattonjp/localcluster/pkg/cluster"
	"github.com/pattonjp/localcluster/pkg/utils"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:          "localk8Dev",
		Short:        "local k8 dev is a k3d helper to create a local development environment",
		SilenceUsage: true,
	}
	config cluster.Config
)

func Run() error {
	var err error
	config, err = cluster.GetConfig()
	if err != nil {
		return err
	}
	utils.MustCheckAllDeps()
	rootCmd.AddCommand(k3dCommandRoot())
	rootCmd.AddCommand(deployCommandRoot())
	return rootCmd.Execute()
}
