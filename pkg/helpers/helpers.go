package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// fileExists checks if a file exists and is not a directory before returning a bool
func FileExists(file string) bool {
	// Ensure path exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	// Ensure path is not a directory
	if fi, err := os.Stat(file); err == nil && fi.IsDir() {
		return false
	}

	// File exists
	return true
}

// SplitOnLast splits a string on the last instance of a substring
func SplitOnLast(s, sep string) (string, string) {
	index := strings.LastIndex(s, sep)
	if index == -1 {
		return s, ""
	}
	return s[:index], s[index+len(sep):]
}

// Helper function to convert []interface{} to []string
func ConvertToStringSlice(interfaces []interface{}) []string {
	strings := make([]string, len(interfaces))
	for i, v := range interfaces {
		strings[i] = v.(string)
	}
	return strings
}

// ReadLinesFromFile reads a file and returns a slice of strings
func ReadLinesFromFile(file string) ([]string, error) {
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
	return lines, nil
}

// CreateOutputFile checks if the filename exists and if so, creates a new filename by adding a number to the end
func CreateOutputFile(file string, ext string) string {
	file = fmt.Sprintf("%s%s", file, ext)
	if FileExists(file) {
		// Split the filename on the last period
		filename, ext := SplitOnLast(file, ".")
		for i := 1; ; i++ {
			newFile := fmt.Sprintf("%s_%d.%s", filename, i, ext)
			if !FileExists(newFile) {
				return newFile
			}
		}
	}
	return file
}
