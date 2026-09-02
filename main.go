package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
)

// LogEntry represents a parsed log line
type LogEntry struct {
	IP        string `json:"ip"`
	Timestamp string `json:"timestamp"`
	Method    string `json:"method"`
	URL       string `json:"url"`
	Status    string `json:"status"`
	Size      string `json:"size"`
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
	formatPtr := flag.String("output", "text", "Output format: text or json")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: go-log-analyzer [-output=json] <log-file>")
		os.Exit(1)
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var entries []*LogEntry

	for scanner.Scan() {
		entry, err := parseLine(scanner.Text())
		if err != nil {
			continue
		}
		entries = append(entries, entry)
	}

	if *formatPtr == "json" {
		jsonData, err := json.MarshalIndent(entries, "", "  ")
		if err != nil {
			fmt.Printf("Error encoding JSON: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(jsonData))
	} else {
		fmt.Printf("Analysis Complete:\n")
		fmt.Printf("Successfully parsed: %d lines\n", len(entries))
	}
}
