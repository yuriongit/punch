/*
Package cmd contains all of Punch's commands.

----

yuri.dev44@outlook.com | @yuriongit

Copyright © 2026 Yuri Okeren
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "punch",
	Short: "punch is a customizable, simple and lightweight HTTP load testing CLI tool.",
	Long: `Punch is a customizable, simple and lightweight HTTP load testing CLI tool.

You define a load test in a punch.json config file, specifying a
target and the requests to send. Punch then runs the test and
reports how the target performed.

To get started, create a config file; see the 'Usage' section or
https://github.com/yuriongit/punch for more info.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.punch.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
