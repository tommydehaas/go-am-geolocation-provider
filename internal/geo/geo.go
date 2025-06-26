package geo

import (
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"sync/atomic"
	"time"
)

var cd countryDatabase

type countryDatabase struct {
	db atomic.Pointer[geoip2.Reader]
}

func GetCountryCodeByIp(ip string) string {
	reader := cd.db.Load()
	if reader == nil {
		return "ZZ"
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "ZZ"
	}

	rec, err := reader.Country(parsedIP)
	if err != nil || len(rec.Country.IsoCode) == 0 {
		return "ZZ"
	}

	return rec.Country.IsoCode
}

func init() {
	loadDB := func() *geoip2.Reader {
		db, err := geoip2.Open("/usr/share/GeoIP/GeoLite2-Country.mmdb")
		if err != nil {
			log.Printf("[db] Failed to load GeoIP DB: %v\n", err)
			return nil
		}

		log.Println("[db] GeoIP DB refreshed")
		return db
	}

	go func() {
		for {
			newDB := loadDB()
			if newDB != nil {
				oldDB := cd.db.Swap(newDB)
				if oldDB != nil {
					oldDB.Close()
				}
			}
			time.Sleep(time.Hour)
		}
	}()
}
