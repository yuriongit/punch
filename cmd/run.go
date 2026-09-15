/*
Package cmd contains all of Punch's commands.

----

yuri.dev44@outlook.com | @yuriongit

Copyright © 2026 Yuri Okeren
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yuriongit/punch/internal/loadtest/controller"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <directory>",
	Short: "Run a load test",
	Long: `The run command starts a load test with Punch. 

To use the run command, specify a directory that contains a configuration file.
If no directory is specified, Punch will return an error.`,
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("directory not provided: specify a directory")
		}
		if len(args) > 1 {
			return fmt.Errorf("too many arguments: specify only one directory")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		return controller.RunTest(args[0])
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
