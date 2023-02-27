package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file [filename]",
	Short: "Download filename from the application",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("missing filename")
		}

		dir, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}

		target, err := cmd.Flags().GetString("app")
		if err != nil {
			return err
		}

		found := false
		for _, val := range []string{"B", "D", "L"} {
			if dir == val {
				found = true
				break
			}
		}

		if !found {
			return errors.New("directory not supported")
		}

		name, length, err := download(target, dir, args[0], true)
		if err != nil {
			return err
		}

		logger.Info("Saved \"%s\" (%d bytes)", name, length)

		return nil
	},
}

func init() {
	fileCmd.Flags().StringP("dir", "d", "B", "dir where is file located(Bundle, Library or Directory)")
	fileCmd.Flags().StringP("app", "a", "Gadget", "application to attach to")
	dlCmd.AddCommand(fileCmd)
}
