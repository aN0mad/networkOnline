package nmap

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
)

var (
	DOCTYPE = "nmaprun"
)

// extractDoctype extracts the DOCTYPE from an XML file.
// Nmap files should have a DOCTYPE of 'nmaprun'
func extractDoctype(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "<!DOCTYPE") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return strings.TrimSuffix(parts[1], ">"), nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("DOCTYPE not found")
}

// verifyFile checks if the file is an Nmap XML file
func verifyFile(file string) error {
	doctype, err := extractDoctype(file)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	if doctype != DOCTYPE || filepath.Ext(file) != ".xml" {
		return fmt.Errorf("file is not an Nmap XML file")
	}
	return nil
}

// ReadNmapOutput reads an Nmap XML file and returns a string slice of IP addresses
func ReadNmapXMLIPs(file string) ([]string, error) {
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
	for _, reportHost := range doc.FindElements("//host") {
		saveHost := false
		livePort := false
		addr := ""
		addrType := ""
		for _, hostChild := range reportHost.ChildElements() {
			if hostChild.Tag == "status" {
				if strings.ToLower(hostChild.SelectAttrValue("state", "")) == "up" {
					// Host is alive
					saveHost = true
				}
			}
			if hostChild.Tag == "address" {
				if strings.ToLower(hostChild.SelectAttrValue("addrtype", "")) == "ipv4" {
					// Host has an IPv4 address
					addr = hostChild.SelectAttrValue("addr", "")
					addrType = "ipv4"
				}
				if strings.ToLower(hostChild.SelectAttrValue("addrtype", "")) == "ipv6" {
					// Host has an IPv6 address
					addr = hostChild.SelectAttrValue("addr", "")
					addrType = "ipv6"
				}
			}
			if hostChild.Tag == "ports" {
				for _, port := range hostChild.ChildElements() {
					if strings.ToLower(port.SelectElement("state").SelectAttrValue("state", "")) == "open" {
						// Host has an open port
						livePort = true
					}
					if livePort {
						break
					}
				}
				if !livePort {
					saveHost = false
				}
			}
			// If we have a live host with an open port and ipv4 address, save the IP address
			if saveHost && livePort && addrType == "ipv4" {
				ips = append(ips, addr)
			}
		}
	}
	return ips, nil
}
