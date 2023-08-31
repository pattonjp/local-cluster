package cmd

import (
	"fmt"

	"github.com/pattonjp/localcluster/pkg/cluster"
	"github.com/pattonjp/localcluster/pkg/utils"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(list())
	rootCmd.AddCommand(create())
	rootCmd.AddCommand(recreate())
	rootCmd.AddCommand(delete())
	rootCmd.AddCommand(start())
	rootCmd.AddCommand(stop())
	rootCmd.AddCommand(use())
	rootCmd.AddCommand(setup())
}

func list() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "lists all local clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cluster.List(true)
		},
	}
	return cmd
}

func use() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "use",
		Short: "sets the local environment to a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				config.Name = args[0]
			}
			return cluster.SetContext(true, config.Name)

		},
	}
	return cmd
}
func start() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start the cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				config.Name = args[0]

			}
			if err := cluster.Start(true, config.Name); err != nil {
				return err
			}
			cluster.SetContext(false, config.Name)

			return nil
		},
	}
	return cmd
}

func stop() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop the cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				config.Name = args[0]
			}
			if err := cluster.Stop(true, config.Name); err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}

func create() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "creates a new local cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				config.Name = args[0]
			}

			if err := cluster.Start(false, config.Name); err == nil {
				fmt.Println("cluster with that name was already created. exiting after starting the cluster")
				cluster.Deploy(config)
				return nil
			}

			if err := cluster.Create(true, config); err != nil {
				return err
			}

			return cluster.Deploy(config)

		},
	}
	cmd.Flags().AddFlagSet(config.GetFlagSet())
	return cmd
}

func delete() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "deletes the local cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				config.Name = args[0]
			}
			return cluster.Delete(true, args[0])
		},
	}

	return cmd
}

func recreate() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "recreate",
		Short: "deletes and creates the local cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				config.Name = args[0]
			}
			delCMD := delete()
			delCMD.SetArgs([]string{config.Name})
			if err := delCMD.Execute(); err != nil {
				return err
			}
			createCMd := create()
			createCMd.SetArgs([]string{config.Name})
			if err := createCMd.Execute(); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().AddFlagSet(config.GetFlagSet())
	return cmd
}

func setup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "ensures the charts are available locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := utils.InstallAllDeps(); err != nil {
				return err
			}
			if err := cluster.Setup(true, config); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}
