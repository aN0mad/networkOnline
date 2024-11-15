/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	outFileExt = ".csv"
	debug      bool
	VERSION    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "networkOnline",
	Short: "A parser for multiple output formats to compare against CIDR ranges",
	Long: `A parser for multiple output formats to compare against CIDR ranges in order to determine
which networks are online. This tool is useful for creating a list of online networks for further
testing or analysis within the output CSV file. 
The tool currently supports the following formats:
- masscan
- nmap
- nessus
- text`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("version").Value.String() == "true" {
			cmd.Println("Version:", VERSION)
			os.Exit(0)
		}
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v string) {
	VERSION = v
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug output")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "", false, "Print the version number")
	// rootCmd.Flags().BoolP("debug", "d", false, "Enable debug output")
}
