package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// LogRecord represents a normalized log entry across different source formats.
type LogRecord struct {
	RemoteAddr    string        `json:"remote_addr"`
	Timestamp     time.Time     `json:"timestamp"`
	Method        string        `json:"method"`
	Path          string        `json:"path"`
	Protocol      string        `json:"protocol"`
	StatusCode    int           `json:"status_code"`
	BodyBytesSent int64         `json:"body_bytes_sent"`
	Duration      time.Duration `json:"duration"`
	Referer       string        `json:"referer,omitempty"`
	UserAgent     string        `json:"user_agent,omitempty"`
}

// LogFormat designates supported web server log formats.
type LogFormat int

const (
	FormatCLF LogFormat = iota
	FormatNginxCombined
	FormatJSON
)

// Parser handles line-by-line parsing of structured log streams.
type Parser struct {
	format      LogFormat
	clfRegex    *regexp.Regexp
	timeLayout  string
}

// NewParser initializes a parser instance for the requested format.
func NewParser(format LogFormat) *Parser {
	// Standard Apache/Nginx Combined Regex pattern:
	// ^(\S+) \S+ \S+ \[([^\]]+)\] "(\S+) (\S+) (\S+)" (\d{3}) (\d+|-)(?: "([^"]*)" "([^"]*)")?
	pattern := `^(\S+)\s+\S+\s+\S+\s+\[([^\]]+)\]\s+"(\S+)\s+(\S+)\s+(\S+)"\s+(\d{3})\s+(\d+|-)(?:\s+"([^"]*)"\s+"([^"]*)")?`
	return &Parser{
		format:     format,
		clfRegex:   regexp.MustCompile(pattern),
		timeLayout: "02/Jan/2006:15:04:05 -0700",
	}
}

// ParseLine parses a single string line into a LogRecord.
func (p *Parser) ParseLine(line string) (*LogRecord, error) {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return nil, errors.New("empty log line")
	}

	switch p.format {
	case FormatJSON:
		return p.parseJSON(line)
	case FormatCLF, FormatNginxCombined:
		return p.parseRegex(line)
	default:
		return nil, fmt.Errorf("unsupported log format: %d", p.format)
	}
}

func (p *Parser) parseJSON(line string) (*LogRecord, error) {
	var raw struct {
		RemoteAddr string  `json:"remote_addr"`
		TimeLocal  string  `json:"time_local"`
		Request    string  `json:"request"`
		Status     int     `json:"status"`
		BytesSent  int64   `json:"bytes_sent"`
		ReqTime    float64 `json:"request_time"`
		Referer    string  `json:"http_referer"`
		UserAgent  string  `json:"http_user_agent"`
	}

	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	parsedTime, err := time.Parse(time.RFC3339, raw.TimeLocal)
	if err != nil {
		parsedTime = time.Now()
	}

	method, path, proto := parseRequestLine(raw.Request)

	return &LogRecord{
		RemoteAddr:    raw.RemoteAddr,
		Timestamp:     parsedTime,
		Method:        method,
		Path:          path,
		Protocol:      proto,
		StatusCode:    raw.Status,
		BodyBytesSent: raw.BytesSent,
		Duration:      time.Duration(raw.ReqTime * float64(time.Second)),
		Referer:       raw.Referer,
		UserAgent:     raw.UserAgent,
	}, nil
}

func (p *Parser) parseRegex(line string) (*LogRecord, error) {
	matches := p.clfRegex.FindStringSubmatch(line)
	if len(matches) < 8 {
		return nil, fmt.Errorf("line does not match CLF regex pattern: %s", line)
	}

	remoteAddr := matches[1]
	rawTime := matches[2]
	method := matches[3]
	path := matches[4]
	protocol := matches[5]
	statusStr := matches[6]
	bytesStr := matches[7]

	parsedTime, err := time.Parse(p.timeLayout, rawTime)
	if err != nil {
		parsedTime = time.Now()
	}

	status, _ := strconv.Atoi(statusStr)
	var bytesSent int64
	if bytesStr != "-" {
		bytesSent, _ = strconv.ParseInt(bytesStr, 10, 64)
	}

	var referer, userAgent string
	if len(matches) >= 9 {
		referer = matches[8]
	}
	if len(matches) >= 10 {
		userAgent = matches[9]
	}

	return &LogRecord{
		RemoteAddr:    remoteAddr,
		Timestamp:     parsedTime,
		Method:        method,
		Path:          path,
		Protocol:      protocol,
		StatusCode:    status,
		BodyBytesSent: bytesSent,
		Referer:       referer,
		UserAgent:     userAgent,
	}, nil
}

func parseRequestLine(req string) (string, string, string) {
	parts := strings.Fields(req)
	if len(parts) == 3 {
		return parts[0], parts[1], parts[2]
	} else if len(parts) == 2 {
		return parts[0], parts[1], "HTTP/1.1"
	}
	return "GET", "/", "HTTP/1.1"
}
