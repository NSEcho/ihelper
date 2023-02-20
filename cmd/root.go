package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:          "ihelper",
	Short:        "iOS penetration testing helpers",
	SilenceUsage: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() error {
	return rootCmd.Execute()
}
