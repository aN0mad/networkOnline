/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"networkOnline/pkg/cidrs"
	"networkOnline/pkg/helpers"
	"networkOnline/pkg/nessus"
	"os"

	"github.com/aN0mad/golog/log"
	"github.com/spf13/cobra"
)

// nessusCmd represents the nessus command
var nessusCmd = &cobra.Command{
	Use:   "nessus",
	Short: "Parse a Nessus XML file and compare against CIDR ranges",
	Long: `Parse a Nessus XML file and compare against the CIDR range file.
	This will create a CSV output for easier consumption of network data.
	Example:
	networkOnline nessus -f google.nessus -c ranges.txt
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set up logging
		if cmd.Flag("debug").Value.String() == "true" {
			log.EnableDebug()
			log.Debug("Debug mode enabled")
		}

		fmt.Println("nessus called")
		log.Debug("Toggle:", cmd.Flag("toggle").Value)
		log.Debug("CIDR:", cmd.Flag("cidr").Value)
		log.Debug("File:", cmd.Flag("file").Value)
		log.Debug("Output:", cmd.Flag("output").Value)
		log.Debug("Debug:", cmd.Flag("debug").Value)

		var outFile string
		fileCidr := cmd.Flag("cidr").Value.String()
		fileNessus := cmd.Flag("file").Value.String()
		fileOut := cmd.Flag("output").Value.String()

		// Check if files exists
		if !helpers.FileExists(fileCidr) {
			log.Fatalf("File does not exist: %s", fileCidr)
		}

		if !helpers.FileExists(fileNessus) {
			log.Fatalf("File does not exist: %s", fileNessus)
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

		ips, err := nessus.ReadNessusXMLIPs(fileNessus)
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
	rootCmd.AddCommand(nessusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nessusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nessusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	nessusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	nessusCmd.Flags().StringP("cidr", "c", "", "Text file of CIDR ranges to cross-reference with the Nessus XML file addresses")
	nessusCmd.Flags().StringP("file", "f", "", "Nessus XML file to parse for IP addresses")
	nessusCmd.Flags().StringP("output", "o", "output", "Output file name")
	nessusCmd.MarkFlagRequired("cidr")
	nessusCmd.MarkFlagRequired("file")
}
