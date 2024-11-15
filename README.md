# networkOnline
A go based project that ingests a list of CIDR blocks and a list of IP addresses and returns a CSV marking a CIDR as alive if an IP address was within the block.

## Background/Logic
There are a few things that a user should be aware of when using this program.
1. The nmap and masscan modules verify a port is 'open' before including the host in the slice to be compared. If a host is found but no port is open then the host will be discarded.
2. Output file names will increment automatically if a file with the same name exists. Example: using the default file name of `output.csv`, when incremented will become `output_1.csv` then `output_2.csv` etc.


## Makefile
There is a makefile in the repository that makes building and cleaning the folder structure easier

```
root@f6a960422947:/workspace# make help
 help: print this help message
 tidy: format code and tidy modfile
 build: build the unix version
 buildwin: build the windows version
 all: build all applications for unix and windows
 clean: clean the repository
```

## Build
The easiest way is to compile with the included Docker container or just on Linux in general. Golang cross-compilation makes this easy.
### Linux
```bash
make build
```

### Windows
```bash
make buildwin
```

## Usage
### File: `small.json`
```json
[
{   "ip": "10.50.2.118",   "timestamp": "1731615638", "ports": [ {"port": 5985, "proto": "tcp", "status": "open", "reason": "rst-ack", "ttl": 64} ] }
,
{   "ip": "10.50.1.43",   "timestamp": "1731615639", "ports": [ {"port": 3306, "proto": "tcp", "status": "closed", "reason": "rst-ack", "ttl": 64} ] }
]
```

### File: `ranges.txt`
```text
10.50.1.0/24
10.50.2.0/24
10.50.3.0/24
```

### Masscan example
```
$ ./bin/networkOnline.elf masscan -f small.json -c ranges.txt 
2024/11/15 20:00:44 INFO masscan called
2024/11/15 20:00:44 INFO Output file: output_1.csv
2024/11/15 20:00:44 INFO CIDRs created: 3
2024/11/15 20:00:44 INFO IPs read: 1
2024/11/15 20:00:44 INFO IPs mapped to CIDRs
2024/11/15 20:00:44 INFO CIDRs written to CSV: output.csv
```

### Output: `output.csv`
```csv
CIDR,Alive,TotalLive,IPs
10.50.1.0/24,false,0,
10.50.2.0/24,true,1,10.50.2.118
10.50.3.0/24,false,0,
```

## Help
```
$ ./bin/networkOnline.elf -h                                                
A parser for multiple output formats to compare against CIDR ranges in order to determine
which networks are online. This tool is useful for creating a list of online networks for further
testing or analysis within the output CSV file. 
The tool currently supports the following formats:
- masscan
- nmap
- nessus
- text

Usage:
  networkOnline [flags]
  networkOnline [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  masscan     Parse a masscan json file and compare against CIDR ranges
  nessus      Parse a Nessus XML file and compare against CIDR ranges
  nmap        Parse a Nmap XML file and compare against CIDR ranges
  text        Parse a text file and compare against CIDR ranges

Flags:
      --debug     Enable debug output
  -h, --help      help for networkOnline
  -t, --toggle    Help message for toggle
      --version   Print the version number

Use "networkOnline [command] --help" for more information about a command.
```