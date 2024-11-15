/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"networkOnline/pkg/cidrs"
	"networkOnline/pkg/helpers"
	"os"

	"github.com/aN0mad/golog/log"
	"github.com/spf13/cobra"
)

// textCmd represents the text command
var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Parse a text file and compare against CIDR ranges",
	Long: `Parse a text file and compare against the CIDR range file.
	This will create a CSV output for easier consumption of network data.
	Example:
	networkOnline text -f google.txt -c ranges.txt
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set up logging
		if cmd.Flag("debug").Value.String() == "true" {
			log.EnableDebug()
			log.Debug("Debug mode enabled")
		}

		fmt.Println("text called")
		log.Debug("Toggle:", cmd.Flag("toggle").Value)
		log.Debug("CIDR:", cmd.Flag("cidr").Value)
		log.Debug("File:", cmd.Flag("file").Value)
		log.Debug("Output:", cmd.Flag("output").Value)
		log.Debug("Debug:", cmd.Flag("debug").Value)

		var outFile string
		fileCidr := cmd.Flag("cidr").Value.String()
		fileText := cmd.Flag("file").Value.String()
		fileOut := cmd.Flag("output").Value.String()

		// Check if files exists
		if !helpers.FileExists(fileCidr) {
			log.Fatalf("File does not exist: %s", fileCidr)
		}

		if !helpers.FileExists(fileText) {
			log.Fatalf("File does not exist: %s", fileText)
		}

		outFile = helpers.CreateOutputFile(fileOut, outFileExt)
		log.Infof("Output file: %s", outFile)

		// Read CIDR file and create a struct for each CIDR
		CIDRS, err := cidrs.ReadCidrsFromFile(fileCidr)
		if err != nil {
			log.Errorf("error reading CIDR file: %s", err)
			os.Exit(1)
		}
		log.Infof("CIDRs created: %d", len(CIDRS.Cidrs))

		ips, err := helpers.ReadLinesFromFile(fileText)
		if err != nil {
			log.Errorf("error reading IP file: %s", err)
			os.Exit(1)
		}
		log.Infof("IPs read: %d", len(ips))

		// Verify if IPs are in CIDRs and modify CIDRs struct accordingly
		err = CIDRS.MapIPToCIDRs(ips)
		if err != nil {
			log.Errorf("error mapping IPs to CIDRs: %s", err)
			os.Exit(1)
		}
		log.Infof("IPs mapped to CIDRs")

		// Write CIDRs to CSV
		_, err = CIDRS.ToCSV(outFile)
		if err != nil {
			log.Errorf("error writing CIDRs to CSV: %s", err)
			os.Exit(1)
		}
		log.Infof("CIDRs written to CSV: %s", outFile)
	},
}

func init() {
	rootCmd.AddCommand(textCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// textCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// textCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	textCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	textCmd.Flags().StringP("cidr", "c", "", "Text file of CIDR ranges to cross-reference with the Text file addresses")
	textCmd.Flags().StringP("file", "f", "", "Text file to parse for IP addresses")
	textCmd.Flags().StringP("output", "o", "output", "Output file name")
	textCmd.MarkFlagRequired("cidr")
	textCmd.MarkFlagRequired("file")
}
