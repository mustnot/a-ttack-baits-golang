package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"time"

	"github.com/IncSW/geoip2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hpcloud/tail"
)

// ErrorCheck is check for error
func ErrorCheck(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// Log is struct of nginx access log
type Log struct {
	IPAddress  string
	Port       int
	Datetime   string
	Method     string
	URL        string
	StatusCode int
	SentBytes  int
	Referrer   string
	UserAgent  string
}

// NewLog is inistalize access log
func NewLog(text string) *Log {
	parseResult := Parse(text)

	parseBytes, err := json.Marshal(parseResult)
	ErrorCheck(err)

	log := Log{}
	if err := json.Unmarshal(parseBytes, &log); err != nil {
		ErrorCheck(err)
	}
	return &log
}

// ParseResult is type of parse result
type ParseResult map[string]interface{}

// Parse of Log text
func Parse(text string) *ParseResult {
	nginxFullRegex := regexp.MustCompile(`(?P<ipaddress>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) (?P<port>\d+) - ([a-z\-]+) \[(?P<datetime>\d{2}/[a-zA-Z]{3}/\d{4}:\d{2}:\d{2}:\d{2} [+-]\d{4})] "(?P<method>\w+)? (?P<url>.+) ((http|HTTP)/1\.[0-1])?" (?P<statuscode>\d{3}) (?P<sentbytes>\d+) ["](?P<referrer>(-)|(.+))["] ["](?P<useragent>.+)["]`)
	match := nginxFullRegex.FindStringSubmatch(text)

	parseResult := make(ParseResult)
	for index, name := range nginxFullRegex.SubexpNames() {
		if name != "" && index > 0 && index <= len(match) {
			value := match[index]
			if intValue, err := strconv.Atoi(value); err == nil {
				parseResult[name] = intValue
			} else if timeValue, err := time.Parse("02/Jan/2006:15:04:05 -0700", value); err == nil {
				parseResult[name] = timeValue.Format("2006-01-02 15:04:05")
			} else {
				parseResult[name] = value
			}
		}
	}
	return &parseResult
}

// Lookup struct
type Lookup struct {
	reader *geoip2.CityReader
}

// NewLookup is initializer
func NewLookup() *Lookup {
	cityReader, err := geoip2.NewCityReaderFromFile("db/GeoLite2-City.mmdb")
	ErrorCheck(err)

	return &Lookup{
		reader: cityReader,
	}
}

// Country is get country in geolocation
func (l *Lookup) Country(ipaddress string) string {
	record, err := l.reader.Lookup(net.ParseIP(ipaddress))
	ErrorCheck(err)
	return record.Country.ISOCode
}

// OpenDB is sql open to mysql database
func OpenDB() *sql.DB {
	db, err := sql.Open("mysql", "nxlogd_user:nxlogd_pw@tcp(db:3306)/nxlogd_db")
	ErrorCheck(err)
	return db
}

// OpenLogFile is open nginx log file
func OpenLogFile(filepath string) *tail.Tail {
	tailConf := tail.Config{
		Follow: true,
		ReOpen: true,
	}
	t, err := tail.TailFile(filepath, tailConf)
	if err != nil {
		fmt.Println(err)
	}
	return t
}

func main() {
	db := OpenDB()
	defer db.Close()

	stmt, _ := db.Prepare(`INSERT INTO access_log (ipaddress, port, country, datetime, method, url, status_code, sent_bytes, referrer, user_agent)
						   VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	stdout := OpenLogFile("/var/log/nginx/access.log")
	lookup := NewLookup()
	for line := range stdout.Lines {
		log := NewLog(line.Text)
		fmt.Println(log)

		if log.IPAddress != "" {
			country := lookup.Country(log.IPAddress)
			_, err := stmt.Exec(log.IPAddress, log.Port, country, log.Datetime, log.Method, log.URL, log.StatusCode, log.SentBytes, log.Referrer, log.UserAgent)
			ErrorCheck(err)
		}
	}
}
