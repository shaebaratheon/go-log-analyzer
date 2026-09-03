package parser

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var CommonPatterns = map[string]string{
	"COMMONAPACHELOG": `^(?P<client_ip>\S+) \S+ (?P<user>\S+) \[(?P<timestamp>[^\]]+)\] "(?P<method>\S+) (?P<path>\S+) (?P<proto>[^"]+)" (?P<status>\d{3}) (?P<size>\d+)`,
	"SYSLOG":          `^(?P<timestamp>[A-Z][a-z]{2}\s+\d+\s+\d+:\d+:\d+) (?P<host>\S+) (?P<process>[^:\[]+)(?:\[(?P<pid>\d+)\])?: (?P<message>.*)$`,
}

type GrokParser struct {
	compiledRegex *regexp.Regexp
	subexpNames   []string
}

func NewGrokParser(pattern string) (*GrokParser, error) {
	regexPattern := pattern
	for k, v := range CommonPatterns {
		regexPattern = strings.ReplaceAll(regexPattern, "%{"+k+"}", v)
	}

	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, err
	}
	return &GrokParser{
		compiledRegex: re,
		subexpNames:   re.SubexpNames(),
	}, nil
}

func (g *GrokParser) Parse(line string) (map[string]interface{}, error) {
	matches := g.compiledRegex.FindStringSubmatch(line)
	if matches == nil {
		return nil, errors.New("pattern mismatch")
	}

	result := make(map[string]interface{})
	for i, name := range g.subexpNames {
		if i != 0 && name != "" {
			val := matches[i]
			if num, err := strconv.Atoi(val); err == nil {
				result[name] = num
			} else {
				result[name] = val
			}
		}
	}
	return result, nil
}

type JSONParser struct{}

func (j *JSONParser) Parse(line string) (map[string]interface{}, error) {
	var out map[string]interface{}
	err := json.Unmarshal([]byte(line), &out)
	return out, err
}
