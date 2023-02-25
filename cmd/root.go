package cmd

import (
	_ "embed"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

//go:embed scripts/script.js
var scriptJS string
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
