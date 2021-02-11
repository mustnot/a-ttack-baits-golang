package main

import (
	"net"
	"regexp"
)

// Log is struct of nginx access log
type Log struct {
	ipaddress net.IP
	// country    string
	datetime   string
	method     string
	url        string
	statuscode string
	sentbytes  string
	referrer   string
	useragent  string
}

// NewLog is inistalize access log
func NewLog(text string) *Log {
	parseResult := *Parse(text)

	return &Log{
		ipaddress:  net.ParseIP(parseResult["ipaddress"]),
		datetime:   parseResult["datetime"],
		url:        parseResult["url"],
		statuscode: parseResult["statuscode"],
		sentbytes:  parseResult["bytessent"],
		referrer:   parseResult["referrer"],
		useragent:  parseResult["useragent"],
	}
}

// ParseResult is type of parse result
type ParseResult map[string]string

// Parse of Log text
func Parse(line string) *ParseResult {
	nginxFullRegex := regexp.MustCompile(`(?P<ipaddress>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) - ([a-z\-]+) \[(?P<datetime>\d{2}/[a-zA-Z]{3}/\d{4}:\d{2}:\d{2}:\d{2} [+-]\d{4})] "(?P<method>\w+) (?P<url>.+) (http|HTTP)/1\.[0-1]" (?P<statuscode>\d{3}) (?P<bytessent>\d+) ["](?P<referrer>(-)|(.+))["] ["](?P<useragent>.+)["]`)
	match := nginxFullRegex.FindStringSubmatch(line)

	parseResult := ParseResult{}
	for index, name := range nginxFullRegex.SubexpNames() {
		if name != "" && index > 0 && index <= len(match) {
			parseResult[name] = match[index]
		}
	}
	return &parseResult
}
