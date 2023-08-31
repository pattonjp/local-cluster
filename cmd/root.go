package cmd

import (
	"github.com/pattonjp/localcluster/pkg/cluster"
	"github.com/pattonjp/localcluster/pkg/updater"
	"github.com/pattonjp/localcluster/pkg/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:          "localcluster",
		Short:        "manages local kubernetes development environment",
		SilenceUsage: true,
	}
	config cluster.Config
)

func Run(vMgr updater.VersionManager) error {
	setUsageTemplate()
	vMgr.CheckForUpdateGuarded()

	var err error
	config, err = cluster.GetConfig()
	if err != nil {
		return err
	}
	utils.MustCheckAllDeps()
	rootCmd.AddCommand(deployCommandRoot())
	rootCmd.AddCommand(versionCmd(vMgr))
	return rootCmd.Execute()
}

func setUsageTemplate() {
	cobra.AddTemplateFunc("StyleHeading", color.New(color.FgGreen).SprintFunc())
	cobra.AddTemplateFunc("StyleCommand", color.New(color.FgBlue).SprintFunc())
	cobra.AddTemplateFunc("StyleFlags", color.New(color.FgMagenta).SprintFunc())

	rootCmd.SetUsageTemplate(usageTemplate)
}

var usageTemplate = `{{ StyleHeading "Usage:"}}{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

{{ StyleHeading "Aliases:"}}
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

{{ StyleHeading "Examples:"}}
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

{{ StyleHeading "Available Commands:"}}{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding | StyleCommand }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

{{ StyleHeading "Additional Commands:"}}{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

{{ StyleHeading "Flags:"}}
{{.LocalFlags.FlagUsages | StyleFlags | trimTrailingWhitespaces }}{{end}}{{if .HasAvailableInheritedFlags}}

{{ StyleHeading "Global Flags:"}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces | StyleFlags}}{{end}}{{if .HasHelpSubCommands}}

{{ StyleHeading "Additional help topics:"}}{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} {{ StyleCommand "[command]"}} --help" for more information about a command.
{{end}}
`
