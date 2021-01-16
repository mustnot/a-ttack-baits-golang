package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// Fields of Nginx log format
type Fields map[string]string

// Log of Nginx
type Log struct {
	text   string
	fields Fields
}

// NewLog is Initialize AccessLog struct
func NewLog(text string) *Log {
	fields := *Parse(text)
	log := Log{text: text, fields: fields}

	return &log
}

// Parse of Log text
func Parse(line string) *Fields {
	nginxFullRegex := regexp.MustCompile(`(?P<ipaddress>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) - ([a-z\-]+) \[(?P<datetime>\d{2}\/[a-zA-Z]{3}\/\d{4}:\d{2}:\d{2}:\d{2} [\+\-]\d{4})\] \"(?P<method>\w+) (?P<url>.+) (http|HTTP)\/1\.[0-1]" (?P<statuscode>\d{3}) (?P<bytessent>\d+) ["](?P<refferer>(\-)|(.+))["] ["](?P<useragent>.+)["]`)
	match := nginxFullRegex.FindStringSubmatch(line)

	fields := Fields{}
	for index, name := range nginxFullRegex.SubexpNames() {
		if name != "" && index > 0 && index <= len(match) {
			fields[name] = match[index]
		}
	}
	return &fields
}

func main() {
	fo, err := os.Open("example/access.log")
	if err != nil {
		fmt.Println(err)
	}
	defer fo.Close()

	reader := bufio.NewReader(fo)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		log := NewLog(line)
		fmt.Println(log.fields)
	}
}
