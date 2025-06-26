package service

import (
	"github.com/tommydehaas/go-am-geolocation-provider/internal/geo"
	"github.com/tommydehaas/go-am-geolocation-provider/internal/render"
	"log"
	"net/http"
)

type response struct {
	Ipinfo struct {
		IPAddress string `json:"ip_address,omitempty"`
		Location  struct {
			CountryData struct {
				CountryCode string `json:"country_code,omitempty"`
			} `json:"CountryData"`
		} `json:"Location"`
	} `json:"ipinfo"`
}

func GetCountyCodeByIp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	ip := r.PathValue("ip")
	country := geo.GetCountryCodeByIp(ip)

	log.Printf("[service] IP: %s, Country: %s\n", ip, country)

	resp := &response{}
	resp.Ipinfo.IPAddress = ip
	resp.Ipinfo.Location.CountryData.CountryCode = country

	render.JSON(w, resp, http.StatusOK)
	return
}
