/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"networkOnline/pkg/cidrs"
	"networkOnline/pkg/helpers"
	"networkOnline/pkg/nmap"
	"os"

	"github.com/aN0mad/golog/log"

	"github.com/spf13/cobra"
)

// nmapCmd represents the nmap command
var nmapCmd = &cobra.Command{
	Use:   "nmap",
	Short: "Parse a Nmap XML file and compare against CIDR ranges",
	Long: `Parse a Nmap XML file and compare against the CIDR range file.
	This will create a CSV output for easier consumption of network data.
	Example:
	networkOnline nmap -f google.xml -c ranges.txt
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set up logging
		if cmd.Flag("debug").Value.String() == "true" {
			log.EnableDebug()
			log.Debug("Debug mode enabled")
		}

		log.Debug("nmap called")

		log.Debugf("Toggle: %s", cmd.Flag("toggle").Value.String())
		log.Debugf("CIDR: %s", cmd.Flag("cidr").Value.String())
		log.Debugf("File: %s", cmd.Flag("file").Value.String())
		log.Debugf("Output: %s", cmd.Flag("output").Value.String())
		log.Debugf("Debug: %s", cmd.Flag("debug").Value.String())

		var outFile string
		fileCidr := cmd.Flag("cidr").Value.String()
		fileNmap := cmd.Flag("file").Value.String()
		fileOut := cmd.Flag("output").Value.String()

		// Check if files exists
		if !helpers.FileExists(fileCidr) {
			log.Fatalf("File does not exist: %s", fileCidr)
		}

		if !helpers.FileExists(fileNmap) {
			log.Fatalf("File does not exist: %s", fileNmap)
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

		ips, err := nmap.ReadNmapXMLIPs(fileNmap)
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
	rootCmd.AddCommand(nmapCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nmapCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nmapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	nmapCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	nmapCmd.Flags().StringP("cidr", "c", "", "Text file of CIDR ranges to cross-reference with the Nmap XML file addresses")
	nmapCmd.Flags().StringP("file", "f", "", "Nmap XML file to parse for IP addresses")
	nmapCmd.Flags().StringP("output", "o", "output", "Output file name")
	nmapCmd.MarkFlagRequired("cidr")
	nmapCmd.MarkFlagRequired("file")
}
