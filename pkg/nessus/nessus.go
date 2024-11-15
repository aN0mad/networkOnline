package nessus

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
)

var (
	DOCTYPE = "NessusClientData_v2"
)

// extractDoctype looks for an element with string 'NessusClientData_v2' in the XML file.
// Nessus files should have an element with string 'NessusClientData_v2'
func extractDoctype(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "<NessusClientData_v2>") {
			parts := strings.Fields(line)
			if len(parts) == 1 {
				return strings.TrimPrefix(strings.TrimSuffix(parts[0], ">"), "<"), nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("DOCTYPE not found")
}

// verifyFile checks if the file is a Nessus XML file
func verifyFile(file string) error {
	doctype, err := extractDoctype(file)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	if doctype != DOCTYPE || filepath.Ext(file) != ".nessus" {
		return fmt.Errorf("file is not a Nessus XML file")
	}
	return nil
}

// ReadNessusXMLIPs reads a Nessus XML file and returns a string slice of IP addresses
func ReadNessusXMLIPs(file string) ([]string, error) {
	// File type check
	err := verifyFile(file)
	if err != nil {
		return nil, err
	}

	// Create IP slice
	ips := []string{}

	// Parse XML
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(file); err != nil {
		return nil, fmt.Errorf("reading file: %s", err.Error())
	}

	// Parse XML
	for _, reportHost := range doc.FindElements("//ReportHost") {
		for _, hostChild := range reportHost.ChildElements() {
			if hostChild.Tag == "HostProperties" {
				for _, hostProp := range hostChild.ChildElements() {
					if hostProp.Tag == "tag" {
						if hostProp.SelectAttrValue("name", "") == "host-ip" {
							// Host has an IP address
							ips = append(ips, hostProp.Text())
						}
					}
				}
			}
		}
	}

	return ips, nil
}
