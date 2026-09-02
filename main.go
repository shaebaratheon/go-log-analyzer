package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// LogEntry represents a parsed log line
type LogEntry struct {
	IP        string
	Timestamp string
	Method    string
	URL       string
	Status    string
	Size      string
}

// Log regex pattern for standard combined format
var logPattern = regexp.MustCompile(`^(\S+) \S+ \S+ \[([^\]]+)\] "(\S+) (\S+) \S+" (\d{3}) (\d+|-)`)

func parseLine(line string) (*LogEntry, error) {
	matches := logPattern.FindStringSubmatch(line)
	if len(matches) < 7 {
		return nil, fmt.Errorf("failed to parse line")
	}

	return &LogEntry{
		IP:        matches[1],
		Timestamp: matches[2],
		Method:    matches[3],
		URL:       matches[4],
		Status:    matches[5],
		Size:      matches[6],
	}, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-log-analyzer <log-file>")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	successCount := 0
	errorCount := 0

	for scanner.Scan() {
		_, err := parseLine(scanner.Text())
		if err != nil {
			errorCount++
			continue
		}
		successCount++
	}

	fmt.Printf("Analysis Complete:\n")
	fmt.Printf("Successfully parsed: %d lines\n", successCount)
	fmt.Printf("Failed to parse:     %d lines\n", errorCount)
}
