package masscan

import (
	"encoding/json"
	"io"
	"os"
)

// Port struct for masscan JSON
type Port struct {
	Port     int    `json:"port"`
	Protocol string `json:"proto"`
	Status   string `json:"status"`
	Reason   string `json:"reason"`
	TTL      int    `json:"ttl"`
}

// Host struct for masscan JSON
type Host struct {
	IP        string `json:"ip"`
	Timestamp string `json:"timestamp"`
	Ports     []Port `json:"ports"`
}

func ReadMasscanJSONIPs(file string) ([]string, error) {
	// Read masscan JSON file and map to Host struct
	masscanHosts := []Host{}

	// Open file
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	// Read JSON file
	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &masscanHosts)
	if err != nil {
		return nil, err
	}

	// Create slice of IPs
	ips := make([]string, 0)
	for _, host := range masscanHosts {
		if host.IP != "" {
			for _, port := range host.Ports {
				if port.Status == "open" {
					ips = append(ips, host.IP)
					break
				}
			}
		}
	}
	return ips, nil
}
