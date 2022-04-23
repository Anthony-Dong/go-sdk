package utils

import (
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func AddCmd(cmd *cobra.Command, foo func() (*cobra.Command, error)) error {
	subCmd, err := foo()
	if err != nil {
		return err
	}
	cmd.AddCommand(subCmd)
	return nil
}

func ExitError(err error) {
	if err != nil {
		color.Red("%s", err.Error())
	}
	os.Exit(1)
}

const (
	HelpTmpl = `{{with (or .Long .Short)}}Name: {{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
	UsageTmpl = `Usage: {{if .Runnable}}{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}{{.CommandPath}} [OPTIONS] COMMAND{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Options:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Options:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} COMMAND --help" for more information about a command.{{end}}

To get more help with gotool, check out our guides at https://docs.docker.com/go/guides/
`
)
