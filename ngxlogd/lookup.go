package main

import (
	"github.com/IncSW/geoip2"
	"net"
)

// Lookup struct
type Lookup struct {
	reader *geoip2.CityReader
}

// NewLookup sss
func NewLookup() *Lookup {
	cityReader, err := geoip2.NewCityReaderFromFile("GeoLite2-City.mmdb")
	if err != nil {
		panic(err)
	}

	return &Lookup{
		reader: cityReader,
	}
}

// Location is
func (l *Lookup) Location(ipaddress string) *geoip2.CityResult {
	record, err := l.reader.Lookup(net.ParseIP(ipaddress))
	if err != nil {
		panic(err)
	}
	return record
}
