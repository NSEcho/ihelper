package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var binCmd = &cobra.Command{
	Use:   "bin [AppName]",
	Short: "Download CFBundleExecutable from the application",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("missing app name")
		}

		appName := args[0]

		name, length, err := download(appName, "", "", false)
		if err != nil {
			return err
		}

		logger.Info("Saved \"%s\" (%d bytes)", name, length)

		return nil
	},
}

func init() {
	dlCmd.AddCommand(binCmd)
}
