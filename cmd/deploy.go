package cmd

import (
	"github.com/pattonjp/localcluster/pkg/cluster"

	"github.com/spf13/cobra"
)

func deployCommandRoot() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "manages helm packages deployed to your dev cluster",
	}

	cmd.AddCommand(update())
	cmd.AddCommand(dployeChart())
	cmd.AddCommand(availableCharts())
	return cmd

}

func update() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "deploy all with current config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cluster.Deploy(config)
		},
	}

	return cmd
}

func dployeChart() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chart",
		Short: "deploy a single chart",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cluster.DeployChart(args[0], config)
		},
	}

	return cmd
}

func availableCharts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "available",
		Short: "list all available charts per configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cluster.AvailableCharts(config)

		},
	}
	return cmd
}
