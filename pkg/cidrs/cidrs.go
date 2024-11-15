package cidrs

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net/netip"
	"os"
	"strings"
)

var (
	csvHeader = []string{"CIDR", "Alive", "TotalLive", "IPs"}
)

// Cidrs is a struct that holds a list of CIDRs
type Cidrs struct {
	Cidrs []Cidr
}

// Cidr is a struct that holds a CIDR and related data
type Cidr struct {
	Cidr      netip.Prefix
	Alive     bool
	IPs       []string
	TotalLive int
}

// ReadCidrsFromFile reads a file and returns a Cidrs object
func ReadCidrsFromFile(file string) (*Cidrs, error) {
	// Create Cidrs object
	CIDRS := Cidrs{}

	// Read file
	lines := []string{}
	handle, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	scanner := bufio.NewScanner(handle)
	// optionally, resize scanner's capacity for lines over 64K, see https://stackoverflow.com/a/16615559
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Parse CIDRs
	for iter, line := range lines {
		network, err := netip.ParsePrefix(line)
		if err != nil {
			return nil, fmt.Errorf("error parsing CIDR '%s' on line %d: %s", line, iter+1, err)
		}

		// Create Cidr object and add to Cidrs object
		cidr := Cidr{network, false, []string{}, 0}
		CIDRS.Cidrs = append(CIDRS.Cidrs, cidr)
	}
	return &CIDRS, nil
}

// MapIPToCIDRs maps a list of IPs to the CIDRs structs
func (c *Cidrs) MapIPToCIDRs(ips []string) error {
	for _, ip := range ips {
		for iter, cidr := range c.Cidrs {
			ipaddr, err := netip.ParseAddr(ip)
			if err != nil {
				return fmt.Errorf("error parsing IP '%s': %s", ip, err)
			}
			if cidr.Cidr.Contains(ipaddr) {
				c.Cidrs[iter].Alive = true
				c.Cidrs[iter].IPs = append(c.Cidrs[iter].IPs, ip)
				c.Cidrs[iter].TotalLive++
			}
		}
	}
	return nil
}

// ToCSV writes the CIDRs to a CSV file
func (c *Cidrs) ToCSV(file string) (string, error) {
	// Open file
	handle, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer handle.Close()

	// Write the CSV data
	writer := csv.NewWriter(handle)

	// Create lines
	lines := [][]string{csvHeader}
	for _, cidr := range c.Cidrs {
		lines = append(lines, []string{cidr.Cidr.String(), fmt.Sprintf("%t", cidr.Alive), fmt.Sprintf("%d", cidr.TotalLive), strings.Join(cidr.IPs, ",")})
	}
	// Write lines and flush
	writer.WriteAll(lines)
	writer.Flush()

	return file, nil
}
