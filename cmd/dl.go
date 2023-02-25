package cmd

import (
	"github.com/spf13/cobra"
)

var dlCmd = &cobra.Command{
	Use:   "dl [filepath]",
	Short: "Download file/binary from the application",
}

func init() {
	rootCmd.AddCommand(dlCmd)
}
