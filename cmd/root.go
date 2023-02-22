package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var logger log.Logger

var rootCmd = &cobra.Command{
	Use:           "ihelper",
	Short:         "iOS penetration testing helpers",
	SilenceUsage:  true,
	SilenceErrors: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute(lg log.Logger) error {
	logger = lg
	return rootCmd.Execute()
}
